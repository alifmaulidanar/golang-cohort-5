package repository

import (
	"database/sql"
	"final-project/domain"
)

// Function to get all products with pagination
func GetAllProducts(db *sql.DB, limit int, offset int, search string) ([]domain.Product, error) {
	query := "SELECT id, uuid, name, image_url, admin_id, created_at, updated_at FROM products WHERE name LIKE ? LIMIT ? OFFSET ?"
	searchTerm := "%" + search + "%"

	rows, err := db.Query(query, searchTerm, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		var product domain.Product
		if err := rows.Scan(&product.ID, &product.UUID, &product.Name, &product.ImageURL, &product.AdminID, &product.CreatedAt, &product.UpdatedAt); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return products, nil
}

// Function to get a product by UUID
func GetProductByUUID(db *sql.DB, uuid string) (domain.Product, error) {
	var product domain.Product
	query := "SELECT id, uuid, name, image_url, admin_id, created_at, updated_at FROM products WHERE uuid = ?"

	err := db.QueryRow(query, uuid).Scan(&product.ID, &product.UUID, &product.Name, &product.ImageURL, &product.AdminID, &product.CreatedAt, &product.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return product, nil
		}
		return product, err
	}
	return product, nil
}

// Function to get a product by ID
func GetProductByID(db *sql.DB, id int) (domain.Product, error) {
	var product domain.Product
	query := "SELECT id, uuid, name, image_url, admin_id, created_at, updated_at FROM products WHERE id = ?"

	err := db.QueryRow(query, id).Scan(&product.ID, &product.UUID, &product.Name, &product.ImageURL, &product.AdminID, &product.CreatedAt, &product.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return product, nil
		}
		return product, err
	}
	return product, nil
}

// Function to insert a new product into the database and return the full product data
func InsertProduct(db *sql.DB, product *domain.Product) error {
	query := "INSERT INTO products (uuid, name, image_url, admin_id, created_at, updated_at) VALUES (?, ?, ?, ?, NOW(), NOW())"

	result, err := db.Exec(query, product.UUID, product.Name, product.ImageURL, product.AdminID)
	if err != nil {
		return err
	}

	// Get the last inserted ID
	productID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	return db.QueryRow("SELECT id, uuid, name, image_url, admin_id, created_at, updated_at FROM products WHERE id = ?", productID).
		Scan(&product.ID, &product.UUID, &product.Name, &product.ImageURL, &product.AdminID, &product.CreatedAt, &product.UpdatedAt)
}

// Function to update an existing product in the database
func UpdateProduct(db *sql.DB, uuid string, product *domain.Product) error {
	query := `UPDATE products SET name = ?, image_url = ?, updated_at = NOW() WHERE uuid = ? AND admin_id = ?`
	_, err := db.Exec(query, product.Name, product.ImageURL, uuid, product.AdminID)
	if err != nil {
		return err
	}

	return db.QueryRow("SELECT id, uuid, name, image_url, admin_id, created_at, updated_at FROM products WHERE uuid = ?", uuid).
		Scan(&product.ID, &product.UUID, &product.Name, &product.ImageURL, &product.AdminID, &product.CreatedAt, &product.UpdatedAt)
}

// Function to delete a product from the database
func DeleteProduct(db *sql.DB, uuid string) error {
	query := "DELETE FROM products WHERE uuid = ?"
	_, err := db.Exec(query, uuid)
	return err
}
