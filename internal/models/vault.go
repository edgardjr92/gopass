package models

import "gorm.io/gorm"

type Vault struct {
	gorm.Model
	Name   string
	UserID uint
}
