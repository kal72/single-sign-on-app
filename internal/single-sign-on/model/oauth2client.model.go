package model

import "encoding/json"

type Oauth2Clients struct {
	ID          string          `json:"client_id"`
	Secret      string          `json:"secret"`
	Domain      string          `json:"domain" binding:"required" form:"domain" binding:"required"`
	Name        string          `gorm:"-" json:"name,omitempty" form:"name" binding:"required"`
	CallbackUrl string          `gorm:"-" json:"callback_url,omitempty" form:"callback_url" binding:"required"`
	Description string          `gorm:"-" json:"description,omitempty" form:"description" binding:"required"`
	Data        json.RawMessage `json:"data,omitempty"`
}
