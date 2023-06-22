package middleware

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"ppob/app/configs"
	"strings"

	"github.com/golang-jwt/jwt"
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

func (hv *HeaderValidator) JWTValidation(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		
		deviceID := c.Request().Header.Get("device_id")
		notificationID := c.Request().Header.Get("notification_id")
		auth := c.Request().Header.Get("Authorization");
		isBearerAuth := strings.HasPrefix(auth, "Bearer ");

		if !isBearerAuth {
			response := &configs.Response{
				Status: http.StatusUnauthorized,
				Code: 5694,
				DB: hv.DB,
				Type: 98,
			}
			return response.ResponseMiddleware(c);
		}

		tokens := strings.TrimPrefix(auth, "Bearer ");

		var UserStatusCode string
		sqlQuery := fmt.Sprintf(`
		SELECT
			tu.user_status_code 
		FROM
			security.t_session_app tsa
		JOIN
			user_management.t_users tu 
		ON
			tu.id_user = tsa.id_user
		WHERE
			token = '%s'
			AND status = 1
			AND device_id = '%s'
			AND notification_id = '%s'
		`, tokens, deviceID, notificationID);
		
		if err := hv.DB.QueryRow(sqlQuery).Scan(&UserStatusCode); err != nil {
			if err == sql.ErrNoRows {
				response := &configs.Response{
					Status: http.StatusUnauthorized,
					Code: 5694,
					DB: hv.DB,
					Type: 98,
				}
				return response.ResponseMiddleware(c);
			} else {
				response := &configs.Response{
					Status: http.StatusInternalServerError,
					Code: 2037,
					DB: hv.DB,
					Error: err.Error(),
					Type: 98,
				}
				return response.ResponseMiddleware(c);
			}
		}

		if UserStatusCode == os.Getenv("STATUS_ACCOUNT_BLOCKED") {
			response := &configs.Response{
				Status: http.StatusForbidden,
				Code: 6492,
				DB: hv.DB,
				Type: 98,
			}
			return response.ResponseMiddleware(c);
		}

		token, err := jwt.Parse(tokens, func (token *jwt.Token)  ( interface{}, error ) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"]);
			}

			return []byte(os.Getenv("JWT_KEY")), nil;
		});
		
		if err != nil {
			if ve, ok := err.(*jwt.ValidationError); ok {
				if ve.Errors&jwt.ValidationErrorMalformed != 0 {
					response := &configs.Response{
						Status: http.StatusBadRequest,
						Code: 4535,
						DB: hv.DB,
						Type: 98,
					}
					return response.ResponseMiddleware(c);
				} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
					response := &configs.Response{
						Status: http.StatusForbidden,
						Code: 5875,
						DB: hv.DB,
						Type: 98,
					}
					return response.ResponseMiddleware(c);
				} else {
					response := &configs.Response{
						Status: http.StatusUnauthorized,
						Code: 5694,
						DB: hv.DB,
						Type: 98,
					}
					return response.ResponseMiddleware(c);
				}
			}
		}

		claims, ok := token.Claims.(jwt.MapClaims);

		if !ok {
			response := &configs.Response{
				Status: http.StatusInternalServerError,
				Code: 3288,
				DB: hv.DB,
				Type: 98,
			}
			return response.ResponseMiddleware(c);
		}

		msisdn := claims["Msisdn"].(string);
		IdUser := int(claims["IdUser"].(float64));

		c.Set("msisdn", msisdn);
		c.Set("id_user", IdUser);

		return next(c);
	}
}

func HeaderContentType(next echo.HandlerFunc) echo.HandlerFunc {
	return func (c echo.Context) error {
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON);
		return next(c);
	}
}