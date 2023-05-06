package models

import "gorm.io/gorm"

type Vault struct {
	gorm.Model
	Name   string
	UserID uint
}

type VaultDetail struct {
	ID     uint
	Name   string
	UserID uint
}
