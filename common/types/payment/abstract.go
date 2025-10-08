package types

import "time"

type DbQuote struct {
	ID         string    `db:"id"`
	CustomerID string    `db:"customer_id"`
	ServiceID  string    `db:"service_id"`
	Subtotal   float32   `db:"subtotal"`
	AddonTotal float32   `db:"addon_total"`
	TotalPrice float32   `db:"total_price"`
	IsValid    bool      `db:"is_valid"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
	Addons     []DbQuoteAddon
}

type DbQuoteAddon struct {
	ID         string    `db:"id"`
	QuoteID    string    `db:"quote_id"`
	AddonID    string    `db:"addon_id"`
	AddonName  string    `db:"addon_name"`
	AddonPrice float32   `db:"addon_price"`
	CreatedAt  time.Time `db:"created_at"`
}
type QuoteResponse struct {
	Addons []DbQuoteAddon
}
