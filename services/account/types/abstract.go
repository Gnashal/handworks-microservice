package types

import (
	"handworks/common/grpc/account"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
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
}

func (dbEmp DbEmployee) ToProto() *account.Employee {
	return &account.Employee{
		Id:               dbEmp.ID,
		Account:          nil, // can join and set later
		Position:         dbEmp.Position,
		Status:           account.EmployeeStatus(account.EmployeeStatus_value[dbEmp.Status]),
		PerformanceScore: dbEmp.Performance,
		HireDate:         timestamppb.New(dbEmp.HireDate),
	}
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

func (dbAcc DbAccount) ToProto() *account.Account {
	return &account.Account{
		Id:        dbAcc.ID,
		FirstName: dbAcc.FirstName,
		LastName:  dbAcc.LastName,
		Email:     dbAcc.Email,
		Provider:  dbAcc.Provider,
		ClerkId:   dbAcc.ClerkID,
		Role:      dbAcc.Role,
		CreatedAt: timestamppb.New(dbAcc.CreatedAt),
		UpdatedAt: timestamppb.New(dbAcc.UpdatedAt),
	}
}
