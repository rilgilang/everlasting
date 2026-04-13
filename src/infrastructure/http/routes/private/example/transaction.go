package example

import "time"

type (
	TransactionItem struct {
		ID              string    `json:"_id" example:"8f364610-6f12-4bed-b7d1-7ea1892803c7"`
		Name            string    `json:"name" example:"markonah"`
		Balance         float64   `json:"balance" example:"240000"`
		Status          string    `json:"status" example:"active"`
		LastTransaction time.Time `json:"last_transaction"`
		CreatedAt       time.Time `json:"created_at"`
		UpdatedAt       time.Time `json:"updated_at"`
	}

	TransactionResponse struct {
		Data TransactionItem `json:"data"`
		Meta struct {
			Code    string `json:"code" example:"ok"`
			Message string `json:"message" example:"Ok"`
		} `json:"meta"`
	}
)
