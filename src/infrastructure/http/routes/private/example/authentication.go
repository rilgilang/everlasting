package example

type (
	SuccessfulLoginResponse struct {
		Data struct {
			AccessToken  string `json:"access_token" example:"your_access_token"`
			RefreshToken string `json:"refresh_token" example:"your_refresh_token"`
		} `json:"data"`
		Meta struct {
			Code    string `json:"code" example:"ok"`
			Message string `json:"message" example:"Ok"`
		} `json:"meta"`
	}

	InvalidCredentialResponse struct {
		Meta struct {
			Code    string `json:"code" example:"bad_request"`
			Message string `json:"message" example:"invalid credential"`
		} `json:"meta"`
	}

	SuccessfulLogoutResponse struct {
		Meta struct {
			Code    string `json:"code" example:"ok"`
			Message string `json:"message" example:"Logout successful"`
		} `json:"meta"`
	}
)
