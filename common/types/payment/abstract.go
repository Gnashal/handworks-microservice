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

type AddonCleaningPrice struct {
	AddonName  string
	AddonPrice float32
}
type CleaningPrices struct {
	MainServicePrice float32
	AddonPrices      []AddonCleaningPrice
}
type QuoteResponse struct {
	Addons []DbQuoteAddon
}

type BookingReply struct {
	Source string         `json:"source"`
	Prices CleaningPrices `json:"prices,omitempty"`
	Error  string         `json:"error,omitempty"`
}
