package tests

import (
	"github.com/12ilya12/go-proj-mng/models"
	"github.com/12ilya12/go-proj-mng/pagination"
	"github.com/12ilya12/go-proj-mng/services"

	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockStatusRepo - мок репозитория статусов
type MockStatusRepo struct {
	mock.Mock
}

func (m *MockStatusRepo) GetAll(pagingOptions pagination.PagingOptions) (pagination.Paging[models.Status], error) {
	args := m.Called()
	return args.Get(0).(pagination.Paging[models.Status]), args.Error(1)
}

func (m *MockStatusRepo) GetById(id uint) (models.Status, error) {
	args := m.Called(id)
	return args.Get(0).(models.Status), args.Error(1)
}

func (m *MockStatusRepo) Create(status *models.Status) error {
	args := m.Called(status)
	return args.Error(0)
}

func (m *MockStatusRepo) Update(status *models.Status) (models.Status, error) {
	args := m.Called(status)
	return args.Get(0).(models.Status), args.Error(1)
}

func (m *MockStatusRepo) HasTasks(statusId uint) bool {
	args := m.Called(statusId)
	return args.Get(0).(bool)
}

func (m *MockStatusRepo) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestStatusService_GetAll(t *testing.T) {
	mockRepo := new(MockStatusRepo)
	service := services.NewStatusServiceImpl(mockRepo)

	expectedStatuses := pagination.Paging[models.Status]{
		Items: []models.Status{
			{ID: 1, Name: "Active"},
			{ID: 2, Name: "Resolved"},
		},
		Pagination: pagination.Pagination{},
	}

	mockRepo.On("GetAll").Return(expectedStatuses, nil)

	statuses, err := service.GetAll(pagination.PagingOptions{})

	assert.NoError(t, err)
	assert.Equal(t, expectedStatuses, statuses)
	mockRepo.AssertExpectations(t)
}

func TestStatusService_GetById(t *testing.T) {
	mockRepo := new(MockStatusRepo)
	service := services.NewStatusServiceImpl(mockRepo)

	expectedStatus := models.Status{ID: 1, Name: "Active"}
	mockRepo.On("GetById", uint(1)).Return(expectedStatus, nil)

	status, err := service.GetById(1)

	assert.NoError(t, err)
	assert.Equal(t, expectedStatus, status)
	mockRepo.AssertExpectations(t)
}

func TestStatusService_Create(t *testing.T) {
	mockRepo := new(MockStatusRepo)
	service := services.NewStatusServiceImpl(mockRepo)

	newStatus := &models.Status{Name: "New Status"}

	mockRepo.On("Create", newStatus).Return(nil)

	err := service.Create(newStatus)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestStatusService_Update(t *testing.T) {
	mockRepo := new(MockStatusRepo)
	service := services.NewStatusServiceImpl(mockRepo)

	updateData := &models.Status{ID: 1, Name: "Updated Status"}
	updatedStatus := models.Status{ID: 1, Name: "Updated Status"}

	mockRepo.On("Update", updateData).Return(updatedStatus, nil)

	result, err := service.Update(updateData)

	assert.NoError(t, err)
	assert.Equal(t, updatedStatus, result)
	mockRepo.AssertExpectations(t)
}

func TestStatusService_Delete(t *testing.T) {
	mockRepo := new(MockStatusRepo)
	service := services.NewStatusServiceImpl(mockRepo)

	id := uint(1)
	expectedStatus := models.Status{ID: 1, Name: "Active"}
	mockRepo.On("GetById", id).Return(expectedStatus, nil)
	mockRepo.On("HasTasks", id).Return(false)
	mockRepo.On("Delete", id).Return(nil)

	err := service.Delete(id)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
