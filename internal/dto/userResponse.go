package dto

import "time"

type LoginResponse struct {
	ID           uint   `db:"id" json:"id"`
	PhoneNumber  string `db:"phone_number" json:"phoneNumber"`
	Username     string `db:"username" json:"username"`
	Avatar       string `db:"avatar" json:"avatar"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type RenewTokenResponse struct {
	AccessToken string    `json:"accessToken"`
	ExpiresAt   time.Time `json:"expiredAt"`
}
