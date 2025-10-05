package types

import (
	"handworks/common/grpc/account"

	"google.golang.org/protobuf/types/known/timestamppb"
)

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

func (dbEmp DbEmployee) ToProto() *account.Employee {
	return &account.Employee{
		Id:               dbEmp.ID,
		Account:          nil, // can join and set later
		Position:         dbEmp.Position,
		Status:           account.EmployeeStatus(account.EmployeeStatus_value[dbEmp.Status]),
		PerformanceScore: dbEmp.Performance,
		HireDate:         timestamppb.New(dbEmp.HireDate),
		NumRatings:       dbEmp.NumRatings,
	}
}
