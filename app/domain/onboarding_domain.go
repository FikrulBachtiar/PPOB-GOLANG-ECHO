package domain

import "github.com/golang-jwt/jwt"

type CheckPayload struct {
	Msisdn string `json:"msisdn" validate:"required"`
}

type LoginPayload struct {
	Msisdn         string `json:"msisdn" validate:"required"`
	Pin            string `json:"pin" validate:"required"`
	NotificationID string `json:"notificationID"`
}

type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
	FullName     string `json:"fullname"`
}

type OnboardingHeader struct {
	DeviceID 		string `json:"device_id" validate:"required"`
	OsName         	string `json:"os_name" validate:"required"`
	OsVersion      	string `json:"os_version" validate:"required"`
	DeviceModel    	string `json:"device_model" validate:"required"`
	AppVersion     	string `json:"app_version" validate:"required"`
	Longitude      	string `json:"longitude" validate:"required"`
	Latitude       	string `json:"latitude" validate:"required"`
	NotificationID 	string `json:"notification_id" validate:"required"`
}

type LogoutPayload struct {
	Msisdn string `json:"msisdn" validate:"required"`
}

type CheckResponse struct {
	State       string `json:"state"`
	IsBlocked   bool   `json:"isBlocked"`
	IsQuestion  bool   `json:"isQuestion"`
	EmailStatus string `json:"emailStatus"`
}

type GetUserByNumber struct {
	IdUser         int     `json:"id_user"`
	Fullname       string  `json:"fullname"`
	Uuid		   string `json:"uuid"`
	IsQuestion     bool    `json:"is_question"`
	LoginStatus    *int    `json:"login_status"`
	LastLogin      *string `json:"last_login"`
	UserStatusCode string  `json:"user_status_code"`
	EmailStatus    bool    `json:"email_status"`
	Pin            *string `json:"pin"`
}

type InsertUser struct {
	IdUser int `json:"id_user"`
	IsQuestion      bool `json:"is_question"`
	IsEmailVerified bool `json:"is_email_verified"`
}

type InsertToken struct {
	IdUser 			int		`json:"id_user"`
	Msisdn 			string	`json:"msisdn"`
	Token 			string	`json:"token"`
	RefreshToken 	string	`json:"refresh_token"`
	NotificationID 	string	`json:"notification_id"`
	ExpiredDate 	string	`json:"expired_date"`
	CreatedOn 		string	`json:"created_on"`
	DeviceID 		string	`json:"device_id"`
}

type InsertLogin struct {
	IdUser         int 	  `json:"id_user"`
	OsName         string `json:"os_name"`
	OsVersion      string `json:"os_version"`
	DeviceModel    string `json:"device_model"`
	DeviceId       string `json:"device_id"`
	AppVersion     string `json:"app_version"`
	Longitude      string `json:"longitude"`
	Latitude       string `json:"latitude"`
	NotificationID string `json:"notification_id"`
	Pin            string `json:"pin"`
	LoginOn        string `json:"login_on"`
	StatusCode     int	  `json:"status_code"`
	Fullname       string `json:"fullname"`
	Msisdn         string `json:"msisdn"`
	Uuid           string `json:"uuid"`
}

type ClaimsJWT struct {
	jwt.StandardClaims
	IdUser int
	Msisdn         string
	NotificationID string
	DeviceID       string
}