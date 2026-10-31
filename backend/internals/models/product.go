package models

import "time"

// Product is a billable product. Its tiers live in ProductVariant.
type Product struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Active      bool      `json:"active"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// ProductVariant is one tier of a product, e.g. Basic or Pro.
type ProductVariant struct {
	ID              int64     `json:"id"`
	ProductID       int64     `json:"productId"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	PriceCents      int64     `json:"priceCents"`
	Currency        string    `json:"currency"`
	BillingInterval string    `json:"billingInterval"` // month | year
	Active          bool      `json:"active"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}
