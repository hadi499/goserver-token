package models

import "time"

type Product struct {
	Id        int64  `gorm:"primaryKey" json:"id"`
	Name      string `gorm:"type:varchar(255)" json:"name"`
	Price     int64  `gorm:"type:int" json:"price"`
	Image     string `gorm:"type:varchar(255)" json:"image"`
	UserId    int64  `gorm:"index" json:"user_id"`
	User      User   `gorm:"foreignKey:UserId"` // Menyatakan relasi dengan model User
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ProductResponse struct {
	ID        int64       `json:"id"`
	Name      string      `json:"name"`
	Price     int64       `json:"price"`
	Image     string      `json:"image"`
	User      UserMinimal `json:"user"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserMinimal struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}
