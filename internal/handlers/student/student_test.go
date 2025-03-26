package student

import (
	"Group03-EX-StudentManagementAppBE/common"
	"Group03-EX-StudentManagementAppBE/internal/handlers/base"
	"Group03-EX-StudentManagementAppBE/internal/models"
	adminModels "Group03-EX-StudentManagementAppBE/internal/models/admin"
	studentModels "Group03-EX-StudentManagementAppBE/internal/models/student"
	studentStatusModels "Group03-EX-StudentManagementAppBE/internal/models/student_status"
	"Group03-EX-StudentManagementAppBE/internal/services"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockStudentService is a mock implementation of the student service
type MockStudentService struct {
	mock.Mock
}

func (m *MockStudentService) GetStudentByID(ctx context.Context, id string) (*studentModels.StudentResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*studentModels.StudentResponse), args.Error(1)
}

func (m *MockStudentService) GetStudentList(ctx context.Context, req *studentModels.ListStudentRequest) (*models.BaseListResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.BaseListResponse), args.Error(1)
}

func (m *MockStudentService) CreateAStudent(ctx context.Context, userID string, req *studentModels.CreateStudentRequest) error {
	args := m.Called(ctx, userID, req)
	return args.Error(0)
}

func (m *MockStudentService) UpdateStudent(ctx context.Context, userID string, studentID string, req *studentModels.UpdateStudentRequest) error {
	args := m.Called(ctx, userID, studentID, req)
	return args.Error(0)
}

func (m *MockStudentService) DeleteStudentByID(ctx context.Context, userID string, studentID string) error {
	args := m.Called(ctx, userID, studentID)
	return args.Error(0)
}

func (m *MockStudentService) GetStudentStatuses(ctx context.Context, req *studentStatusModels.ListStudentStatusRequest) ([]*studentModels.StudentStatus, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*studentModels.StudentStatus), args.Error(1)
}

func (m *MockStudentService) CreateStudentStatus(ctx context.Context, studentStatus *studentStatusModels.CreateStudentStatusRequest) error {
	args := m.Called(ctx, studentStatus)
	return args.Error(0)
}

func (m *MockStudentService) UpdateStudentStatus(ctx context.Context, id string, req *studentStatusModels.UpdateStudentStatusRequest) (*studentModels.StudentStatus, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*studentModels.StudentStatus), args.Error(1)
}

func (m *MockStudentService) DeleteStudentStatus(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockStudentService) ImportStudentsFromFile(ctx context.Context, userID string, fileURL string) (*adminModels.ImportResult, error) {
	args := m.Called(ctx, userID, fileURL)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*adminModels.ImportResult), args.Error(1)
}

func (m *MockStudentService) ExportStudentsToCSV(ctx context.Context) ([]byte, error) {
	args := m.Called(ctx)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockStudentService) ExportStudentsToJSON(ctx context.Context) ([]byte, error) {
	args := m.Called(ctx)
	return args.Get(0).([]byte), args.Error(1)
}

// Mock the JWT profile middleware
func mockJWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create the profile with the correct structure
		profile := &common.UserJWTProfile{
			Id:          "8a0f7a89-cac7-48b3-8f6e-cdb1786fa953",
			Role:        "a1b2c3d4-e5f6-47a8-b9c0-d1e2f3a4b5c6",
			AppAccess:   true,
			AdminAccess: true,
			Iat:         1615000000,
			Exp:         time.Now().Add(24 * time.Hour).Unix(), // Make sure it's not expired
			Iss:         "test-issuer",
		}

		// Set it in the context with the EXACT same key as the actual middleware
		c.Set(common.USER_JWT_KEY, profile)

		// Also set the user ID since your middleware does this too
		c.Set(common.UserId, profile.Id)

		c.Next()
	}
}

func createMockStudentService() *MockStudentService {
	return new(MockStudentService)
}

func setupTestRouter(mockStudentService *MockStudentService) (*gin.Engine, *Handler) {
	gin.SetMode(gin.TestMode)

	service := &services.Service{
		Student: mockStudentService,
	}

	handler := &Handler{
		Handler: base.NewHandler(service),
	}

	router := gin.Default()
	return router, handler
}

func TestGetStudentList_Success(t *testing.T) {
	mockService := createMockStudentService()
	router, handler := setupTestRouter(mockService)

	// Mock JWT authentication middleware
	router.GET("/students", mockJWTMiddleware(), handler.GetStudentList)

	// Create mock student list response that matches the actual API response structure
	mockStudentResponses := []*studentModels.StudentListResponse{
		{
			ID:          uuid.New(),
			Fullname:    "John Doe",
			StudentCode: 20083,
			Email:       "johndoe@student.com",
			FacultyID:   1,
			Gender:      "Male",
		},
		{
			ID:          uuid.New(),
			Fullname:    "Jane Smith",
			StudentCode: 20086,
			Email:       "janesmith@student.com",
			FacultyID:   2,
			Gender:      "Female",
		},
	}

	// Create the expected response that matches your API format
	expectedResponse := &models.BaseListResponse{
		Total:    2,
		Page:     1,
		PageSize: 10,
		Items:    mockStudentResponses,
		Extra:    nil,
	}

	// Setup the mock to return the expected response for any arguments
	mockService.On("GetStudentList", mock.Anything, mock.Anything).Return(expectedResponse, nil)

	// Create the request with query parameters
	req, _ := http.NewRequest(http.MethodGet, "/students?page=1&page_size=10", nil)
	w := httptest.NewRecorder()

	// Serve the request
	router.ServeHTTP(w, req)

	// Check response status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse the response
	var response models.BaseListResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verify the response structure
	assert.Equal(t, expectedResponse.Total, response.Total)
	assert.Equal(t, expectedResponse.Page, response.Page)
	assert.Equal(t, expectedResponse.PageSize, response.PageSize)
	assert.NotNil(t, response.Items)
}

func TestGetStudentList_Unauthorized(t *testing.T) {
	router := gin.Default()

	// Use the actual handler without middleware
	handler := &Handler{
		Handler: base.NewHandler(&services.Service{}),
	}

	router.GET("/students", handler.GetStudentList) // No JWT middleware

	req, _ := http.NewRequest(http.MethodGet, "/students", nil)
	w := httptest.NewRecorder()

	// Serve the request
	router.ServeHTTP(w, req)

	// Should return unauthorized (your implementation returns 400)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetStudentList_InvalidInput(t *testing.T) {
	mockService := createMockStudentService()
	router, handler := setupTestRouter(mockService)

	// Mock JWT authentication middleware
	router.GET("/students", mockJWTMiddleware(), handler.GetStudentList)

	// Setup the mock to return a valid response
	// This is important - we need to setup the mock even though we expect validation to fail
	mockService.On("GetStudentList", mock.Anything, mock.MatchedBy(func(req *studentModels.ListStudentRequest) bool {
		return req.PageSize < 0
	})).Return(nil, common.ErrInvalidInput)

	// Using an invalid parameter that should trigger validation error
	req, _ := http.NewRequest(http.MethodGet, "/students?page_size=-1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Check for bad request
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetStudentList_ServiceError(t *testing.T) {
	mockService := createMockStudentService()
	router, handler := setupTestRouter(mockService)

	// Mock JWT authentication middleware
	router.GET("/students", mockJWTMiddleware(), handler.GetStudentList)

	// Setup the mock to return an error
	mockService.On("GetStudentList", mock.Anything, mock.Anything).Return(nil, errors.New("service error"))

	req, _ := http.NewRequest(http.MethodGet, "/students?page=1&page_size=10", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Your handler returns 400 for service errors rather than 500
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
