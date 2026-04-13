package example

type (
	SuccessfulResetPassword struct {
		Meta struct {
			Code    string `json:"code" example:"ok"`
			Message string `json:"message" example:"Ok"`
		} `json:"meta"`
	}

	SuccessfulRedeemToken struct {
		Data string `json:"data" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDQzNjI3MTEsImlkIjoiN2IyYWE2MTEtMmViMS00NDI0LWJjNWEtYjg2N2I3MTBmYWQ3Iiwic3ViamVjdCI6InJlc2V0X3Bhc3N3b3JkIn0.G_Uofrg80H4WwoRZzjUH6-UbT_4otd8VwMqNOZOyB5w"`
		Meta struct {
			Code    string `json:"code" example:"ok"`
			Message string `json:"message" example:"Ok"`
		} `json:"meta"`
	}
)
