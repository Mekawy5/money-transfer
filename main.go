// Package main the main app package
package main

import (
	"net/http"

	"github.com/Mekawy5/money-transfer/accounts"
	"github.com/Mekawy5/money-transfer/internals/appctx"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
)

var (
	ctx appctx.Context
)

// main function entrypoint of the app
func main() {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	e.HideBanner = true
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	e.GET("/accounts", func(c echo.Context) error {
		accounts, err := accounts.GetAccounts(c.Request(), ctx)

		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": true,
			"result":  accounts,
		})
	})

	e.POST("/accounts/transfer", func(c echo.Context) error {
		var r accounts.TransferRequest

		if err := c.Bind(&r); err != nil {
			return c.JSON(400, map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			})
		}

		if err := c.Validate(r); err != nil {
			return c.JSON(422, map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			})
		}

		if r.ReceiverID == r.SenderID {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"error":   "cannot transfer money for the same account",
			})
		}

		from, to, err := accounts.Transfer(r, ctx)

		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": true,
			"result": map[string]interface{}{
				"from": from,
				"to":   to,
			},
		})
	})

	e.Logger.Fatal(e.Start(":8080"))
}
