package accounts

import (
	"errors"

	"github.com/Mekawy5/money-transfer/internals/appctx"
)

func GetAccounts(ctx appctx.Context) ([]Account, error) {
	var accounts []Account

	stmt, err := ctx.DBConn.Prepare("SELECT id, name, balance FROM accounts")
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var account Account
		err := rows.Scan(&account.ID, &account.Name, &account.Balance)
		if err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}

func Transfer(r TransferRequest, ctx appctx.Context) (Account, Account, error) {
	// begin transaction
	tx, err := ctx.DBConn.Begin()
	if err != nil {
		return Account{}, Account{}, err
	}

	var from, to Account

	err = tx.QueryRow("SELECT id, name, balance FROM accounts WHERE id = ?", r.SenderID).Scan(&from.ID, &from.Name, &from.Balance)
	if err != nil {
		tx.Rollback()
		return Account{}, Account{}, errors.New("sender doesn't exist")
	}

	// validate the balance more than deducted amount.
	if from.Balance-r.Amount < 0 {
		tx.Rollback()
		return Account{}, Account{}, errors.New("no enough balance")
	}

	err = tx.QueryRow("SELECT id, name, balance FROM accounts WHERE id = ?", r.ReceiverID).Scan(&to.ID, &to.Name, &to.Balance)
	if err != nil {
		tx.Rollback()
		return Account{}, Account{}, errors.New("receiver doesn't exist")
	}

	// update balance in the structs
	from.Balance -= r.Amount
	to.Balance += r.Amount

	_, err = tx.Exec("UPDATE accounts SET balance = ? WHERE id = ?", from.Balance, from.ID)
	if err != nil {
		tx.Rollback()
		return Account{}, Account{}, errors.New("balance transfer failed")
	}
	_, err = tx.Exec("UPDATE accounts SET balance = ? WHERE id = ?", to.Balance, to.ID)
	if err != nil {
		tx.Rollback()
		return Account{}, Account{}, errors.New("balance transfer failed")
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return Account{}, Account{}, errors.New("balance transfer failed")
	}

	return from, to, nil
}
