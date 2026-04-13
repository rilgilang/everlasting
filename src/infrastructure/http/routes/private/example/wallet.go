package example

import (
	"time"
)

type (
	WalletAccount struct {
		Type string `json:"type" example:"doctor"`
		ID   string `json:"id" example:"8f364610-6f12-4bed-b7d1-7ea1892803c9"`
	}

	WalletItem struct {
		ID              string        `json:"_id" example:"8f364610-6f12-4bed-b7d1-7ea1892803c7"`
		Account         WalletAccount `json:"account"`
		Name            string        `json:"name" example:"markonah"`
		Balance         float64       `json:"balance" example:"240000"`
		Status          string        `json:"status" example:"active"`
		LastTransaction time.Time     `json:"last_transaction"`
		CreatedAt       time.Time     `json:"created_at"`
		UpdatedAt       time.Time     `json:"updated_at"`
	}

	WalletResponse struct {
		Data WalletItem `json:"data"`
		Meta struct {
			Code    string `json:"code" example:"ok"`
			Message string `json:"message" example:"Ok"`
		} `json:"meta"`
	}

	WalletsResponse struct {
		Data []WalletItem
		Meta struct {
			Code    string `json:"code" example:"ok"`
			Message string `json:"message" example:"Ok"`
			Detail  struct {
				Pagination struct {
					CurrentPage uint `json:"current_page" example:"1"`
					MaxPage     uint `json:"max_page" example:"2"`
					TotalData   uint `json:"total_data" example:"10"`
				} `json:"pagination"`
			} `json:"detail"`
		} `json:"meta"`
	}
)
