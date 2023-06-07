package domain

type CheckPayload struct {
	Msisdn string `json:"msisdn" validate:"required"`
}

type CheckResponse struct {
	State       string `json:"state"`
	IsBlocked   bool   `json:"isBlocked"`
	IsQuestion  bool   `json:"isQuestion"`
	EmailStatus string `json:"emailStatus"`
}

type GetUserByNumber struct {
	Fullname       string  `json:"fullname"`
	IsQuestion     bool    `json:"is_question"`
	LoginStatus    *int    `json:"login_status"`
	LastLogin      *string `json:"last_login"`
	UserStatusCode string  `json:"user_status_code"`
	EmailStatus    bool    `json:"email_status"`
}

type InsertUser struct {
	IsQuestion      bool `json:"is_question"`
	IsEmailVerified bool `json:"is_email_verified"`
}