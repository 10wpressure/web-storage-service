package models

type Session struct {
	ID        string `json:"id"`
	UID       int    `json:"uid"`
	CreatedAt string `json:"created_at"`
	ExpiresAt string `json:"expires_at"`
	IPAddress string `json:"ip_address"`
	Active    bool   `json:"active"`
}

func (s Session) TableName() string {
	return "sessions"
}
