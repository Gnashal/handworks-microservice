package helpers

import (
	model "handworks-gateway/graph/generated/models"
	"handworks/common/grpc/account"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func MapAccount(a *account.Account) *model.Account {
	if a == nil {
		return nil
	}

	return &model.Account{
		ID:        a.Id,
		FirstName: a.FirstName,
		LastName:  a.LastName,
		Email:     a.Email,
		Provider:  &a.Provider,
		Role:      a.Role,
		ClerkID:   a.ClerkId,
		CreatedAt: timestampToTime(a.CreatedAt),
		UpdatedAt: timestampToTime(a.UpdatedAt),
	}
}

func MapCustomer(c *account.Customer) *model.Customer {
	if c == nil {
		return nil
	}

	return &model.Customer{
		ID:      c.Id,
		Account: MapAccount(c.Account),
	}
}

func MapEmployee(e *account.Employee) *model.Employee {
	if e == nil {
		return nil
	}

	return &model.Employee{
		ID:               e.Id,
		Account:          MapAccount(e.Account),
		Position:         e.Position,
		Status:           model.EmployeeStatus(e.Status.String()), // map enum -> GraphQL enum
		PerformanceScore: float64(e.PerformanceScore),
		HireDate:         timestampToTime(e.HireDate),
		NumRatings:       int32(e.NumRatings),
	}
}

// Convert protobuf Timestamp to Go time.Time
func timestampToTime(ts *timestamppb.Timestamp) time.Time {
	if ts == nil {
		return time.Time{}
	}
	return ts.AsTime()
}
