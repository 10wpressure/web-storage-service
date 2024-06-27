package dto

type GetAssetByName struct {
	Name   string `json:"name"`
	UserID int    `json:"user_id"`
}
