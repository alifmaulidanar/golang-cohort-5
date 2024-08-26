package models

import "time"

// strukutur data Order
type Order struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	CustomerName string    `json:"customerName"`
	OrderedAt    time.Time `json:"orderedAt"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	Items        []Item    `json:"items" gorm:"foreignKey:OrderID"`
}

// strukutur data Item
type Item struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Quantity    uint      `json:"quantity"`
	OrderID     uint      `json:"orderId"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
