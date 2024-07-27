package feedbacktypes

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Satisfaction struct {
	gorm.Model
	UUID      uuid.UUID `json:"uuid_v4" binding:"required" gorm:"uniqueIndex"`
	Satisfied *bool     `json:"satisfied" binding:"required"`
}
