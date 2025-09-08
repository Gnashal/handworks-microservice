package service

import (
	"context"
	"fmt"
	"handworks-services-account/types"
	"handworks/common/grpc/account"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func (a *AccountService) withTx(ctx context.Context, db *pgxpool.Pool, fn func(pgx.Tx) error) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := fn(tx); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (a *AccountService) FetchCustomerData(c context.Context, tx pgx.Tx, ID string) (*types.DbCustomer, *types.DbAccount, error) {
	var acc types.DbAccount
	var customer types.DbCustomer

	if err := tx.QueryRow(c,
		`SELECT id, account_id FROM account.customers WHERE id = $1`,
		ID,
	).Scan(&customer.ID, &customer.AccountID); err != nil {
		return nil, nil, fmt.Errorf("could not query customer table: %w", err)
	}

	if err := tx.QueryRow(c,
		`SELECT first_name, last_name, email, provider, clerk_id, role, id, created_at, updated_at
		 FROM account.accounts WHERE id = $1`,
		customer.AccountID,
	).Scan(
		&acc.FirstName,
		&acc.LastName,
		&acc.Email,
		&acc.Provider,
		&acc.ClerkID,
		&acc.Role,
		&acc.ID,
		&acc.CreatedAt,
		&acc.UpdatedAt,
	); err != nil {
		return nil, nil, fmt.Errorf("could not query accounts table: %w", err)
	}

	return &customer, &acc, nil
}
func (a *AccountService) FetchEmployeeData(c context.Context, tx pgx.Tx, ID string) (*types.DbEmployee, *types.DbAccount, error) {
	var acc types.DbAccount
	var emp types.DbEmployee

	if err := tx.QueryRow(c,
		`SELECT id, account_id, position, status, performance_score, hire_date
		 FROM account.employees WHERE id = $1`,
		ID,
	).Scan(
		&emp.ID,
		&emp.AccountID,
		&emp.Position,
		&emp.Status,
		&emp.Performance,
		&emp.HireDate,
	); err != nil {
		return nil, nil, fmt.Errorf("could not query employees table: %w", err)
	}

	if err := tx.QueryRow(c,
		`SELECT first_name, last_name, email, provider, clerk_id, role, id, created_at, updated_at
		 FROM account.accounts WHERE id = $1`,
		emp.AccountID,
	).Scan(
		&acc.FirstName,
		&acc.LastName,
		&acc.Email,
		&acc.Provider,
		&acc.ClerkID,
		&acc.Role,
		&acc.ID,
		&acc.CreatedAt,
		&acc.UpdatedAt,
	); err != nil {
		return nil, nil, fmt.Errorf("could not query accounts table: %w", err)
	}

	return &emp, &acc, nil
}

func (a *AccountService) CreateAccount(c context.Context, tx pgx.Tx, FirstName, LastName, Email, Provider, ClerkId, Role string) (*types.DbAccount, error) {
	var acc types.DbAccount
	err := tx.QueryRow(c,
		`INSERT INTO account.accounts (first_name, last_name, email, provider, clerk_id, role)
		VALUES ($1,$2, $3, $4, $5, $6)
		RETURNING first_name, last_name, email, provider, clerk_id, role, id, created_at, updated_at`,
		FirstName, LastName, Email, Provider, ClerkId, Role,
	).Scan(
		&acc.FirstName,
		&acc.LastName,
		&acc.Email,
		&acc.Provider,
		&acc.ClerkID,
		&acc.Role,
		&acc.ID,
		&acc.CreatedAt,
		&acc.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("could not create account: %w", err)
	}
	return &acc, nil
}

func (a *AccountService) CreateEmployee(c context.Context, tx pgx.Tx, id, position string, hireDate time.Time) (*types.DbEmployee, error) {
	var emp types.DbEmployee
	if err := tx.QueryRow(c,
		`INSERT INTO account.employee (account_id, position, status, performance_score, hire_date)
	VALUES ($1, $2, $3, $4, $5) 
	RETURNING id, position, status, performance_score, hire_date`, id, position, account.EmployeeStatus_INACTIVE.String(), 5.0, hireDate).Scan(
		&emp.ID, &emp.Position, &emp.Status, &emp.Performance,
		&emp.HireDate); err != nil {
		return nil, fmt.Errorf("could not insert into employee table: %w", err)
	}
	return &emp, nil
}

func (a *AccountService) CreateCustomer(c context.Context, tx pgx.Tx, id string) (string, error) {
	var customerId string
	if err := tx.QueryRow(c,
		`INSERT INTO account.customers (account_id)
		VALUES ($1)
		RETURNING id`, id).Scan(&customerId); err != nil {
		return "", fmt.Errorf("could not insert into customer table: %w", err)
	}
	return customerId, nil
}
