package entities

type User struct {
	Id              uint   `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	RoleId          uint   `json:"role_id"`
	Name            string `gorm:"type:varchar(50)" json:"name" form:"name" binding:"required"`
	Phone           string `gorm:"type:varchar(13)" json:"phone" form:"phone" binding:"required"`
	Address         string `gorm:"type:varchar(255)" json:"address" form:"address" binding:"required"`
	Email           string `gorm:"type:varchar(50)" json:"email" form:"email" binding:"required"`
	Password        string `gorm:"type:varchar(255);not null" json:"-" form:"password" binding:"required"`
	ConfirmPassword string `gorm:"-" json:"confirm_password" form:"confirm_password" binding:"required"`
	Roles           *Role  `json:"roles,omitempty"`
	Token           string `gorm:"-" json:"token,omitempty"`
}

//this method is helper for get table name
func (m User) TableName() string {
	return "users"
}
