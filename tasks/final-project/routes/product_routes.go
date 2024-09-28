package routes

import (
	"database/sql"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"final-project/config"
	"final-project/domain"
	"final-project/middleware"
	"final-project/repository"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func ProductRoutes(r *gin.Engine, db *sql.DB) {
	// Get all products
	r.GET("/products", func(c *gin.Context) {
		limitStr := c.DefaultQuery("limit", "10")  // Default limit is 10 if not specified
		offsetStr := c.DefaultQuery("offset", "0") // Default offset is 0 if not specified
		search := c.DefaultQuery("search", "")     // Default search is empty if not specified
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
			return
		}

		offset, err := strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset parameter"})
			return
		}

		products, err := repository.GetAllProducts(db, limit, offset, search)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get products"})
			return
		}
		c.JSON(http.StatusOK, products)
	})

	// Get a product by UUID
	r.GET("/products/:uuid", func(c *gin.Context) {
		uuid := c.Param("uuid")
		if uuid == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "UUID is required"})
			return
		}

		// Call the repository function to get the product by UUID
		product, err := repository.GetProductByUUID(db, uuid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get product"})
			return
		}
		if product.UUID == "" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		c.JSON(http.StatusOK, product)
	})

	// Middleware to protect routes
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())

	// Initialize validator
	validate := validator.New()

	// Create product
	protected.POST("/products", func(c *gin.Context) {
		cld, err := config.InitializeCloudinary()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize Cloudinary"})
			return
		}

		name := c.PostForm("name")
		if name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Product name is required"})
			return
		}

		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
			return
		}

		allowedExtensions := map[string]bool{
			".jpg":  true,
			".jpeg": true,
			".png":  true,
			".svg":  true,
		}
		ext := filepath.Ext(file.Filename)
		if !allowedExtensions[ext] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Only JPG, JPEG, PNG, and SVG are allowed."})
			return
		}
		if file.Size > 5*1024*1024 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "File size exceeds 5 MB"})
			return
		}

		// Create the temporary file path
		tempFilePath := "./temp/files/" + file.Filename
		if err := os.MkdirAll(filepath.Dir(tempFilePath), 0755); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create temp directory"})
			return
		}

		// Save uploaded file to the temporary directory
		if err := c.SaveUploadedFile(file, tempFilePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file locally"})
			return
		}

		// Upload the image to Cloudinary
		uploadResult, err := config.UploadImage(cld, tempFilePath, "products")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image to Cloudinary"})
			os.Remove(tempFilePath)
			return
		}

		// Remove the temporary file after successful upload
		os.Remove(tempFilePath)

		productUUID := uuid.New().String()
		adminID := c.MustGet("admin_id").(int)
		product := domain.Product{
			UUID:     productUUID,
			Name:     name,
			ImageURL: uploadResult.SecureURL,
			AdminID:  adminID,
		}

		if err := validate.Struct(product); err != nil {
			validationErrors := err.(validator.ValidationErrors)
			errors := make(map[string]string)
			for _, fieldError := range validationErrors {
				errors[fieldError.Field()] = fieldError.Error()
			}
			c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
			return
		}

		// Insert into database
		err = repository.InsertProduct(db, &product)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Product created successfully", "product": product})
	})

	// Update product
	protected.PUT("/products/:uuid", func(c *gin.Context) {
		productUUID := c.Param("uuid")
		if productUUID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Product UUID is required"})
			return
		}

		adminID := c.MustGet("admin_id").(int)
		existingProduct, err := repository.GetProductByUUID(db, productUUID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve product"})
			return
		}
		if existingProduct.ID == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		if existingProduct.AdminID != adminID {
			c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to update this product"})
			return
		}

		name := c.PostForm("name")
		file, _ := c.FormFile("file")
		if name == "" && file == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Either product name or file must be provided"})
			return
		}

		updatedProduct := existingProduct
		if name != "" {
			updatedProduct.Name = name
		}

		if file != nil {
			allowedExtensions := map[string]bool{
				".jpg":  true,
				".jpeg": true,
				".png":  true,
				".svg":  true,
			}
			ext := filepath.Ext(file.Filename)
			if !allowedExtensions[ext] {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Only JPG, JPEG, PNG, and SVG are allowed."})
				return
			}
			if file.Size > 5*1024*1024 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "File size exceeds 5 MB"})
				return
			}

			// Create the temporary file path
			tempFilePath := "./temp/files/" + file.Filename
			if err := os.MkdirAll(filepath.Dir(tempFilePath), 0755); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create temp directory"})
				return
			}

			// Save the uploaded file to the temporary directory
			if err := c.SaveUploadedFile(file, tempFilePath); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file locally"})
				return
			}

			// Upload the image to Cloudinary
			cld, err := config.InitializeCloudinary()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize Cloudinary"})
				return
			}
			uploadResult, err := config.UploadImage(cld, tempFilePath, "products")
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image to Cloudinary"})
				os.Remove(tempFilePath)
				return
			}

			// Remove the temporary file after successful upload
			os.Remove(tempFilePath)

			updatedProduct.ImageURL = uploadResult.SecureURL
		}

		// Update the product in the database
		err = repository.UpdateProduct(db, productUUID, &updatedProduct)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Product updated successfully", "product": updatedProduct})
	})

	// Delete product
	protected.DELETE("/products/:uuid", func(c *gin.Context) {
		productUUID := c.Param("uuid")
		if productUUID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Product UUID is required"})
			return
		}

		adminID := c.MustGet("admin_id").(int)
		existingProduct, err := repository.GetProductByUUID(db, productUUID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve product"})
			return
		}

		// Check if the product exists and the user is the owner
		if existingProduct.ID == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		if existingProduct.AdminID != adminID {
			c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to delete this product"})
			return
		}

		// Delete the product from the database
		err = repository.DeleteProduct(db, productUUID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
	})
}
