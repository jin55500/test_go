package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username   string  `gorm:"column:username;type:varchar(255);not null"`
	Password   string  `gorm:"column:password;type:varchar(255);not null"`
	Name       string  `gorm:"column:name;type:varchar(255);not null"`
	Surname    string  `gorm:"column:surname;type:varchar(255);not null"`
	BankNumber string  `gorm:"column:bank_number;type:varchar(10);not null"`
	Credit     float64 `gorm:"column:credit;type:decimal(12,3);default:1000"`
}

type UserCreateForm struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	Name       string `json:"name" binding:"required"`
	Surname    string `json:"surname" binding:"required"`
	BankNumber string `json:"bankNumber" binding:"required,len=10,number"`
}

type UserUpdateForm struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" `
	Name       string `json:"name" binding:"required"`
	Surname    string `json:"surname" binding:"required"`
	BankNumber string `json:"bankNumber" binding:"required,len=10,number"`
}

// UserSwagger is used for Swagger documentation purposes
type UserSwagger struct {
	ID         uint    `json:"id"`
	Username   string  `json:"Username"`
	Name       string  `json:"Name"`
	Surname    string  `json:"Surname"`
	BankNumber string  `json:"BankNumber"`
	Credit     float64 `json:"Credit"`
}
