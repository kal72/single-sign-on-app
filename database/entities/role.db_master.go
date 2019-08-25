package entities

type Role struct {
	ID   uint   `gorm:"primary_key;AUTO_INCREMENT" json:"role_id"`
	Name string `gorm:"type:varchar(30)" json:"role"`
}

//this method is helper for get table name
func (m Role) TableName() string {
	return "roles"
}
