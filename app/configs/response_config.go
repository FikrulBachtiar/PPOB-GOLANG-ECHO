package configs

import (
	"github.com/labstack/echo/v4"
)

type Response struct {
	Status  int         `json:"status"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Error   string      `json:"error"`
}

type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Error   string      `json:"error,omitempty"`
}

func (res *Response) ResponseMiddleware(next echo.Context) error {
	var data interface{}
	var responses APIResponse
	
	if res.Error == "" {
		if res.Data == nil {
			data = make(map[string]interface{})
		} else {
			data = res.Data;
		}

		responses = APIResponse{
			Code: res.Code,
			Message: res.Message,
			Data: data,
		}
	} else {
		if res.Data == nil {
			data = make(map[string]interface{})
		} else {
			data = res.Data;
		}

		responses = APIResponse{
			Code: res.Code,
			Message: "Unexpected error occurred on the server.",
			Data: data,
			Error: res.Error,
		}
	}

	return next.JSON(res.Status, responses);
}