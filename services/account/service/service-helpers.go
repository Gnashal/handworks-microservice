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

func (a *AccountService) withTx(
	ctx context.Context,
	db *pgxpool.Pool,
	fn func(pgx.Tx) error,
) (err error) {
	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				a.L.Error("rollback failed: %v", rbErr)
			}
		} else {
			err = tx.Commit(ctx)
		}
	}()
	return fn(tx)
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
func (a *AccountService) FetchOnlyCustomerData(c context.Context, tx pgx.Tx, accountId string) (*types.DbCustomer, error) {
	var cust types.DbCustomer
	err := tx.QueryRow(c,
		`SELECT id, account_id
		 FROM account.customers
		 WHERE account_id = $1`,
		accountId,
	).Scan(
		&cust.ID,
		&cust.AccountID,
	)
	if err != nil {
		return nil, fmt.Errorf("could not fetch customer: %w", err)
	}
	return &cust, nil
}
func (a *AccountService) FetchOnlyEmployeeData(c context.Context, tx pgx.Tx, accountId string) (*types.DbEmployee, error) {
	var emp types.DbEmployee
	err := tx.QueryRow(c,
		`SELECT id, account_id, position, status,performance_score, hire_date, num_ratings
		 FROM account.employees
		 WHERE account_id = $1`,
		accountId,
	).Scan(
		&emp.ID,
		&emp.AccountID,
		&emp.Position,
		&emp.Status,
		&emp.Performance,
		&emp.HireDate,
		&emp.NumRatings,
	)
	if err != nil {
		return nil, fmt.Errorf("could not fetch employee: %w", err)
	}
	return &emp, nil
}

func (a *AccountService) FetchEmployeeData(c context.Context, tx pgx.Tx, ID string) (*types.DbEmployee, *types.DbAccount, error) {
	var acc types.DbAccount
	var emp types.DbEmployee

	if err := tx.QueryRow(c,
		`SELECT id, account_id, position, status, performance_score, hire_date, num_ratings
		 FROM account.employees WHERE id = $1`,
		ID,
	).Scan(
		&emp.ID,
		&emp.AccountID,
		&emp.Position,
		&emp.Status,
		&emp.Performance,
		&emp.HireDate,
		&emp.NumRatings,
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
		`INSERT INTO account.employees (account_id, position, status, performance_score, hire_date, num_ratings)
	VALUES ($1, $2, $3, $4, $5, $6) 
	RETURNING id, position, status, performance_score, hire_date, num_ratings`, id, position, account.EmployeeStatus_INACTIVE.String(), 5.0, hireDate, 0).Scan(
		&emp.ID, &emp.Position, &emp.Status, &emp.Performance,
		&emp.HireDate, &emp.NumRatings); err != nil {
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

func (a *AccountService) UpdateAccount(c context.Context, tx pgx.Tx, id, firstName, lastName, email string) (*types.DbAccount, error) {
	var acc types.DbAccount
	err := tx.QueryRow(c,
		`UPDATE account.accounts
		 SET first_name = $1, last_name = $2, email = $3, updated_at = NOW()
		 WHERE id = $4
		 RETURNING id, first_name, last_name, email, role, provider, created_at, updated_at, clerk_id`,
		firstName, lastName, email, id,
	).Scan(
		&acc.ID,
		&acc.FirstName,
		&acc.LastName,
		&acc.Email,
		&acc.Role,
		&acc.Provider,
		&acc.CreatedAt,
		&acc.UpdatedAt,
		&acc.ClerkID,
	)
	if err != nil {
		return nil, fmt.Errorf("could not update account: %w", err)
	}
	return &acc, nil
}

func (a *AccountService) AddPerformanceScore(c context.Context, tx pgx.Tx, newScore float32, empId string) error {
	a.L.Debug("New Score: %.2f", newScore)
	a.L.Debug("UUID: %s", empId)
	args := pgx.NamedArgs{
		"newScore": newScore,
		"id":       empId,
	}
	cmdTag, err := tx.Exec(c,
		`UPDATE account.employees
	 SET performance_score = ((performance_score * num_ratings) + @newScore) / (num_ratings + 1),
	     num_ratings = num_ratings + 1, updated_at = NOW()
	 WHERE id = @id::uuid`,
		args,
	)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("no employee found with id %s", empId)
	}
	return nil
}
func (a *AccountService) UpdateStatus(c context.Context, tx pgx.Tx, status, empId string) error {
	newStatus := DetermineStatusType(status)
	args := pgx.NamedArgs{
		"newStatus": newStatus,
		"id":        empId,
	}
	_, err := tx.Exec(c,
		`UPDATE account.employees
	SET status = @newStatus
	WHERE id = @id`, args)
	if err != nil {
		return err
	}
	return nil
}
func (a *AccountService) DeleteEmployeeData(c context.Context, tx pgx.Tx, empId, accId string) (*types.DbEmployee, *types.DbAccount, error) {
	var emp types.DbEmployee
	var acc types.DbAccount

	// Delete employee, return row
	if err := tx.QueryRow(c,
		`DELETE FROM account.employees
		 WHERE id = $1
		 RETURNING id, account_id, position, status, performance_score, hire_date, num_ratings`,
		empId,
	).Scan(
		&emp.ID,
		&emp.AccountID,
		&emp.Position,
		&emp.Status,
		&emp.Performance,
		&emp.HireDate,
		&emp.NumRatings,
	); err != nil {
		return nil, nil, fmt.Errorf("could not delete employee: %w", err)
	}
	if err := tx.QueryRow(c,
		`DELETE FROM account.accounts
		 WHERE id = $1
		 RETURNING first_name, last_name, email, provider, clerk_id, role, id, created_at, updated_at`,
		accId,
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
		return nil, nil, fmt.Errorf("could not delete account: %w", err)
	}

	return &emp, &acc, nil
}

func (a *AccountService) DeleteCustomerData(c context.Context, tx pgx.Tx, customerId, accId string) (*types.DbCustomer, *types.DbAccount, error) {
	var cust types.DbCustomer
	var acc types.DbAccount

	if err := tx.QueryRow(c, `
		DELETE FROM account.customers
		WHERE id = $1
		RETURNING id, account_id, status, loyalty_points, created_at, updated_at
	`, customerId).Scan(
		&cust.ID,
		&cust.AccountID,
	); err != nil {
		return nil, nil, fmt.Errorf("could not delete customer with id %s: %w", customerId, err)
	}

	if err := tx.QueryRow(c, `
		DELETE FROM account.accounts
		WHERE id = $1
		RETURNING first_name, last_name, email, provider, clerk_id, role, id, created_at, updated_at
	`, accId).Scan(
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
		return nil, nil, fmt.Errorf("could not delete account with id %s: %w", accId, err)
	}

	return &cust, &acc, nil
}
