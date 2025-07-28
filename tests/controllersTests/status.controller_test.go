package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/12ilya12/go-proj-mng/controllers"
	"github.com/12ilya12/go-proj-mng/models"
	"github.com/12ilya12/go-proj-mng/pagination"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockStatusService struct {
	mock.Mock
}

func (m *MockStatusService) GetAll(pagingOptions pagination.PagingOptions) (pagination.Paging[models.Status], error) {
	args := m.Called(pagingOptions)
	return args.Get(0).(pagination.Paging[models.Status]), args.Error(1)
}

func (m *MockStatusService) GetById(id uint) (models.Status, error) {
	args := m.Called(id)
	return args.Get(0).(models.Status), args.Error(1)
}

func (m *MockStatusService) Create(status *models.Status) error {
	args := m.Called(status)
	return args.Error(0)
}

func (m *MockStatusService) Update(status *models.Status) (models.Status, error) {
	args := m.Called(status)
	return args.Get(0).(models.Status), args.Error(1)
}

func (m *MockStatusService) HasTasks(id uint) (bool, error) {
	args := m.Called(id)
	return args.Get(0).(bool), args.Error(1)
}

func (m *MockStatusService) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockStatusService) DeleteForce(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestStatusController_GetAll(t *testing.T) {
	mockService := new(MockStatusService)
	controller := controllers.NewStatusController(mockService)

	expectedStatuses := pagination.Paging[models.Status]{
		Items: []models.Status{
			{ID: 1, Name: "New"},
			{ID: 2, Name: "Active"},
		},
		Pagination: pagination.Pagination{},
	}

	mockService.On("GetAll", pagination.PagingOptions{}).Return(expectedStatuses, nil)

	req, err := http.NewRequest("GET", "/statuses", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/statuses", controller.GetAll).Methods("GET")

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response pagination.Paging[models.Status]
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedStatuses, response)
	mockService.AssertExpectations(t)
}

func TestStatusController_GetById(t *testing.T) {
	mockService := new(MockStatusService)
	controller := controllers.NewStatusController(mockService)

	expectedStatus := models.Status{ID: 1, Name: "New"}

	mockService.On("GetById", uint(1)).Return(expectedStatus, nil)

	req, err := http.NewRequest("GET", "/statuses/1", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/statuses/{id}", controller.GetById).Methods("GET")

	vars := map[string]string{
		"id": "1",
	}
	req = mux.SetURLVars(req, vars)

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response models.Status
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedStatus, response)
	mockService.AssertExpectations(t)
}

func TestStatusController_Create(t *testing.T) {
	mockService := new(MockStatusService)
	controller := controllers.NewStatusController(mockService)

	newStatus := models.Status{Name: "New Status"}
	requestBody, err := json.Marshal(newStatus)
	assert.NoError(t, err)

	mockService.On("Create", &models.Status{Name: "New Status"}).Return(nil)

	req, err := http.NewRequest("POST", "/statuses", bytes.NewBuffer(requestBody))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/statuses", controller.Create).Methods("POST")

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	mockService.AssertExpectations(t)
}

func TestStatusController_Update(t *testing.T) {
	mockService := new(MockStatusService)
	controller := controllers.NewStatusController(mockService)

	updateData := models.Status{ID: 1, Name: "Updated Status"}
	expectedStatus := models.Status{ID: 1, Name: "Updated Status"}

	requestBody, err := json.Marshal(updateData)
	assert.NoError(t, err)

	mockService.On("Update", &models.Status{ID: 1, Name: "Updated Status"}).Return(expectedStatus, nil)

	req, err := http.NewRequest("PUT", "/statuses/1", bytes.NewBuffer(requestBody))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/statuses/{id}", controller.Update).Methods("PUT")

	vars := map[string]string{
		"id": "1",
	}
	req = mux.SetURLVars(req, vars)

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response models.Status
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedStatus, response)
	mockService.AssertExpectations(t)
}

func TestStatusController_Delete(t *testing.T) {
	mockService := new(MockStatusService)
	controller := controllers.NewStatusController(mockService)

	mockService.On("Delete", uint(1)).Return(nil)

	req, err := http.NewRequest("DELETE", "/statuses/1", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/statuses/{id}", controller.Delete).Methods("DELETE")

	vars := map[string]string{
		"id": "1",
	}
	req = mux.SetURLVars(req, vars)

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockService.AssertExpectations(t)
}

func TestStatusController_GetById_NotFound(t *testing.T) {
	mockService := new(MockStatusService)
	controller := controllers.NewStatusController(mockService)

	mockService.On("GetById", uint(1)).Return(models.Status{}, gorm.ErrRecordNotFound)

	req, err := http.NewRequest("GET", "/statuses/1", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/statuses/{id}", controller.GetById).Methods("GET")

	vars := map[string]string{"id": "1"}
	req = mux.SetURLVars(req, vars)

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	mockService.AssertExpectations(t)
}
