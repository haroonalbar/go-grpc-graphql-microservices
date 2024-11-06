package order

import "time"

type Order struct {
	ID         string           `json:"id"`
	AccountID  string           `json:"account_id"`
	CreatedAt  time.Time        `json:"created_at"`
	TotalPrice float64          `json:"tatal_price"`
	Products   []OrderedProduct `json:"products"`
}

type OrderedProduct struct {
	ID          string  `json:"product_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    uint32  `json:"quantity"`
}
