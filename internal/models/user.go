package models

import "github.com/fzzp/hotel-booking-api/internal/dto"

type User struct {
	ID           uint   `db:"id" json:"id"`
	PhoneNumber  string `db:"phone_number" json:"phoneNumber"`
	PasswordHash string `db:"password_hash" json:"-"`
	Username     string `db:"username" json:"username"`
	Avatar       string `db:"avatar" json:"avatar"`
	Role         int8   `db:"role" json:"role"`
	CreatedAt    string `db:"created_at" json:"createdAt"`
	UpdatedAt    string `db:"updated_at" json:"updatedAt"`
	IsDeleted    int8   `db:"is_deleted" json:"-"`
}

func (u *User) ToDto(aToken, rToken string) dto.LoginResponse {
	return dto.LoginResponse{
		ID:           u.ID,
		PhoneNumber:  u.PhoneNumber,
		Username:     u.Username,
		Avatar:       u.Avatar,
		AccessToken:  aToken,
		RefreshToken: rToken,
	}
}
