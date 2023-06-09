// Package main the main app package
package main

import (
	"database/sql"
	"fmt"

	"github.com/Mekawy5/money-transfer/accounts"
)

// init function runs first in the package, setup app dependencies.
func init() {
	// setup database connection
	{
		conn, err := sql.Open("sqlite3", ":memory:")
		if err != nil {
			panic(err)
		}

		ctx.DBConn = conn
	}

	fmt.Println("Connection established.")

	// load accounts data
	{
		accounts.LoadAccounts(ctx)
	}

	fmt.Println("Accounts Loaded, Ready to make transfers.")
}
