package models

import "time"

type Asset struct {
	Name      string    `json:"name"`
	Uid       int       `json:"uid"`
	Data      string    `json:"data"`
	CreatedAt time.Time `json:"created_at"`
	Deleted   bool      `json:"deleted"`
}

func (a Asset) TableName() string {
	return "assets"
}
