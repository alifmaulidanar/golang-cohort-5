package repository

import (
	"database/sql"
	"final-project/domain"
)

func InsertAdmin(db *sql.DB, admin domain.Admin) error {
	query := "INSERT INTO admins (uuid, name, email, password, created_at, updated_at) VALUES (?, ?, ?, ?, NOW(), NOW())"
	_, err := db.Exec(query, admin.UUID, admin.Name, admin.Email, admin.Password)
	return err
}

func FindAdminByEmail(db *sql.DB, email string) (domain.Admin, error) {
	var admin domain.Admin
	query := "SELECT id, uuid, name, email, password FROM admins WHERE email = ?"
	err := db.QueryRow(query, email).Scan(&admin.ID, &admin.UUID, &admin.Name, &admin.Email, &admin.Password)
	return admin, err
}
