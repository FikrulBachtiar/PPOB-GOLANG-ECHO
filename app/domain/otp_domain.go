package domain

type RequestOtpPayload struct {
	Msisdn string `json:"msisdn" validate:"required"`
	Type   string `json:"type" validate:"required"`
}

type RequestOtpHeader struct {
	DeviceID string `json:"device_id"`
}

type VerificationOtpPayload struct {
	Msisdn   string `json:"msisdn" validate:"required"`
	Type     string `json:"type" validate:"required"`
	OtpToken string `json:"otpToken" validate:"required"`
}

type RequestOtpResponse struct {
	Otp        string `json:"otp"`
	Attempt    int    `json:"attempt"`
	MaxAttempt int    `json:"max_attempt"`
}

type ExpiredDurationOTP struct {
	Duration     int    `json:"duration"`
	DurationType string `json:"duration_type"`
}

type ParamValue struct {
	Value   *string `json:"value"`
	Measure *string `json:"measure"`
}

type GetUserAttempt struct {
	IdUser              int     `json:"id_user"`
	UserStatusCode      string  `json:"user_status_code"`
	AttemptRequestOtp   int     `json:"attempt_request_otp"`
	GroupId             *string `json:"group_id"`
	AttemptVerification *int    `json:"attempt"`
	Otp                 *string `json:"otp"`
	ExpiredOn           *string `json:"expired_on"`
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
	DeviceID  string `json:"device_id"`
}