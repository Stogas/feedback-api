package dto

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type ReportRequest struct {
	UUID      uuid.UUID       `json:"uuid" binding:"required"`
	Satisfied *bool           `json:"satisfied" binding:"required"`
	Comment   string          `json:"comment" binding:"max=1000"`
	IssueID   *int            `json:"issue_id"`
	Metadata  *datatypes.JSON `json:"metadata" binding:"max=2048"`
}
