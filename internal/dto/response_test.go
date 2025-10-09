package dto

import (
	"testing"
	"time"

	"github.com/Stogas/feedback-api/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

const testMetadataJSON = `{"source": "web", "rating": 5}`

func TestMapReportToReportResponse_BasicMapping(t *testing.T) {
	t.Run("maps complete Report to ReportResponse", func(t *testing.T) {
		testUUID := uuid.New()
		satisfied := true
		issueID := 42
		metadata := datatypes.JSON(testMetadataJSON)

		report := models.Report{
			Model: gorm.Model{
				ID:        1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			UUID:      testUUID,
			Satisfied: &satisfied,
			Comment:   "Great service, very satisfied!",
			IssueID:   &issueID,
			Metadata:  &metadata,
		}

		response := MapReportToReportResponse(report)

		// Verify all fields are correctly mapped
		assert.Equal(t, testUUID, response.UUID)
		assert.NotNil(t, response.Satisfied)
		assert.True(t, *response.Satisfied)
		assert.Equal(t, "Great service, very satisfied!", response.Comment)
		assert.NotNil(t, response.IssueID)
		assert.Equal(t, 42, *response.IssueID)
		assert.NotNil(t, response.Metadata)
		assert.Equal(t, metadata, *response.Metadata)
	})

	t.Run("maps Report with minimal fields", func(t *testing.T) {
		testUUID := uuid.New()
		satisfied := false

		report := models.Report{
			Model: gorm.Model{
				ID: 2,
			},
			UUID:      testUUID,
			Satisfied: &satisfied,
		}

		response := MapReportToReportResponse(report)

		assert.Equal(t, testUUID, response.UUID)
		assert.NotNil(t, response.Satisfied)
		assert.False(t, *response.Satisfied)
		assert.Empty(t, response.Comment)
		assert.Nil(t, response.IssueID)
		assert.Nil(t, response.Metadata)
	})
}

func TestMapReportToReportResponse_EdgeCases(t *testing.T) {
	t.Run("maps Report with nil optional fields", func(t *testing.T) {
		testUUID := uuid.New()
		satisfied := true

		report := models.Report{
			UUID:      testUUID,
			Satisfied: &satisfied,
			Comment:   "",
			IssueID:   nil,
			Metadata:  nil,
		}

		response := MapReportToReportResponse(report)

		assert.Equal(t, testUUID, response.UUID)
		assert.NotNil(t, response.Satisfied)
		assert.True(t, *response.Satisfied)
		assert.Empty(t, response.Comment)
		assert.Nil(t, response.IssueID)
		assert.Nil(t, response.Metadata)
	})

	t.Run("handles zero values correctly", func(t *testing.T) {
		testUUID := uuid.New()
		satisfied := false
		issueID := 0

		report := models.Report{
			UUID:      testUUID,
			Satisfied: &satisfied,
			IssueID:   &issueID,
		}

		response := MapReportToReportResponse(report)

		assert.Equal(t, testUUID, response.UUID)
		assert.NotNil(t, response.Satisfied)
		assert.False(t, *response.Satisfied)
		assert.NotNil(t, response.IssueID)
		assert.Equal(t, 0, *response.IssueID)
	})
}

func TestMapIssuesToIssueResponses_BasicCases(t *testing.T) {
	t.Run("maps multiple Issues to IssueResponses", func(t *testing.T) {
		issues := []models.Issue{
			{
				Model: gorm.Model{ID: 1},
				Name:  "Bug Report",
			},
			{
				Model: gorm.Model{ID: 2},
				Name:  "Feature Request",
			},
			{
				Model: gorm.Model{ID: 3},
				Name:  "Improvement Suggestion",
			},
		}

		responses := MapIssuesToIssueResponses(issues)

		assert.Len(t, responses, 3)

		assert.Equal(t, uint(1), responses[0].ID)
		assert.Equal(t, "Bug Report", responses[0].Name)

		assert.Equal(t, uint(2), responses[1].ID)
		assert.Equal(t, "Feature Request", responses[1].Name)

		assert.Equal(t, uint(3), responses[2].ID)
		assert.Equal(t, "Improvement Suggestion", responses[2].Name)
	})

	t.Run("maps single Issue to IssueResponse", func(t *testing.T) {
		issues := []models.Issue{
			{
				Model: gorm.Model{ID: 42},
				Name:  "Critical Bug",
			},
		}

		responses := MapIssuesToIssueResponses(issues)

		assert.Len(t, responses, 1)
		assert.Equal(t, uint(42), responses[0].ID)
		assert.Equal(t, "Critical Bug", responses[0].Name)
	})
}

func TestMapIssuesToIssueResponses_EdgeCases(t *testing.T) {
	t.Run("returns empty slice for empty input", func(t *testing.T) {
		issues := []models.Issue{}

		responses := MapIssuesToIssueResponses(issues)

		assert.Len(t, responses, 0)
		assert.NotNil(t, responses) // Should be empty slice, not nil
	})

	t.Run("handles nil slice input", func(t *testing.T) {
		var issues []models.Issue = nil

		responses := MapIssuesToIssueResponses(issues)

		assert.Len(t, responses, 0)
		assert.NotNil(t, responses) // Should create empty slice even from nil input
	})

	t.Run("handles Issues with empty names", func(t *testing.T) {
		issues := []models.Issue{
			{
				Model: gorm.Model{ID: 1},
				Name:  "",
			},
			{
				Model: gorm.Model{ID: 2},
				Name:  "Valid Name",
			},
		}

		responses := MapIssuesToIssueResponses(issues)

		assert.Len(t, responses, 2)
		assert.Equal(t, uint(1), responses[0].ID)
		assert.Empty(t, responses[0].Name)
		assert.Equal(t, uint(2), responses[1].ID)
		assert.Equal(t, "Valid Name", responses[1].Name)
	})

	t.Run("handles Issues with zero IDs", func(t *testing.T) {
		issues := []models.Issue{
			{
				Model: gorm.Model{ID: 0},
				Name:  "Zero ID Issue",
			},
		}

		responses := MapIssuesToIssueResponses(issues)

		assert.Len(t, responses, 1)
		assert.Equal(t, uint(0), responses[0].ID)
		assert.Equal(t, "Zero ID Issue", responses[0].Name)
	})
}

func TestReportResponse_Structure(t *testing.T) {
	t.Run("ReportResponse embeds ReportRequest correctly", func(t *testing.T) {
		testUUID := uuid.New()
		satisfied := true

		response := ReportResponse{
			ReportRequest: ReportRequest{
				UUID:      testUUID,
				Satisfied: &satisfied,
				Comment:   "Test comment",
			},
		}

		// Verify that ReportResponse has access to all ReportRequest fields
		assert.Equal(t, testUUID, response.UUID)
		assert.NotNil(t, response.Satisfied)
		assert.True(t, *response.Satisfied)
		assert.Equal(t, "Test comment", response.Comment)
	})
}

func TestIssueResponse_Structure(t *testing.T) {
	t.Run("IssueResponse contains correct fields", func(t *testing.T) {
		response := IssueResponse{
			ID:   123,
			Name: "Test Issue",
		}

		assert.Equal(t, uint(123), response.ID)
		assert.Equal(t, "Test Issue", response.Name)
	})

	t.Run("IssueResponse handles zero values", func(t *testing.T) {
		response := IssueResponse{
			ID:   0,
			Name: "",
		}

		assert.Equal(t, uint(0), response.ID)
		assert.Empty(t, response.Name)
	})
}
