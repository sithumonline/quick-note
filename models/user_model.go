package models

import "github.com/sithumonline/quick-note/transact/base"

type User struct {
	base.Base
	Password     string `gorm:"size:100" json:"password"`
	Email        string `gorm:"type:varchar(100);unique" json:"email"`
	Name         string `json:"name"`
	Verification bool
}
