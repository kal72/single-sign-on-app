package entities

import "time"

type BaseEntity struct {
	//CreatedAt *time.Time `json:"created_at,omitempty" sql:"DEFAULT:current_timestamp"`
	//UpdatedAt *time.Time `json:"updated_at,omitempty" sql:"DEFAULT:current_timestamp"`
	//DeletedAt *time.Time `sql:"index" json:"deleted_at,omitempty"`
	CreatedAt *time.Time `json:"-" sql:"DEFAULT:current_timestamp"`
	UpdatedAt *time.Time `json:"-" sql:"DEFAULT:current_timestamp"`
	DeletedAt *time.Time `sql:"index" json:"-"`

	CreatedTime *string `gorm:"-" json:"created_time,omitempty"`
	UpdatedTime *string `gorm:"-" json:"updated_time,omitempty"`
	DeletedTime *string `gorm:"-" json:"deleted_time,omitempty"`
}
