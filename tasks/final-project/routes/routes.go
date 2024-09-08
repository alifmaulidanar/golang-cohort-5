package routes

import (
	"database/sql"
	"net/http"
	"path/filepath"
	"strconv"

	"final-project/config"
	"final-project/domain"
	"final-project/middleware"
	"final-project/repository"
	"final-project/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func SetupRoutes(r *gin.Engine, db *sql.DB) {
	// Register Admin
	r.POST("/auth/register", func(c *gin.Context) {
		var admin domain.Admin

		if err := c.ShouldBind(&admin); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validate := validator.New()
		if err := validate.Struct(admin); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
			return
		}

		admin.UUID = uuid.New().String()

		hashedPassword, err := service.HashPassword(admin.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		admin.Password = hashedPassword

		if err := repository.InsertAdmin(db, admin); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register admin"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Admin registered successfully", "uuid": admin.UUID})
	})

	// Login Admin
	r.POST("/auth/login", func(c *gin.Context) {
		var loginData domain.LoginRequest

		// Bind form data to LoginRequest struct
		if err := c.ShouldBind(&loginData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Find the admin by email
		admin, err := repository.FindAdminByEmail(db, loginData.Email)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		// Check if the provided password matches the stored hashed password
		if !service.CheckPasswordHash(loginData.Password, admin.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		// Generate JWT token for authenticated admin
		token, err := service.GenerateJWT(admin)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": token})
	})

	// Get all products with pagination
	r.GET("/products", func(c *gin.Context) {
		limitStr := c.DefaultQuery("limit", "10")  // Default limit is 10 if not specified
		offsetStr := c.DefaultQuery("offset", "0") // Default offset is 0 if not specified

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

		products, err := repository.GetAllProducts(db, limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get products"})
			return
		}

		c.JSON(http.StatusOK, products)
	})

	// Get a product by UUID
	r.GET("/products/:uuid", func(c *gin.Context) {
		uuid := c.Param("uuid") // Get the UUID from path parameter

		// Check if UUID is provided
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

		// Check if the product was not found
		if product.UUID == "" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}

		c.JSON(http.StatusOK, product)
	})

	// Gunakan middleware untuk endpoint yang membutuhkan otorisasi
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())

	// POST route to create a new product
	// POST route to create a new product
	protected.POST("/products", func(c *gin.Context) {
		// Initialize Cloudinary client
		cld, err := config.InitializeCloudinary()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize Cloudinary"})
			return
		}

		// Get the product name from the form-data
		name := c.PostForm("name")
		if name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Product name is required"})
			return
		}

		// Get the file from the form-data
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
			return
		}

		// Validate file type and size
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

		// Save the file locally temporarily
		tempFilePath := "./" + file.Filename
		if err := c.SaveUploadedFile(file, tempFilePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file locally"})
			return
		}

		// Upload the image to Cloudinary
		uploadResult, err := config.UploadImage(cld, tempFilePath, "products")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image to Cloudinary"})
			return
		}

		// Generate a new UUID for the product
		productUUID := uuid.New().String()

		// Extract admin_id from JWT
		adminID := c.MustGet("admin_id").(int)

		// Create a new product
		product := domain.Product{
			UUID:     productUUID,
			Name:     name,
			ImageURL: uploadResult.SecureURL,
			AdminID:  adminID,
		}

		// Insert the product into the database and get the complete data
		err = repository.InsertProduct(db, &product)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Product created successfully", "product": product})
	})
}
