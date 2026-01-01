package entity

import "time"

type Token struct {
	ID        string    `json:"id"`
	UserID    int       `json:"user_id"`
	TokenType string    `json:"token_type"`
	HashToken string    `json:"hash_token"`
	ExpiredAt time.Time `json:"expired_at"`
	IsRevoked bool      `json:"is_revoked"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
