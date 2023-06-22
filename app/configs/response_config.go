package configs

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/labstack/echo/v4"
)

type Response struct {
	Status  int
	Code    int
	Data    interface{}
	Error   string
	DB 		*sql.DB
	Type    int
}

type APIResponse struct {
	Code    int         `json:"code"`
	Message Messages  	`json:"message"`
	Data    interface{} `json:"data"`
	Error   string      `json:"error,omitempty"`
}

type Messages struct {
	MessageIna string `json:"message_ina"`
	MessageEn string `json:"message_en"`
}

func (res *Response) ResponseMiddleware(next echo.Context) error {
	var data interface{}
	var responses APIResponse
	var message, message_en string
	
	if res.Error == "" {
		sqlQuery := fmt.Sprintf("SELECT message, message_en FROM master.t_message_response WHERE status = 1 AND code = %d AND endpoint_type = %d", res.Code, res.Type);
		err := res.DB.QueryRow(sqlQuery).Scan(&message, &message_en);
		if err != nil {
			if err == sql.ErrNoRows {
				res.Error = errors.New("Haven't set [message] in the database yet").Error();
			} else {
				res.Error = err.Error();
			}
			responses = res.errorInternal();
		} else {
			if res.Data == nil {
				data = make(map[string]interface{})
			} else {
				data = res.Data;
			}
	
			responses = APIResponse{
				Code: res.Code,
				Message: Messages{
					MessageIna: message,
					MessageEn: message_en,
				},
				Data: data,
			}
		}
	} else {
		responses = res.errorInternal();
	}

	return next.JSON(res.Status, responses);
}

func (res *Response) errorInternal() APIResponse {
	data := make(map[string]interface{})
	responses := APIResponse{
		Code: res.Code,
		Message: Messages{
			MessageIna: "Terjadi kesalahan tak terduga di server",
			MessageEn: "Unexpected error occurred on the server",
		},
		Data: data,
		Error: res.Error,
	}

	return responses;
}