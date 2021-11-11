package models

import (
	"github.com/google/uuid"

	"github.com/sithumonline/quick-note/transact/base"
)

type Note struct {
	base.Base
	Note     string    `json:"note,omitempty"`
	UserId   uuid.UUID `json:"userId,omitempty"`
	Complete bool      `gorm:"default:false" json:"complete"`
	Public   bool      `gorm:"default:false" json:"public"`
}
