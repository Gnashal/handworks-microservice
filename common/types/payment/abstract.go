package types

import (
	"encoding/json"
	"time"
)

type DbQuote struct {
	ID          string
	CustomerID  string
	MainService string
	Subtotal    float32
	AddonTotal  float32
	TotalPrice  float32
	IsValid     bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Addons      []*DbQuoteAddon
}

type DbQuoteAddon struct {
	ID            string
	QuoteID       string
	ServiceType   string
	ServiceDetail json.RawMessage // serialized ServicesRequest
	AddonPrice    float32
	CreatedAt     time.Time
}

type QuoteResponse struct {
	Addons []DbQuoteAddon
}
