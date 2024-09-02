package dto

type TelephoneRequest struct {
	PhoneNumber string `json:"phoneNumber" validate:"required,vphone"`
}

type CreateUserRequest struct {
	PhoneNumber  string `json:"phoneNumber" validate:"required,vphone"`
	SMSCode      string `json:"code" validate:"required,min=4,max=6"`
	Username     string `json:"username"`
	PasswordHash string `json:"passwordHash"`
	Avatar       string `json:"avatar"`
}

type LoginRequest struct {
	PhoneNumber  string `json:"phoneNumber" validate:"required,vphone"`
	LoginType    string `json:"loginType" validate:"oneof=sms psw"`
	SMSCode      string `json:"code" validate:"required_if=LoginType sms"`
	PasswordHash string `json:"passwordHash" validate:"required_if=LoginType psw"`

	UserAgent string
	ClientIP  string
}

type UpdateUserRequest struct {
	ID       int    `json:"id" validate:"required,min=1"`
	Username string `json:"username" validate:"required"`
	Avatar   string `json:"avatar" validate:"required"`
}
