package models

type User struct {
	ID     uint   `json:"id" gorm:"primaryKey;autoIncrement:true"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	AuthID uint   `json:"auth_id"`
}
