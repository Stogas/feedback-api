package feedbacktypes

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Satisfaction struct {
	gorm.Model
	UUID      uuid.UUID `json:"uuid" binding:"required" gorm:"uniqueIndex"`
	Satisfied *bool     `json:"satisfied" binding:"required"`
}
