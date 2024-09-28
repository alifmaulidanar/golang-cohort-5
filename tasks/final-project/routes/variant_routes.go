package routes

import (
	"database/sql"
	"net/http"
	"strconv"

	"final-project/domain"
	"final-project/middleware"
	"final-project/repository"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func VariantRoutes(r *gin.Engine, db *sql.DB) {
	r.GET("/products/variant", func(c *gin.Context) {
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

		variants, err := repository.GetAllVariants(db, limit, offset, search)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get variants"})
			return
		}
		c.JSON(http.StatusOK, variants)
	})

	// Get a variant by UUID
	r.GET("/products/variant/:uuid", func(c *gin.Context) {
		uuid := c.Param("uuid")
		if uuid == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "UUID is required"})
			return
		}

		// Call the repository function to get the product by UUID
		product, err := repository.GetVariantByUUID(db, uuid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get product"})
			return
		}
		if product.UUID == "" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Variant not found"})
			return
		}
		c.JSON(http.StatusOK, product)
	})

	// Gunakan middleware untuk endpoint yang membutuhkan otorisasi
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())

	// Create route for a variant
	protected.POST("/products/variants", func(c *gin.Context) {
		adminID := c.MustGet("admin_id").(int)
		variantName := c.PostForm("variant_name")
		quantityStr := c.PostForm("quantity")
		productIDStr := c.PostForm("product_id")
		if variantName == "" || quantityStr == "" || productIDStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "All fields (variant_name, quantity, product_id) are required"})
			return
		}

		// Convert quantity to integers
		quantity, err := strconv.Atoi(quantityStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quantity format"})
			return
		}

		// Validate that the product exists and the user is the owner
		existingProduct, err := repository.GetProductByUUID(db, productIDStr)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve variant"})
			return
		}
		if existingProduct.ID == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Variant not found"})
			return
		}
		if existingProduct.AdminID != adminID {
			c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to add a variant to this product"})
			return
		}

		variantUUID := uuid.New().String()
		variant := domain.Variant{
			UUID:        variantUUID,
			VariantName: variantName,
			Quantity:    quantity,
			ProductID:   existingProduct.ID,
		}

		validate := validator.New()
		if err := validate.Struct(variant); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err = repository.InsertVariant(db, &variant)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create variant"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Variant created successfully", "variant": variant})
	})

	// Update route for a variant
	protected.PUT("/products/variants/:uuid", func(c *gin.Context) {
		variantUUID := c.Param("uuid")
		if variantUUID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Variant UUID is required"})
			return
		}

		adminID := c.MustGet("admin_id").(int)
		existingVariant, err := repository.GetVariantByUUID(db, variantUUID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve variant"})
			return
		}
		if existingVariant.ID == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Variant not found"})
			return
		}

		// Check if the user is the owner of the product associated with the variant
		product, err := repository.GetProductByID(db, existingVariant.ProductID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve product"})
			return
		}
		if product.AdminID != adminID {
			c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to update this variant"})
			return
		}

		variantName := c.PostForm("variant_name")
		quantityStr := c.PostForm("quantity")
		if variantName == "" && quantityStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Either variant name or quantity must be provided", "req body": c.Request.Body})
			return
		}

		updatedVariant := existingVariant
		if variantName != "" {
			updatedVariant.VariantName = variantName
		}
		if quantityStr != "" {
			quantity, err := strconv.Atoi(quantityStr)
			if err != nil || quantity < 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quantity provided"})
				return
			}
			updatedVariant.Quantity = quantity
		}
		if updatedVariant == existingVariant {
			c.JSON(http.StatusOK, gin.H{"message": "No changes detected. Variant remains unchanged", "variant": updatedVariant})
			return
		}

		err = repository.UpdateVariant(db, variantUUID, &updatedVariant)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update variant"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Variant updated successfully", "variant": updatedVariant})
	})

	// Delete route for a variant
	protected.DELETE("/products/variants/:uuid", func(c *gin.Context) {
		variantUUID := c.Param("uuid")
		if variantUUID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Variant UUID is required"})
			return
		}

		adminID := c.MustGet("admin_id").(int)
		existingVariant, err := repository.GetVariantByUUID(db, variantUUID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve variant"})
			return
		}
		if existingVariant.ID == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Variant not found"})
			return
		}

		product, err := repository.GetProductByID(db, existingVariant.ProductID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve variant"})
			return
		}
		if product.AdminID != adminID {
			c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to delete this variant"})
			return
		}

		err = repository.DeleteVariant(db, variantUUID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete variant"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Variant deleted successfully"})
	})
}
