package models

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
	UUID      uuid.UUID `binding:"required" gorm:"uniqueIndex"`
	Satisfied *bool     `binding:"required"`
	Comment   string    `binding:"max=1000"`
	IssueID   *int
	Issue     *Issue
}
