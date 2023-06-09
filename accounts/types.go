package accounts

import (
	"encoding/json"
	"strconv"
)

type TransferRequest struct {
	SenderID   string  `form:"sender_id" validate:"required,uuid"`
	ReceiverID string  `form:"receiver_id" validate:"required,uuid"`
	Amount     float64 `form:"amount" validate:"required,numeric,min=0"`
}

type GetAccountRequest struct {
	ID string `query:"id" form:"id" path:"id" validate:"required,uuid"`
}

// Account struct holds account data
type Account struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Balance float64 `json:"balance"`
}

func (a *Account) UnmarshalJSON(data []byte) error {
	type Alias Account
	aux := &struct {
		*Alias
		Balance string `json:"balance"`
	}{
		Alias: (*Alias)(a),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	balance, err := strconv.ParseFloat(aux.Balance, 64)
	if err != nil {
		return err
	}
	a.Balance = balance
	return nil
}
