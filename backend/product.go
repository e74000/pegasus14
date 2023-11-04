package main

type Product struct {
	SKU         int     `json:"sku"`
	Title       string  `json:"title"`
	Img         string  `json:"img"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type Impression struct {
	User         string  `json:"user"`
	Product      int     `json:"product"`
	Liked        int     `json:"liked"`
	ViewDuration float64 `json:"viewDuration"`
	Claim        Claim   `json:"claim"`
}
