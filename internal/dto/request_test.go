package dto

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/datatypes"
)

func TestReportRequest_JSONMarshaling(t *testing.T) {
	t.Run("marshals valid ReportRequest to JSON", func(t *testing.T) {
		testUUID := uuid.New()
		satisfied := true
		issueID := 42
		metadata := datatypes.JSON(`{"key": "value"}`)

		req := ReportRequest{
			UUID:      testUUID,
			Satisfied: &satisfied,
			Comment:   "Test comment",
			IssueID:   &issueID,
			Metadata:  &metadata,
		}

		jsonData, err := json.Marshal(req)
		require.NoError(t, err)
		assert.Contains(t, string(jsonData), testUUID.String())
		assert.Contains(t, string(jsonData), "true")
		assert.Contains(t, string(jsonData), "Test comment")
		assert.Contains(t, string(jsonData), "42")
	})

	t.Run("marshals ReportRequest with minimal fields", func(t *testing.T) {
		testUUID := uuid.New()
		satisfied := false

		req := ReportRequest{
			UUID:      testUUID,
			Satisfied: &satisfied,
		}

		jsonData, err := json.Marshal(req)
		require.NoError(t, err)
		assert.Contains(t, string(jsonData), testUUID.String())
		assert.Contains(t, string(jsonData), "false")
	})
}

func TestReportRequest_JSONUnmarshaling_ValidCases(t *testing.T) {
	t.Run("unmarshals valid JSON to ReportRequest", func(t *testing.T) {
		testUUID := uuid.New()
		jsonData := `{
			"uuid": "` + testUUID.String() + `",
			"satisfied": true,
			"comment": "Great service!",
			"issue_id": 123,
			"metadata": {"source": "web", "version": "1.0"}
		}`

		var req ReportRequest
		err := json.Unmarshal([]byte(jsonData), &req)
		require.NoError(t, err)

		assert.Equal(t, testUUID, req.UUID)
		assert.NotNil(t, req.Satisfied)
		assert.True(t, *req.Satisfied)
		assert.Equal(t, "Great service!", req.Comment)
		assert.NotNil(t, req.IssueID)
		assert.Equal(t, 123, *req.IssueID)
		assert.NotNil(t, req.Metadata)
	})

	t.Run("unmarshals JSON with only required fields", func(t *testing.T) {
		testUUID := uuid.New()
		jsonData := `{
			"uuid": "` + testUUID.String() + `",
			"satisfied": false
		}`

		var req ReportRequest
		err := json.Unmarshal([]byte(jsonData), &req)
		require.NoError(t, err)

		assert.Equal(t, testUUID, req.UUID)
		assert.NotNil(t, req.Satisfied)
		assert.False(t, *req.Satisfied)
		assert.Empty(t, req.Comment)
		assert.Nil(t, req.IssueID)
		assert.Nil(t, req.Metadata)
	})
}

func TestReportRequest_JSONUnmarshaling_EdgeCases(t *testing.T) {
	t.Run("handles null values for optional fields", func(t *testing.T) {
		testUUID := uuid.New()
		jsonData := `{
			"uuid": "` + testUUID.String() + `",
			"satisfied": true,
			"comment": "",
			"issue_id": null,
			"metadata": null
		}`

		var req ReportRequest
		err := json.Unmarshal([]byte(jsonData), &req)
		require.NoError(t, err)

		assert.Equal(t, testUUID, req.UUID)
		assert.NotNil(t, req.Satisfied)
		assert.True(t, *req.Satisfied)
		assert.Empty(t, req.Comment)
		assert.Nil(t, req.IssueID)
		assert.Nil(t, req.Metadata)
	})

	t.Run("fails with invalid UUID", func(t *testing.T) {
		jsonData := `{
			"uuid": "invalid-uuid",
			"satisfied": true
		}`

		var req ReportRequest
		err := json.Unmarshal([]byte(jsonData), &req)
		assert.Error(t, err)
	})

	t.Run("fails with invalid JSON", func(t *testing.T) {
		jsonData := `{
			"uuid": "` + uuid.New().String() + `",
			"satisfied": "not-a-boolean"
		}`

		var req ReportRequest
		err := json.Unmarshal([]byte(jsonData), &req)
		assert.Error(t, err)
	})
}

func TestReportRequest_FieldValidation_BasicCases(t *testing.T) {
	t.Run("creates valid ReportRequest with all fields", func(t *testing.T) {
		testUUID := uuid.New()
		satisfied := true
		issueID := 1
		metadata := datatypes.JSON(`{"test": true}`)

		req := ReportRequest{
			UUID:      testUUID,
			Satisfied: &satisfied,
			Comment:   "Valid comment",
			IssueID:   &issueID,
			Metadata:  &metadata,
		}

		assert.Equal(t, testUUID, req.UUID)
		assert.NotNil(t, req.Satisfied)
		assert.True(t, *req.Satisfied)
		assert.Equal(t, "Valid comment", req.Comment)
		assert.NotNil(t, req.IssueID)
		assert.Equal(t, 1, *req.IssueID)
		assert.NotNil(t, req.Metadata)
	})

	t.Run("creates valid ReportRequest with minimal required fields", func(t *testing.T) {
		testUUID := uuid.New()
		satisfied := false

		req := ReportRequest{
			UUID:      testUUID,
			Satisfied: &satisfied,
		}

		assert.Equal(t, testUUID, req.UUID)
		assert.NotNil(t, req.Satisfied)
		assert.False(t, *req.Satisfied)
		assert.Empty(t, req.Comment)
		assert.Nil(t, req.IssueID)
		assert.Nil(t, req.Metadata)
	})
}

func TestReportRequest_FieldValidation_PointerFields(t *testing.T) {
	t.Run("handles pointer fields correctly", func(t *testing.T) {
		testUUID := uuid.New()
		satisfied := true

		// Test that Satisfied field being a pointer works correctly
		req := ReportRequest{
			UUID:      testUUID,
			Satisfied: &satisfied,
		}

		// Verify we can change the value through the pointer
		*req.Satisfied = false
		assert.False(t, *req.Satisfied)
		assert.False(t, satisfied) // Original variable should also be changed

		// Test IssueID pointer behavior
		issueID := 42
		req.IssueID = &issueID
		assert.NotNil(t, req.IssueID)
		assert.Equal(t, 42, *req.IssueID)
	})
}

func TestReportRequest_EdgeCases(t *testing.T) {
	t.Run("handles empty comment", func(t *testing.T) {
		testUUID := uuid.New()
		satisfied := true

		req := ReportRequest{
			UUID:      testUUID,
			Satisfied: &satisfied,
			Comment:   "",
		}

		assert.Empty(t, req.Comment)
	})

	t.Run("handles zero issue ID", func(t *testing.T) {
		testUUID := uuid.New()
		satisfied := true
		issueID := 0

		req := ReportRequest{
			UUID:      testUUID,
			Satisfied: &satisfied,
			IssueID:   &issueID,
		}

		assert.NotNil(t, req.IssueID)
		assert.Equal(t, 0, *req.IssueID)
	})

	t.Run("handles empty metadata JSON", func(t *testing.T) {
		testUUID := uuid.New()
		satisfied := true
		metadata := datatypes.JSON(`{}`)

		req := ReportRequest{
			UUID:      testUUID,
			Satisfied: &satisfied,
			Metadata:  &metadata,
		}

		assert.NotNil(t, req.Metadata)
		assert.Equal(t, datatypes.JSON(`{}`), *req.Metadata)
	})
}
