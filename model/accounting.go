package model

import "gorm.io/gorm"

type Accounting struct {
	gorm.Model
	UserTransfer uint    `gorm:"column:user_transfer;type:varchar(255);not null"`
	Credit       float64 `gorm:"column:credit;type:decimal(12,3);not null"`
	UserReceiver uint    `gorm:"column:user_receiver;type:varchar(255);not null"`
	BankNumber   string  `gorm:"column:bank_number;type:varchar(10);not null"`
}

type TransferCreateForm struct {
	BankNumber string  `json:"BankNumber" binding:"required"`
	Credit     float64 `json:"credit" binding:"required"`
}

type AccountingResult struct {
	ID           uint    `json:"id" `
	UserTransfer uint    `json:"UserTransfer"`
	Credit       float64 `json:"Credit" `
	UserReceiver uint    `json:"UserReceiver"`
	BankNumber   string  `json:"BankNumber"`
	Type         string  `json:"Type" `
}
