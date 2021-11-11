package base

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Base struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" example:"2e36a7c2-46dc-4107-8499-49ccc85adb25"`
	CreatedAt time.Time      `example:"2021-05-13T05:19:18.789332489Z"`
	UpdatedAt time.Time      `example:"2021-05-13T05:19:18.789332489Z"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
