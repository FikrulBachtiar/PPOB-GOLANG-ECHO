package domain

type UserDetailResponse struct {
	Msisdn        string `json:"msisdn"`
	FullName      string `json:"fullname"`
	Email         string `json:"email"`
	EmailStatus   string `json:"email_status"`
	EmailStatusEn string `json:"email_status_en"`
	KtpStatus     string `json:"ktp_status"`
	KtpStatusEn   string `json:"ktp_status_en"`
	Uuid          string `json:"uuid"`
	Question      bool   `json:"question"`
	UrlProPic     string `json:"url_pro_pic"`
}

type UserBalanceResponse struct {
	Currency  string `json:"currency"`
	Value     string `json:"value"`
	KtpStatus bool   `json:"ktp_status"`
}

type UserPointResponse struct {
	Currency   string `json:"currency"`
	TotalPoint string `json:"total_point"`
	IsTap      bool   `json:"is_tap"`
}

type SpGetBalance struct {
	Currency string `json:"currency"`
	Value    string `json:"value"`
}

type SpGetPoint struct {
	Currency string `json:"currency"`
	Value    string `json:"value"`
	IsTap    bool   `json:"is_tap"`
}

type PointTypeResponse struct {
	IdPointType     int    `json:"id_point_type"`
	PointTypeName   string `json:"point_type_name"`
	PointTypeNameEn string `json:"point_type_name_en"`
}

type ListPointResponse struct {
	ProductCategory string `json:"product_category"`
	Direction       string `json:"direction"`
	Currency        string `json:"currency"`
	PointTotal      string `json:"point_total"`
	TransactionDate string `json:"transaction_date"`
}

type PaginationDefault struct {
	TotalPages  int         `json:"total_pages"`
	CurrentPage int         `json:"current_page"`
	TotalData   int         `json:"total_data"`
	ListData    interface{} `json:"list_data"`
}

type ListPointRequest struct {
	IdUser      int `json:"id_user"`
	Page        int `json:"page"`
	Limit       int `json:"limit"`
	IdPointType int `json:"id_point_type"`
}