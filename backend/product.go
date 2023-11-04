package main

type ProductImage struct {
	Original     string `json:"original"`
	Thumbnail    string `json:"thumbnail"`
	LargeProduct string `json:"largeProduct"`
	Zoom         string `json:"zoom"`
}

type Product struct {
	Title  string         `json:"title"`
	SKU    int64          `json:"sku"`
	Images []ProductImage `json:"images"`
}
