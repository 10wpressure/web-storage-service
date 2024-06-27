package dto

type ListAssets struct {
	UserID int `json:"user_id"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}
