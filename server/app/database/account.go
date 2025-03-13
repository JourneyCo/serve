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
		&account.CreatedAt, &account.UpdatedAt,
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
