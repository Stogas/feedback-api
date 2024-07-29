package feedbacktypes

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Issue struct {
	gorm.Model
	Name string
}

type Report struct {
	gorm.Model
	UUID      uuid.UUID `json:"uuid" binding:"required" gorm:"uniqueIndex"`
	Satisfied *bool     `json:"satisfied" binding:"required"`
	Comment   string    `json:"comment" binding:"max=1000"`
	IssueID   *int      `json:"issue_id"`
	Issue     *Issue    `json:"-"`
}
