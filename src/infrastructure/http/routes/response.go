package routes

import (
	"net/http"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type (
	Meta struct {
		Status  int         `json:"-"`
		Code    string      `json:"code"`
		Message string      `json:"message"`
		Detail  interface{} `json:"detail,omitempty"`
	}

	ResponseBody struct {
		Data interface{} `json:"data,omitempty"`
		Meta Meta        `json:"meta"`
	}
)

func NewResponse(data interface{}, message, code string, status int, detail interface{}) *ResponseBody {
	meta := Meta{
		Status:  status,
		Code:    code,
		Message: message,
	}
	if detail != nil {
		meta.Detail = detail
	}
	return &ResponseBody{
		Data: data,
		Meta: meta,
	}
}

func NewErrorResponseFromValidator(errors validator.ValidationErrors) *ResponseBody {
	var detail = map[string]interface{}{}
	for _, i := range errors {
		detail[i.Field()] = i.Value()
	}

	return &ResponseBody{
		Meta: Meta{
			Status:  http.StatusBadRequest,
			Code:    "invalid_input",
			Message: "please check your input",
			Detail:  detail,
		},
	}
}

func (r *ResponseBody) GetStatus() int {
	if r.Meta.Status != 0 {
		return r.Meta.Status
	}

	return http.StatusInternalServerError
}

func JsonResponse(c echo.Context, data interface{}, message, code string, status int, detail interface{}) (err error) {
	if code == "" {
		code = "ok"
	}

	if message == "" {
		code = "Ok"
	}

	return c.JSON(status, NewResponse(data, message, code, status, detail))
}
