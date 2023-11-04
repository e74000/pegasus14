package main

import "database/sql"

func ParseRows(rows *sql.Rows) []Product {
	products := make([]Product, 0)

	for rows.Next() {
		product := *new(Product)

		err := rows.Scan(&product.SKU, &product.Title, &product.Img, &product.Description, &product.Price)
		if err != nil {
			continue
		}

		products = append(products, product)
	}

	return products
}
