package dto

import "github.com/google/uuid"

type ReportRequest struct {
	UUID      uuid.UUID `json:"uuid" binding:"required"`
	Satisfied *bool     `json:"satisfied" binding:"required"`
	Comment   string    `json:"comment" binding:"max=1000"`
	IssueID   *int      `json:"issue_id"`
}
