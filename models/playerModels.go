package models

type RegisterUser struct {
	FullName     string `json:"full_name"`
	MobileNumber string `json:"mobile_number"`
	Password     string `json:"password"`
}
type LoginUser struct {
	MobileNumber string `json:"mobile_number"`
	Password     string `json:"password"`
}

type UserInfo struct {
	FullName     string `db:"full_name"`
	MobileNumber string `db:"mobile_number"`
	Password     string `db:"password_hash"`
	UserID       string `db:"user_id"`
}

type ForgotPassword struct {
	MobileNumber string `json:"mobile_number"`
	Password     string `json:"password"`
	Otp          string `json:"otp"`
}
