package models

type User struct {
	ID       uint   `gorm:"primary key"`
	Name     string `gorm:"type:varchar(255);not null"`
	Surname  string `gorm:"type:varchar(255);not null"`
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
}

type Item struct {
	ID          uint    `gorm:"primary key"`
	Name        string  `gorm:"type:varchar(255);not null"`
	Description string  `gorm:"type:text"`
	Amount      uint    `gorm:"not null"`
	Price       float64 `gorm:"not null"`
	Image       string  `gorm:"type:text"`
	Rating      float64
}
