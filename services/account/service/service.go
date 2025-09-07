package service

import (
	"context"
	"handworks-services-account/types"
	"handworks/common/grpc/account"
	"handworks/common/utils"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AccountService struct {
	L  *utils.Logger
	DB *pgxpool.Pool
	account.UnimplementedAccountServiceServer
}

func (a *AccountService) SignUpCustomer(ctx context.Context, in *account.SignUpCustomerRequest) (*account.SignUpCustomerResponse, error) {
	var dbAcc types.DbAccount
	var customerID string

	err := a.withTx(ctx, a.DB, func(tx pgx.Tx) error {
		acc, err := a.CreateAccount(ctx, tx, in.FirstName, in.LastName, in.Email, in.Provider, in.ClerkId, in.Role)
		if err != nil {
			return err
		}
		dbAcc = *acc

		id, err := a.CreateCustomer(ctx, tx, dbAcc.ID)
		if err != nil {
			return err
		}
		customerID = id
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &account.SignUpCustomerResponse{
		Customer: &account.Customer{
			Id:      customerID,
			Account: dbAcc.ToProto(),
		},
	}, nil
}

func (a *AccountService) SignUpEmployee(ctx context.Context, in *account.SignUpEmployeeRequest) (*account.SignUpEmployeeResponse, error) {
	var dbAcc types.DbAccount
	var dbEmp types.DbEmployee

	err := a.withTx(ctx, a.DB, func(tx pgx.Tx) error {
		acc, err := a.CreateAccount(ctx, tx, in.FirstName, in.LastName, in.Email, in.Provider, in.ClerkId, in.Role)
		if err != nil {
			return err
		}
		dbAcc = *acc
		emp, err := a.CreateEmployee(ctx, tx, dbAcc.ID, in.Position, in.HireDate.AsTime())
		if err != nil {
			return err
		}
		dbEmp = *emp

		return nil
	})
	if err != nil {
		return nil, err
	}
	employee := dbEmp.ToProto()
	employee.Account = dbAcc.ToProto()
	return &account.SignUpEmployeeResponse{
		Employee: employee,
	}, nil
}
