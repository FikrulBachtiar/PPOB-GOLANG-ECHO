package domain

type RequestOtpPayload struct {
	Msisdn string `json:"msisdn" validate:"required"`
	Type   string `json:"type" validate:"required"`
}

type VerificationOtpPayload struct {
	Msisdn   string `json:"msisdn" validate:"required"`
	Type     string `json:"type" validate:"required"`
	OtpToken string `json:"otpToken" validate:"required"`
}

type RequestOtpResponse struct {
	Otp string `json:"otp"`
}

type ExpiredDurationOTP struct {
	Duration     int    `json:"duration"`
	DurationType string `json:"duration_type"`
}

type GetUserAttempt struct {
	IdUser              int     `json:"id_user"`
	AttemptRequestOtp   int     `json:"attempt_request_otp"`
	GroupId             *string `json:"group_id"`
	AttemptVerification *int    `json:"attempt"`
	Otp                 *string `json:"otp"`
}

type GetUserOTP struct {
	IdUser  int    `json:"id_user"`
	Attempt int    `json:"attempt"`
	Otp     string `json:"otp"`
}

type InsertRequestOTP struct {
	IdUser    int    `json:"id_user"`
	Type      string `json:"type"`
	Otp       string `json:"otp"`
	CratedOn  string `json:"created_on"`
	ExpiredOn string `json:"expired_on"`
	GroupId   string `json:"group_id"`
}