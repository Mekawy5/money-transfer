// Package main the main app package
package main

import (
	"github.com/go-playground/validator"
)

// CustomValidator for echo
type CustomValidator struct {
	validator *validator.Validate
}

// Validate function
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}
