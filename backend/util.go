package main

import (
	"bytes"
	"database/sql"
	"encoding/base64"
)

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

func base64Hash(h [32]byte) string {
	buffer := new(bytes.Buffer)
	encoder := base64.NewEncoder(base64.StdEncoding, buffer)

	_, _ = encoder.Write(h[:])
	return string(buffer.Bytes())
}
