package tests

import (
	"testing"

	"github.com/12ilya12/go-proj-mng/common"
	"github.com/12ilya12/go-proj-mng/models"
	"github.com/12ilya12/go-proj-mng/pagination"
	"github.com/12ilya12/go-proj-mng/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCategoryRepo struct {
	mock.Mock
}

func (m *MockCategoryRepo) GetAll(pagingOptions pagination.PagingOptions) (pagination.Paging[models.Category], error) {
	args := m.Called(pagingOptions)
	return args.Get(0).(pagination.Paging[models.Category]), args.Error(1)
}

func (m *MockCategoryRepo) GetById(id uint) (models.Category, error) {
	args := m.Called(id)
	return args.Get(0).(models.Category), args.Error(1)
}

func (m *MockCategoryRepo) Create(category *models.Category) error {
	args := m.Called(category)
	return args.Error(0)
}

func (m *MockCategoryRepo) Update(category *models.Category) (models.Category, error) {
	args := m.Called(category)
	return args.Get(0).(models.Category), args.Error(1)
}

func (m *MockCategoryRepo) HasTasks(id uint) (bool, error) {
	args := m.Called(id)
	return args.Get(0).(bool), args.Error(1)
}

func (m *MockCategoryRepo) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestCategoryService_GetAll(t *testing.T) {
	mockRepo := new(MockCategoryRepo)
	categoryService := services.NewCategoryServiceImpl(mockRepo)

	pagingOptions := pagination.PagingOptions{}
	expectedCategories := pagination.Paging[models.Category]{
		Items: []models.Category{
			{ID: 1, Name: "Bug"},
			{ID: 2, Name: "Test Case"},
		},
		Pagination: pagination.Pagination{},
	}

	mockRepo.On("GetAll", pagingOptions).Return(expectedCategories, nil)

	categories, err := categoryService.GetAll(pagination.PagingOptions{})

	assert.NoError(t, err)
	assert.Equal(t, expectedCategories, categories)
	mockRepo.AssertExpectations(t)
}

func TestCategoryService_GetById(t *testing.T) {
	mockRepo := new(MockCategoryRepo)
	categoryService := services.NewCategoryServiceImpl(mockRepo)

	expectedCategory := models.Category{ID: 1, Name: "Bug"}

	mockRepo.On("GetById", uint(1)).Return(expectedCategory, nil)

	category, err := categoryService.GetById(1)

	assert.NoError(t, err)
	assert.Equal(t, expectedCategory, category)
	mockRepo.AssertExpectations(t)
}

func TestCategoryService_Create(t *testing.T) {
	mockRepo := new(MockCategoryRepo)
	categoryService := services.NewCategoryServiceImpl(mockRepo)

	newCategory := &models.Category{Name: "New Category"}

	mockRepo.On("Create", newCategory).Return(nil)

	err := categoryService.Create(newCategory)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCategoryService_Update(t *testing.T) {
	mockRepo := new(MockCategoryRepo)
	service := services.NewCategoryServiceImpl(mockRepo)

	updateData := &models.Category{ID: 1, Name: "Updated Category"}
	updatedCategory := models.Category{ID: 1, Name: "Updated Category"}

	mockRepo.On("Update", updateData).Return(updatedCategory, nil)

	result, err := service.Update(updateData)

	assert.NoError(t, err)
	assert.Equal(t, updatedCategory, result)
	mockRepo.AssertExpectations(t)
}

func TestCategoryService_Delete(t *testing.T) {
	mockRepo := new(MockCategoryRepo)
	categoryService := services.NewCategoryServiceImpl(mockRepo)

	id := uint(1)
	mockRepo.On("HasTasks", id).Return(false, nil)
	mockRepo.On("Delete", id).Return(nil)

	err := categoryService.Delete(id)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCategoryService_Delete_CategoryHasTasks(t *testing.T) {
	mockRepo := new(MockCategoryRepo)
	service := services.NewCategoryServiceImpl(mockRepo)

	id := uint(1)
	mockRepo.On("HasTasks", id).Return(true, nil)

	err := service.Delete(id)

	assert.Error(t, err)
	assert.Equal(t, common.ErrCategoryHasRelatedTasks, err)
	mockRepo.AssertExpectations(t)
}
