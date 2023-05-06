package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name  string
	Email string
	Psw   string
}

type UserDetail struct {
	Name  string
	Email string
}
