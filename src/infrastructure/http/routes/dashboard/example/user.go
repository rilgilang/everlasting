package example

import (
	"time"
)

type (
	UserItem struct {
		ID              string    `json:"_id" example:"8f364610-6f12-4bed-b7d1-7ea1892803c7"`
		Name            string    `json:"name" example:"markonah"`
		Balance         float64   `json:"balance" example:"240000"`
		Status          string    `json:"status" example:"active"`
		LastTransaction time.Time `json:"last_transaction"`
		CreatedAt       time.Time `json:"created_at"`
		UpdatedAt       time.Time `json:"updated_at"`
	}

	UserResponse struct {
		Data UserItem `json:"data"`
		Meta struct {
			Code    string `json:"code" example:"ok"`
			Message string `json:"message" example:"Ok"`
		} `json:"meta"`
	}

	UsersResponse struct {
		Data []UserItem
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
