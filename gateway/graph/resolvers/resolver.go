package resolvers

import (
	generated "handworks-gateway/graph/generated/models"
	"handworks-gateway/graph/resolvers/helpers"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	accounts  []*generated.Account
	customers []*generated.Customer
	employees []*generated.Employee
	admins    []*generated.Admin
	watchList *generated.WatchList
	helpers   *helpers.Helpers
}
