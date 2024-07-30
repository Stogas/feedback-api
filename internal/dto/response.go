package dto

import (
	"github.com/Stogas/feedback-api/internal/models"
)

type ReportResponse struct {
	ReportRequest
}

func MapReportToReportResponse(report models.Report) ReportResponse {
	return ReportResponse{
		ReportRequest: ReportRequest{
			UUID:      report.UUID,
			Satisfied: report.Satisfied,
			Comment:   report.Comment,
			IssueID:   report.IssueID,
			Metadata:  report.Metadata,
		},
	}
}

type IssueResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func MapIssuesToIssueResponses(issues []models.Issue) []IssueResponse {
	response := make([]IssueResponse, len(issues))
	for i, issue := range issues {
		response[i] = IssueResponse{
			ID:   issue.ID,
			Name: issue.Name,
		}
	}
	return response
}
