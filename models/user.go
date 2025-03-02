package models

import "time"

type User struct {
	Id        int64  `gorm:"primaryKey" json:"id"`
	Username  string `gorm:"type:varchar(100)" json:"username"`
	Email     string `gorm:"type:varchar(100)" json:"email" validate:"required,email"`
	Password  string `gorm:"type:varchar(255)" json:"password" validate:"required,min=6"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
