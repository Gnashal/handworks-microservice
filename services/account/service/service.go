package service

import (
	"context"
	"fmt"
	"handworks/common/grpc/account"
	types "handworks/common/types/account"
	"handworks/common/utils"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
)

type AccountService struct {
	L  *utils.Logger
	DB *pgxpool.Pool
	NC *nats.Conn
	account.UnimplementedAccountServiceServer
}

func (a *AccountService) SignUpCustomer(ctx context.Context, in *account.SignUpCustomerRequest) (*account.SignUpCustomerResponse, error) {
	var dbAcc types.DbAccount
	var customerID string

	if err := a.withTx(ctx, a.DB, func(tx pgx.Tx) error {
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
	}); err != nil {
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

	if err := a.withTx(ctx, a.DB, func(tx pgx.Tx) error {
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
	}); err != nil {
		return nil, err
	}
	employee := dbEmp.ToProto()
	employee.Account = dbAcc.ToProto()
	return &account.SignUpEmployeeResponse{
		Employee: employee,
	}, nil
}
func (a *AccountService) GetCustomer(ctx context.Context, in *account.GetCustomerRequest) (*account.GetCustomerResponse, error) {
	var dbAcc types.DbAccount
	var dbCustomer types.DbCustomer

	if err := a.withTx(ctx, a.DB, func(tx pgx.Tx) error {
		customer, acc, err := a.FetchCustomerData(ctx, tx, in.Id)
		if err != nil {
			return err
		}
		dbAcc = *acc
		dbCustomer = *customer
		return nil
	}); err != nil {
		return nil, err
	}
	customer := dbCustomer.ToProto()
	customer.Account = dbAcc.ToProto()

	return &account.GetCustomerResponse{
		Customer: customer,
	}, nil
}
func (a *AccountService) GetEmployee(ctx context.Context, in *account.GetEmployeeRequest) (*account.GetEmployeeResponse, error) {
	var dbAcc types.DbAccount
	var dbEmp types.DbEmployee

	if err := a.withTx(ctx, a.DB, func(tx pgx.Tx) error {
		emp, acc, err := a.FetchEmployeeData(ctx, tx, in.Id)
		if err != nil {
			return err
		}
		dbAcc = *acc
		dbEmp = *emp
		return nil
	}); err != nil {
		return nil, err
	}
	employee := dbEmp.ToProto()
	employee.Account = dbAcc.ToProto()
	return &account.GetEmployeeResponse{
		Employee: employee,
	}, nil
}
func (a *AccountService) UpdateCustomer(ctx context.Context, in *account.UpdateCustomerRequest) (*account.UpdateCustomerResponse, error) {
	var dbAcc types.DbAccount
	var dbCust types.DbCustomer

	if err := a.withTx(ctx, a.DB, func(tx pgx.Tx) error {
		acc, err := a.UpdateAccount(ctx, tx, in.Id, in.FirstName, in.LastName, in.Email)
		if err != nil {
			return err
		}
		dbAcc = *acc

		cust, err := a.FetchOnlyCustomerData(ctx, tx, in.Id)
		if err != nil {
			return err
		}
		dbCust = *cust

		return nil
	}); err != nil {
		return nil, fmt.Errorf("could not update customer: %w", err)
	}

	customer := dbCust.ToProto()
	customer.Account = dbAcc.ToProto()

	return &account.UpdateCustomerResponse{
		Customer: customer,
	}, nil
}

func (a *AccountService) UpdateEmployee(ctx context.Context, in *account.UpdateEmployeeRequest) (*account.UpdateEmployeeResponse, error) {
	var dbAcc types.DbAccount
	var dbEmp types.DbEmployee
	if err := a.withTx(ctx, a.DB, func(tx pgx.Tx) error {
		acc, err := a.UpdateAccount(ctx, tx, in.Id, in.FirstName, in.LastName, in.Email)
		if err != nil {
			return err
		}
		dbAcc = *acc

		emp, err := a.FetchOnlyEmployeeData(ctx, tx, dbAcc.ID)
		if err != nil {
			return err
		}
		dbEmp = *emp
		return nil
	}); err != nil {
		return nil, fmt.Errorf("could not query employee table")
	}

	employee := dbEmp.ToProto()
	employee.Account = dbAcc.ToProto()
	return &account.UpdateEmployeeResponse{
		Employee: employee,
	}, nil
}

func (a *AccountService) UpdateEmployeePerformanceScore(ctx context.Context, in *account.UpdatePerformanceScoreRequest) (*account.UpdatePerformanceScoreResponse, error) {
	var dbAcc types.DbAccount
	var dbEmp types.DbEmployee
	if err := a.withTx(ctx, a.DB, func(tx pgx.Tx) error {
		err := a.AddPerformanceScore(ctx, tx, in.NewPerformanceScore, in.Id)
		if err != nil {
			return err
		}

		emp, acc, err := a.FetchEmployeeData(ctx, tx, in.Id)
		if err != nil {
			return err
		}
		dbEmp = *emp
		dbAcc = *acc
		return nil
	}); err != nil {
		return nil, fmt.Errorf("could not update performance score in employee table: %w", err)
	}
	employee := dbEmp.ToProto()
	employee.Account = dbAcc.ToProto()
	return &account.UpdatePerformanceScoreResponse{
		Employee: employee,
	}, nil
}
func (a *AccountService) UpdateEmployeeStatus(ctx context.Context, in *account.UpdateEmployeeStatusRequest) (*account.UpdateEmployeeStatusResponse, error) {
	var dbAcc types.DbAccount
	var dbEmp types.DbEmployee
	if err := a.withTx(ctx, a.DB, func(tx pgx.Tx) error {
		err := a.UpdateStatus(ctx, tx, in.Status, in.Id)
		if err != nil {
			return err
		}

		emp, acc, err := a.FetchEmployeeData(ctx, tx, in.Id)
		if err != nil {
			return err
		}
		dbEmp = *emp
		dbAcc = *acc
		return nil
	}); err != nil {
		return nil, fmt.Errorf("could not update performance score in employee table: %w", err)
	}
	employee := dbEmp.ToProto()
	employee.Account = dbAcc.ToProto()
	return &account.UpdateEmployeeStatusResponse{
		Employee: employee,
	}, nil
}
func (a *AccountService) DeleteEmployee(ctx context.Context, in *account.DeleteEmployeeRequest) (*account.DeleteEmployeeResponse, error) {
	var dbAcc types.DbAccount
	var dbEmp types.DbEmployee

	if err := a.withTx(ctx, a.DB, func(tx pgx.Tx) error {
		emp, acc, err := a.DeleteEmployeeData(ctx, tx, in.EmpId, in.Id)
		if err != nil {
			return err
		}
		dbEmp = *emp
		dbAcc = *acc
		return nil
	}); err != nil {
		return nil, err
	}

	employee := dbEmp.ToProto()
	employee.Account = dbAcc.ToProto()

	return &account.DeleteEmployeeResponse{
		Ok:       true,
		Message:  "deleting employee success",
		Employee: employee,
	}, nil
}
func (a *AccountService) DeleteCustomer(ctx context.Context, in *account.DeleteCustomerRequest) (*account.DeleteCustomerResponse, error) {
	var dbCust types.DbCustomer
	var dbAcc types.DbAccount

	if err := a.withTx(ctx, a.DB, func(tx pgx.Tx) error {
		cust, acc, err := a.DeleteCustomerData(ctx, tx, in.CustId, in.Id)
		if err != nil {
			return err
		}
		dbCust = *cust
		dbAcc = *acc
		return nil
	}); err != nil {
		return nil, err
	}
	customer := dbCust.ToProto()
	customer.Account = dbAcc.ToProto()
	return &account.DeleteCustomerResponse{
		Ok:       true,
		Message:  "deleting customer success",
		Customer: customer,
	}, nil
}
