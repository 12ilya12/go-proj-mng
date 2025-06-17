package services

import (
	"github.com/12ilya12/go-proj-mng/models"
	"github.com/12ilya12/go-proj-mng/pagination"
	"github.com/12ilya12/go-proj-mng/repos"
)

type CategoryService struct {
	categoryRepo repos.CategoryRepository
}

func NewCategoryService(categoryRepo repos.CategoryRepository) CategoryService {
	return CategoryService{categoryRepo}
}

func (ss *CategoryService) GetAll(pagingOptions pagination.PagingOptions) (categoriesWithPaging pagination.Paging[models.Category], err error) {
	categoriesWithPaging, err = ss.categoryRepo.GetAll(pagingOptions)
	return
}

func (ss *CategoryService) GetById(id uint) (category models.Category, err error) {
	category, err = ss.categoryRepo.GetById(id)
	return
}

func (ss *CategoryService) Create(category *models.Category) (err error) {
	err = ss.categoryRepo.Create(category)
	return
}

func (ss *CategoryService) Update(paramsForUpdate *models.Category) (updatedCategory models.Category, err error) {
	updatedCategory, err = ss.categoryRepo.Update(paramsForUpdate)
	return
}

func (ss *CategoryService) HasTasks(id uint) (hasTasks bool, err error) {
	hasTasks, err = ss.categoryRepo.HasTasks(id)
	return
}

func (ss *CategoryService) Delete(id uint) (err error) {
	err = ss.categoryRepo.Delete(id)
	return
}

func (ss *CategoryService) DeleteForce(id uint) (err error) {
	err = ss.categoryRepo.DeleteForce(id)
	return
}
