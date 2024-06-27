package models

type User struct {
	ID           int    `json:"id"`
	Login        string `json:"login"`
	PasswordHash string `json:"password_hash"`
	CreatedAt    string `json:"created_at"`
}

func (u User) TableName() string {
	return "users"
}
