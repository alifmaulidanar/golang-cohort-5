package controllers

import (
	"errors"
	"net/http"
	"sesi-7-assignment/database"
	"sesi-7-assignment/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// create order
func CreateOrder(c *gin.Context) {
	var input struct {
		OrderedAt    string        `json:"orderedAt"`
		CustomerName string        `json:"customerName"`
		Items        []models.Item `json:"items"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// datetime menggunakan acuan format RFC3339
	orderedAt, err := time.Parse(time.RFC3339, input.OrderedAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}

	order := models.Order{CustomerName: input.CustomerName, OrderedAt: orderedAt, Items: input.Items}
	db := database.GetDB()

	err = db.Create(&order).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": order})
}

// get all orders
func GetAllOrders(c *gin.Context) {
	var orders []models.Order
	db := database.GetDB()

	err := db.Preload("Items").Find(&orders).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting orders"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": orders})
}

// get order by ID
func GetOrderById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var order models.Order
	db := database.GetDB()
	err = db.Preload("Items").First(&order, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error finding order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": order})
}

// update order
func UpdateOrder(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var input struct {
		OrderedAt    string        `json:"orderedAt"`
		CustomerName string        `json:"customerName"`
		Items        []models.Item `json:"items"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var order models.Order
	db := database.GetDB()
	err = db.Preload("Items").First(&order, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error finding order"})
		return
	}

	orderedAt, err := time.Parse(time.RFC3339, input.OrderedAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}

	err = db.Model(&order).Updates(models.Order{CustomerName: input.CustomerName, OrderedAt: orderedAt}).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating order"})
		return
	}

	err = db.Model(&order).Association("Items").Replace(input.Items)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating items"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": order})
}

// delete order
func DeleteOrder(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var order models.Order
	db := database.GetDB()
	err = db.Preload("Items").First(&order, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error finding order"})
		return
	}

	// harus hapus item yg berelasi dengan order dahulu
	err = db.Where("order_id = ?", order.ID).Delete(&models.Item{}).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting items"})
		return
	}

	// sekarang dapat menghapus order
	err = db.Delete(&order).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": true})
}
