package models

type Auth struct {
	ID       uint   `json:"id" gorm:"primaryKey;autoIncrement:true"`
	Username string `json:"username"`
	Password string `json:"password"`
}
