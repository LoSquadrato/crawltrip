package database

import "time"

// TODO: use another database for user and token management, SQLite?

type Token struct {
	Token     string    `json:"token"`
	UserID    string    `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

func CreateToken(userID string, expiresAt time.Time) (*Token, error) {
	return nil, nil
}

func GetToken(token string) (*Token, error) {
	return nil, nil
}

func RevokeToken(token string) error {
	return nil
}

func UpdateToken(token *Token) error {
	return nil
}
