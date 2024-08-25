CREATE DATABASE go_sql_sesi_6;

USE go_sql_sesi_6;

CREATE TABLE products (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

CREATE TABLE variants (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    variant_name VARCHAR(50) NOT NULL,
    quantity INT NOT NULL,
    product_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY (product_id) REFERENCES products(id)
);

SHOW TABLES;
SELECT * FROM products;
SELECT * FROM variants;