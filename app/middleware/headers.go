package middleware

import (
	"database/sql"
	"net/http"
	"ppob/app/configs"

	"github.com/labstack/echo/v4"
)

type HeaderData struct {
	DeviceID 		string `json:"device_id" validate:"required"`
	OsName         	string `json:"os_name" validate:"required"`
	OsVersion      	string `json:"os_version" validate:"required"`
	DeviceModel    	string `json:"device_model" validate:"required"`
	AppVersion     	string `json:"app_version" validate:"required"`
	Longitude      	string `json:"longitude" validate:"required"`
	Latitude       	string `json:"latitude" validate:"required"`
	NotificationID 	string `json:"notification_id" validate:"required"`
}

type HeaderValidator struct {	
	DB *sql.DB
}

func (hv *HeaderValidator) HeaderValidator(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		
		data := &HeaderData{
			DeviceID: c.Request().Header.Get("device_id"),
			OsName: c.Request().Header.Get("os_name"),
			OsVersion: c.Request().Header.Get("os_version"),
			DeviceModel: c.Request().Header.Get("device_model"),
			AppVersion: c.Request().Header.Get("app_version"),
			Longitude: c.Request().Header.Get("longitude"),
			Latitude: c.Request().Header.Get("latitude"),
			NotificationID: c.Request().Header.Get("notification_id"),
		};

		if err := c.Validate(data); err != nil {
			response := &configs.Response{
				Status: http.StatusBadRequest,
				Code: 6708,
				DB: hv.DB,
				Type: 99,
			}
			return response.ResponseMiddleware(c);
		}

		return next(c);
	}
}

func HeaderContentType(next echo.HandlerFunc) echo.HandlerFunc {
	return func (c echo.Context) error {
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON);
		return next(c);
	}
}