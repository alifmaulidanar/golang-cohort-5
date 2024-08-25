package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

const (
	host     = "localhost"
	port     = 3306
	user     = "root"
	password = "root"
	dbname   = "go_sql_sesi_6"
)

var (
	db  *sql.DB
	err error
)

type Variants struct {
	ID           int
	Variant_name string
	Quantity     int
	Product_id   int
	Created_at   string
	Updated_at   string
}

type Products struct {
	ID         int
	Name       string
	Created_at string
	Updated_at string
}

func main() {
	mysqlInfo := fmt.Sprintf("%s:%s@/%s", user, password, dbname)

	db, err = sql.Open("mysql", mysqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected to database")

	// createProduct()
	// updateProduct()
	// getProductById(2)
	// createVariant()
	// updateVariantById(3)
	// deleteVariantById(4)
	// getProductWithVariant(3)
}

func createProduct() {
	var product = Products{}

	query := `INSERT INTO Products (Name, Created_at, Updated_at) VALUES (?, NOW(), NOW())`
	result, err := db.Exec(query, "Product D")
	if err != nil {
		panic(err)
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		panic(err)
	}

	sqlSelect := `SELECT * FROM Products WHERE ID = ?`
	err = db.QueryRow(sqlSelect, lastInsertID).Scan(&product.ID, &product.Name, &product.Created_at, &product.Updated_at)
	if err != nil {
		panic(err)
	}

	fmt.Println("Product created successfully")
	fmt.Printf("Product: %+v\n", product)
}

func updateProduct() {
	var product = Products{}
	query := `UPDATE Products SET Name = ?, Updated_at = NOW() WHERE ID = ?`
	_, err := db.Exec(query, "Product B V2", 2)
	if err != nil {
		panic(err)
	}

	sqlSelect := `SELECT * FROM Products WHERE ID = ?`
	err = db.QueryRow(sqlSelect, 2).Scan(&product.ID, &product.Name, &product.Created_at, &product.Updated_at)
	if err != nil {
		panic(err)
	}

	fmt.Println("Product updated successfully")
	fmt.Printf("Updated Product: %+v\n", product)
}

func getProductById(id int) {
	var product Products
	query := `SELECT ID, Name, Created_at, Updated_at FROM Products WHERE ID = ?`
	err := db.QueryRow(query, id).Scan(&product.ID, &product.Name, &product.Created_at, &product.Updated_at)
	if err != nil {
		panic(err)
	}

	fmt.Println("Product retrieved successfully")
	fmt.Printf("Product: %+v\n", product)
}

func createVariant() {
	var variant = Variants{}

	query := `INSERT INTO Variants (Variant_name, Quantity, Product_id, Created_at, Updated_at) VALUES (?, ?, ?, NOW(), NOW())`
	result, err := db.Exec(query, "Variant D", 500, 4)
	if err != nil {
		panic(err)
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		panic(err)
	}

	sqlSelect := `SELECT * FROM Variants WHERE ID = ?`
	err = db.QueryRow(sqlSelect, lastInsertID).Scan(&variant.ID, &variant.Variant_name, &variant.Quantity, &variant.Product_id, &variant.Created_at, &variant.Updated_at)
	if err != nil {
		panic(err)
	}

	fmt.Println("Variant created successfully")
	fmt.Printf("Variant: %+v\n", variant)
}

func updateVariantById(id int) {
	var variant = Variants{}
	query := `UPDATE Variants SET Variant_name = ?, Quantity = ?, Updated_at = NOW() WHERE ID = ?`
	_, err := db.Exec(query, "Variant C (updated)", 200, id)
	if err != nil {
		panic(err)
	}

	sqlSelect := `SELECT * FROM Variants WHERE ID = ?`
	err = db.QueryRow(sqlSelect, id).Scan(&variant.ID, &variant.Variant_name, &variant.Quantity, &variant.Product_id, &variant.Created_at, &variant.Updated_at)
	if err != nil {
		panic(err)
	}

	fmt.Println("Variant updated successfully")
	fmt.Printf("Updated Variant: %+v\n", variant)
}

func deleteVariantById(id int) {
	query := `DELETE FROM Variants WHERE ID = ?`
	_, err := db.Exec(query, id)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Variant with ID %d deleted successfully\n", id)
}

func getProductWithVariant(productId int) {
	var product Products
	query := `SELECT ID, Name, Created_at, Updated_at FROM Products WHERE ID = ?`
	err := db.QueryRow(query, productId).Scan(&product.ID, &product.Name, &product.Created_at, &product.Updated_at)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Product: %+v\n", product)

	rows, err := db.Query(`SELECT ID, Variant_name, Quantity, Created_at, Updated_at FROM Variants WHERE Product_id = ?`, productId)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var variant Variants
		err := rows.Scan(&variant.ID, &variant.Variant_name, &variant.Quantity, &variant.Created_at, &variant.Updated_at)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Variant: %+v\n", variant)
	}
	if err = rows.Err(); err != nil {
		panic(err)
	}

	fmt.Println("Product with variants retrieved successfully")
}
