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

type MockCategoryService struct {
	mock.Mock
}

func (m *MockCategoryService) GetAll(pagingOptions pagination.PagingOptions) (pagination.Paging[models.Category], error) {
	args := m.Called(pagingOptions)
	return args.Get(0).(pagination.Paging[models.Category]), args.Error(1)
}

func (m *MockCategoryService) GetById(id uint) (models.Category, error) {
	args := m.Called(id)
	return args.Get(0).(models.Category), args.Error(1)
}

func (m *MockCategoryService) Create(category *models.Category) error {
	args := m.Called(category)
	return args.Error(0)
}

func (m *MockCategoryService) Update(category *models.Category) (models.Category, error) {
	args := m.Called(category)
	return args.Get(0).(models.Category), args.Error(1)
}

func (m *MockCategoryService) HasTasks(id uint) (bool, error) {
	args := m.Called(id)
	return args.Get(0).(bool), args.Error(1)
}

func (m *MockCategoryService) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockCategoryService) DeleteForce(id uint) (err error) {
	args := m.Called(id)
	return args.Error(0)
}

func TestCategoryController_GetAll(t *testing.T) {
	mockService := new(MockCategoryService)
	controller := controllers.NewCategoryController(mockService)

	expectedCategories := pagination.Paging[models.Category]{
		Items: []models.Category{
			{ID: 1, Name: "Bug"},
			{ID: 2, Name: "Test Case"},
		},
		Pagination: pagination.Pagination{},
	}

	mockService.On("GetAll", pagination.PagingOptions{}).Return(expectedCategories, nil)

	req, err := http.NewRequest("GET", "/categories", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/categories", controller.GetAll).Methods("GET")

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response pagination.Paging[models.Category]
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedCategories, response)
	mockService.AssertExpectations(t)
}

func TestCategoryController_GetById(t *testing.T) {
	mockService := new(MockCategoryService)
	controller := controllers.NewCategoryController(mockService)

	expectedCategory := models.Category{ID: 1, Name: "Bug"}

	mockService.On("GetById", uint(1)).Return(expectedCategory, nil)

	req, err := http.NewRequest("GET", "/categories/1", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/categories/{id}", controller.GetById).Methods("GET")

	// Set URL variables
	vars := map[string]string{
		"id": "1",
	}
	req = mux.SetURLVars(req, vars)

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response models.Category
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedCategory, response)
	mockService.AssertExpectations(t)
}

func TestCategoryController_Create(t *testing.T) {
	mockService := new(MockCategoryService)
	controller := controllers.NewCategoryController(mockService)

	newCategory := models.Category{Name: "New Category"}
	requestBody, err := json.Marshal(newCategory)
	assert.NoError(t, err)

	mockService.On("Create", &models.Category{Name: "New Category"}).Return(nil)

	req, err := http.NewRequest("POST", "/categories", bytes.NewBuffer(requestBody))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/categories", controller.Create).Methods("POST")

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	mockService.AssertExpectations(t)
}

func TestCategoryController_Update(t *testing.T) {
	mockService := new(MockCategoryService)
	controller := controllers.NewCategoryController(mockService)

	updateData := models.Category{ID: 1, Name: "Updated Category"}
	expectedCategory := models.Category{ID: 1, Name: "Updated Category"}

	requestBody, err := json.Marshal(updateData)
	assert.NoError(t, err)

	mockService.On("Update", &models.Category{ID: 1, Name: "Updated Category"}).Return(expectedCategory, nil)

	req, err := http.NewRequest("PUT", "/categories/1", bytes.NewBuffer(requestBody))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/categories/{id}", controller.Update).Methods("PUT")

	// Set URL variables
	vars := map[string]string{
		"id": "1",
	}
	req = mux.SetURLVars(req, vars)

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response models.Category
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedCategory, response)
	mockService.AssertExpectations(t)
}

func TestCategoryController_Delete(t *testing.T) {
	mockService := new(MockCategoryService)
	controller := controllers.NewCategoryController(mockService)

	mockService.On("Delete", uint(1)).Return(nil)

	req, err := http.NewRequest("DELETE", "/categories/1", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/categories/{id}", controller.Delete).Methods("DELETE")

	// Set URL variables
	vars := map[string]string{
		"id": "1",
	}
	req = mux.SetURLVars(req, vars)

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockService.AssertExpectations(t)
}

// Error cases
func TestCategoryController_GetById_NotFound(t *testing.T) {
	mockService := new(MockCategoryService)
	controller := controllers.NewCategoryController(mockService)

	mockService.On("GetById", uint(1)).Return(models.Category{}, gorm.ErrRecordNotFound)

	req, err := http.NewRequest("GET", "/categories/1", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/categories/{id}", controller.GetById).Methods("GET")

	vars := map[string]string{"id": "1"}
	req = mux.SetURLVars(req, vars)

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	mockService.AssertExpectations(t)
}
