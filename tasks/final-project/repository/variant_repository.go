package repository

import (
	"database/sql"
	"final-project/domain"
)

// Function to get all variants with pagination
func GetAllVariants(db *sql.DB, limit int, offset int, search string) ([]domain.Variant, error) {
	query := "SELECT id, uuid, variant_name, quantity, product_id, created_at, updated_at FROM variants WHERE variant_name LIKE ? LIMIT ? OFFSET ?"
	searchTerm := "%" + search + "%"

	rows, err := db.Query(query, searchTerm, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var variants []domain.Variant
	for rows.Next() {
		var variant domain.Variant
		if err := rows.Scan(&variant.ID, &variant.UUID, &variant.VariantName, &variant.Quantity, &variant.ProductID, &variant.CreatedAt, &variant.UpdatedAt); err != nil {
			return nil, err
		}
		variants = append(variants, variant)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return variants, nil
}

// Function to get a variant by UUID
func GetVariantByUUID(db *sql.DB, uuid string) (domain.Variant, error) {
	var variant domain.Variant
	query := "SELECT id, uuid, variant_name, quantity, product_id, created_at, updated_at FROM variants WHERE uuid = ?"
	err := db.QueryRow(query, uuid).Scan(&variant.ID, &variant.UUID, &variant.VariantName, &variant.Quantity, &variant.ProductID, &variant.CreatedAt, &variant.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return variant, nil
		}
		return variant, err
	}
	return variant, nil
}

// Function to insert a new variant into the database and return the full variant data
func InsertVariant(db *sql.DB, variant *domain.Variant) error {
	query := "INSERT INTO variants (uuid, variant_name, quantity, product_id, created_at, updated_at) VALUES (?, ?, ?, ?, NOW(), NOW())"
	result, err := db.Exec(query, variant.UUID, variant.VariantName, variant.Quantity, variant.ProductID)
	if err != nil {
		return err
	}

	// Get the last inserted ID
	variantID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	return db.QueryRow("SELECT id, uuid, variant_name, quantity, product_id, created_at, updated_at FROM variants WHERE id = ?", variantID).
		Scan(&variant.ID, &variant.UUID, &variant.VariantName, &variant.Quantity, &variant.ProductID, &variant.CreatedAt, &variant.UpdatedAt)
}

// Function to update an existing variant in the database
func UpdateVariant(db *sql.DB, uuid string, variant *domain.Variant) error {
	query := `UPDATE variants SET variant_name = ?, quantity = ?, updated_at = NOW() WHERE uuid = ? AND product_id = ?`

	// Update the variant in the database
	_, err := db.Exec(query, variant.VariantName, variant.Quantity, uuid, variant.ProductID)
	if err != nil {
		return err
	}

	return db.QueryRow("SELECT id, uuid, variant_name, quantity, product_id, created_at, updated_at FROM variants WHERE uuid = ?", uuid).
		Scan(&variant.ID, &variant.UUID, &variant.VariantName, &variant.Quantity, &variant.ProductID, &variant.CreatedAt, &variant.UpdatedAt)
}

// Function to delete a variant from the database
func DeleteVariant(db *sql.DB, uuid string) error {
	query := "DELETE FROM variants WHERE uuid = ?"
	_, err := db.Exec(query, uuid)
	return err
}
