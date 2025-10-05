package types

import (
	"handworks/common/grpc/account"
	"time"
)

type DbCustomer struct {
	ID        string
	AccountID string
}

func (dbCustomer DbCustomer) ToProto() *account.Customer {
	return &account.Customer{
		Id:      dbCustomer.ID,
		Account: nil,
	}
}

type DbEmployee struct {
	ID          string
	AccountID   string
	Position    string
	Status      string
	Performance float32
	HireDate    time.Time
	NumRatings  int32
}

type DbAccount struct {
	ID        string
	FirstName string
	LastName  string
	Email     string
	Provider  string
	ClerkID   string
	Role      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CleanerAssigned struct {
	ID               string
	CleanerFirstName string
	CleanerLastName  string
	PFPUrl           string
}
type BookingReply struct {
	Source   string            `json:"source"`
	Cleaners []CleanerAssigned `json:"cleaners,omitempty"`
	Error    string            `json:"error,omitempty"`
}
