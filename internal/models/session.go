package models

import "time"

type Session struct {
	ID           uint      `db:"id" json:"id"`
	UserID       uint      `db:"user_id" json:"userId"`
	TokenID      string    `db:"token_id" json:"tokenId"`
	RefreshToken string    `db:"refresh_token" json:"refreshToken"`
	ClientIP     string    `db:"client_ip" json:"clientIp"`
	UserAgent    string    `db:"user_agent" json:"userAgent"`
	ExpiresAt    time.Time `db:"expires_at" json:"expiresAt"`
	CreatedAt    time.Time `db:"created_at" json:"createdAt"`
}
