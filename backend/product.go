package main

type Product struct {
	SKU         int     `json:"sku"`
	Title       string  `json:"title"`
	Img         string  `json:"img"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}
