package helper

import (
	"net/http"
	"reflect"

	"github.com/labstack/echo"
)

// HTTPResponse abstract interface
type HTTPResponse interface {
	JSON(c echo.Context) error
}

type (
	httpResponse struct {
		Success bool        `json:"success"`
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Meta    interface{} `json:"meta,omitempty"`
		Data    interface{} `json:"data,omitempty"`
		Errors  interface{} `json:"errors,omitempty"`
	}

	Meta struct {
		Page         int64 `json:"page"`
		Limit        int64 `json:"limit"`
		TotalRecords int64 `json:"totalRecords"`
		TotalPages   int   `json:"totalPages"`
	}
)

// NewJSONResponse for create common response, data must in first params and meta in second params
func NewJSONResponse(code int, message string, params ...interface{}) HTTPResponse {
	commonResponse := new(httpResponse)

	for _, param := range params {
		refValue := reflect.ValueOf(param)
		if refValue.Kind() == reflect.Ptr {
			refValue = refValue.Elem()
		}
		param = refValue.Interface()

		switch param.(type) {
		case Meta:
			commonResponse.Meta = param
		default:
			commonResponse.Data = param
		}
	}

	commonResponse.Success = code < http.StatusBadRequest
	commonResponse.Code = code
	commonResponse.Message = message
	return commonResponse
}

// JSON for set http JSON response (Content-Type: application/json)
func (resp *httpResponse) JSON(c echo.Context) error {
	if resp.Data == nil {
		resp.Data = struct{}{}
	}
	return c.JSON(resp.Code, resp)
}
