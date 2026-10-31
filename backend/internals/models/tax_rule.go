package models

import "time"

// TaxRule is an admin-defined tax rate applied to a subscription and
// snapshotted onto each invoice.
type TaxRule struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	RateBps   int       `json:"rateBps"` // basis points, 1800 = 18%
	Active    bool      `json:"active"`
	CreatedBy string    `json:"createdBy"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
