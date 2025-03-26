package student

import (
	models "Group03-EX-StudentManagementAppBE/internal/models/student"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock Student Service
type MockStudentService struct {
	mock.Mock
}

func (m *MockStudentService) CreateAStudent(ctx context.Context, userID string, req *models.CreateStudentRequest) error {
	args := m.Called(ctx, userID, req)
	return args.Error(0)
}

// Test for createStudentRequestFromJSON
func TestCreateStudentRequestFromJSON(t *testing.T) {
	service := &studentService{}

	// Test case: Valid JSON input
	t.Run("Valid JSON Input", func(t *testing.T) {
		jsonData := map[string]interface{}{
			"studentCode": 12345,
			"fullname":    "John Doe",
			"email":       "john.doe@example.com",
			"dateOfBirth": "2000-01-01",
			"gender":      "Male",
			"facultyId":   1,
			"batch":       "2023",
			"program":     "Computer Science",
			"address":     "123 Main St",
			"phone":       "1234567890",
			"statusId":    1,
			"programId":   1,
			"nationality": "Vietnam",
		}

		req, err := service.createStudentRequestFromJSON(jsonData)
		assert.NoError(t, err)
		assert.NotNil(t, req)
		assert.Equal(t, 12345, *req.StudentCode)
		assert.Equal(t, "John Doe", *req.Fullname)
		assert.Equal(t, "john.doe@example.com", *req.Email)
		assert.Equal(t, "Male", *req.Gender)
	})

	// Test case: Missing required fields
	t.Run("Missing Required Fields", func(t *testing.T) {
		jsonData := map[string]interface{}{
			"fullname": "John Doe",
		}

		req, err := service.createStudentRequestFromJSON(jsonData)
		assert.Error(t, err)
		assert.Nil(t, req)
	})

	// Test case: Invalid date format
	t.Run("Invalid Date Format", func(t *testing.T) {
		jsonData := map[string]interface{}{
			"studentCode": 12345,
			"fullname":    "John Doe",
			"email":       "john.doe@example.com",
			"dateOfBirth": "invalid-date",
		}

		req, err := service.createStudentRequestFromJSON(jsonData)
		assert.NoError(t, err)
		assert.NotNil(t, req)
		assert.True(t, req.DateOfBirth.Before(time.Now()))
	})

	// Test case: Default values for optional fields
	t.Run("Default Values", func(t *testing.T) {
		jsonData := map[string]interface{}{
			"studentCode": 12345,
			"fullname":    "John Doe",
			"email":       "john.doe@example.com",
		}

		req, err := service.createStudentRequestFromJSON(jsonData)
		assert.NoError(t, err)
		assert.NotNil(t, req)
		assert.Equal(t, "Other", *req.Gender)
		assert.Equal(t, "Vietnam", *req.Nationality)
	})
}
