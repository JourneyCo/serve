package database

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"serve/models"
)

func GetAccount(ctx context.Context, id string) (models.Account, error) {
	account := models.Account{}

	tx, err := DB.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("error beginning tx")
		return account, err
	}

	// Defer a rollback in case anything fails.
	defer tx.Rollback()

	sqlStatement := `
SELECT * FROM accounts WHERE id = $1`
	row := tx.QueryRowContext(ctx, sqlStatement, id)
	if err = row.Scan(
		&account.ID, &account.FirstName, &account.LastName, &account.Email, &account.CellPhone, &account.TextPermission,
		&account.Lead, &account.CreatedAt, &account.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("account does not exist in db %s", id)
		} else {
			log.Printf("error scanning account")
		}
		return account, err
	}

	// Rows.Err will report the last error encountered by Rows.Scan.
	if err = row.Err(); err != nil {
		log.Printf("error row err account")
		return account, err
	}

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		log.Printf("error committing tx for account")
		return account, err
	}

	return account, nil
}

// GetAccounts returns all accounts
func GetAccounts(ctx context.Context) ([]models.Account, error) {
	accounts := []models.Account{}

	tx, err := DB.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("error beginning tx")
		return accounts, err
	}

	// Defer a rollback in case anything fails.
	defer tx.Rollback()

	sqlStatement := `
SELECT * FROM accounts`
	rows, err := tx.QueryContext(ctx, sqlStatement)
	if err != nil {
		log.Printf("error getting accounts: %v", err)
		return accounts, err
	}
	defer rows.Close()

	for rows.Next() {
		var account models.Account
		if err = rows.Scan(
			&account.ID, &account.FirstName, &account.LastName, &account.Email, &account.CellPhone,
			&account.TextPermission, &account.CreatedAt, &account.UpdatedAt,
		); err != nil {
			log.Printf("error scanning while getting accounts")
			return accounts, err
		}
		accounts = append(accounts, account)
	}

	// Rows.Err will report the last error encountered by Rows.Scan.
	if err = rows.Err(); err != nil {
		log.Printf("error row err while getting accounts")
		return accounts, err
	}

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		log.Printf("error committing tx while getting accounts")
		return accounts, err
	}

	return accounts, nil
}

// GetAccountsByProject returns all accounts associated with a project
// TODO: Needs completed
func GetAccountsByProject(ctx context.Context, proj int) ([]models.Account, error) {
	accounts := []models.Account{}

	tx, err := DB.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("error beginning tx")
		return accounts, err
	}

	// Defer a rollback in case anything fails.
	defer tx.Rollback()

	sqlStatement := `SELECT a.id, a.first, a.last, a.email, a.cellphone, a.text_permission, a.created_at, a.updated_at FROM accounts a INNER JOIN registrations r ON a.id = r.account_id INNER JOIN projects p ON r.project_id = p.id WHERE p.id=$1;`
	rows, err := tx.QueryContext(ctx, sqlStatement, proj)
	if err != nil {
		log.Printf("error getting accounts by project: %v", err)
		return accounts, err
	}
	defer rows.Close()

	for rows.Next() {
		var account models.Account
		if err = rows.Scan(
			&account.ID, &account.FirstName, &account.LastName, &account.Email, &account.CellPhone,
			&account.TextPermission, &account.CreatedAt, &account.UpdatedAt,
		); err != nil {
			log.Printf("error scanning while getting accounts")
			return accounts, err
		}
		accounts = append(accounts, account)
	}

	// Rows.Err will report the last error encountered by Rows.Scan.
	if err = rows.Err(); err != nil {
		log.Printf("error row err while getting accounts")
		return accounts, err
	}

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		log.Printf("error committing tx while getting accounts")
		return accounts, err
	}

	return accounts, nil
}

func PostAccount(ctx context.Context, account models.Account) (models.Account, error) {
	tx, err := DB.BeginTx(ctx, nil)
	if err != nil {
		return models.Account{}, err
	}

	// Defer a rollback in case anything fails.
	defer tx.Rollback()

	var id string
	sqlStatement := `
INSERT INTO accounts (id, first, last, email, cellphone, text_permission)
VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	if err = tx.QueryRow(
		sqlStatement, account.ID, account.FirstName, account.LastName, account.Email, account.CellPhone,
		account.TextPermission,
	).
		Scan(&id); err != nil {
		log.Printf("error inserting account: %v", err)
		return account, err
	}

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		return account, err
	}

	log.Printf("inserted account: %v", account.ID)
	return account, nil
}

// PutAccount changes an existing account
func PutAccount(ctx context.Context, account models.Account) (models.Account, error) {
	tx, err := DB.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("error beginning tx while putting account")
		return account, err
	}

	// Defer a rollback in case anything fails.
	defer tx.Rollback()

	// TODO: UpdatedAt not working
	sqlStatement := `
	UPDATE accounts SET (first, last, email, cellphone, text_permission, lead) = ($1, $2, $3, $4, $5, $6) WHERE id = $7`
	_, err = tx.ExecContext(
		ctx, sqlStatement, account.FirstName, account.LastName, account.Email, account.CellPhone,
		account.TextPermission, account.Lead, account.ID,
	)

	if err != nil {
		log.Printf("error executing update while putting account")
		return account, err
	}

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		log.Printf("error committing tx while putting account")
		return account, err
	}

	return account, nil
}
