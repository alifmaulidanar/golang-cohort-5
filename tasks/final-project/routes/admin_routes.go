package routes

import (
	"database/sql"
	"net/http"

	"final-project/domain"
	"final-project/repository"
	"final-project/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func AdminRoutes(r *gin.Engine, db *sql.DB) {
	validate := validator.New()
	validate.RegisterValidation("password", service.PasswordValidator)

	// Register Admin
	r.POST("/auth/register", func(c *gin.Context) {
		var admin domain.Admin
		if err := c.ShouldBind(&admin); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Input validation
		if err := validate.Struct(admin); err != nil {
			validationErrors := err.(validator.ValidationErrors)
			errors := make(map[string]string)
			for _, fieldError := range validationErrors {
				errors[fieldError.Field()] = fieldError.Error()
			}
			c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
			return
		}

		// Hash the password
		admin.UUID = uuid.New().String()
		hashedPassword, err := service.HashPassword(admin.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		admin.Password = hashedPassword

		// Insert the admin into the database
		if err := repository.InsertAdmin(db, admin); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register admin"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Admin registered successfully", "uuid": admin.UUID})
	})

	// Login Admin
	r.POST("/auth/login", func(c *gin.Context) {
		var loginData domain.LoginRequest
		if err := c.ShouldBind(&loginData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Input validation
		if err := validate.Struct(loginData); err != nil {
			validationErrors := err.(validator.ValidationErrors)
			errors := make(map[string]string)
			for _, fieldError := range validationErrors {
				errors[fieldError.Field()] = fieldError.Error()
			}
			c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
			return
		}

		// Check if the email exists
		admin, err := repository.FindAdminByEmail(db, loginData.Email)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		// Check if the password is correct
		if !service.CheckPasswordHash(loginData.Password, admin.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		// Generate JWT
		token, err := service.GenerateJWT(admin)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": token})
	})
}
