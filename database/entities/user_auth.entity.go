package entities

type UserAuth struct {
	UserId   uint   `gorm:"column:id" json:"-"`
	RoleId   uint   `gorm:"-" json:"role_id"`
	Name     string `gorm:"-" json:"name"`
	Fullname string `gorm:"-" json:"fullname,omitempty"`
	Email    string `gorm:"-" json:"email"`
	Username string `gorm:"-" json:"username,omitempty"`
	Password string `gorm:"-" json:"-"`
	Address  string `gorm:"-" json:"address"`
	Phone    string `gorm:"-" json:"phone"`
	Token    string `gorm:"-" json:"token,omitempty"`
}
