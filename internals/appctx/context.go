// Package appctx for app level variables
package appctx

import "database/sql"

// Context serves as general input
type Context struct {
	DBConn *sql.DB
}
