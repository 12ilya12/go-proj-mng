package services

import (
	"github.com/12ilya12/go-proj-mng/common"
	"github.com/12ilya12/go-proj-mng/models"
	"github.com/12ilya12/go-proj-mng/pagination"
	"github.com/12ilya12/go-proj-mng/repos"
)

type CategoryService interface {
	GetAll(pagingOptions pagination.PagingOptions) (statusesWithPaging pagination.Paging[models.Category], err error)
	GetById(id uint) (status models.Category, err error)
	Create(status *models.Category) (err error)
	Update(paramsForUpdate *models.Category) (updatedCategory models.Category, err error)
	HasTasks(id uint) (hasTasks bool, err error)
	Delete(id uint) (err error)
	DeleteForce(id uint) (err error)
}

type CategoryServiceImpl struct {
	categoryRepo repos.CategoryRepository
}

func NewCategoryServiceImpl(categoryRepo repos.CategoryRepository) CategoryService {
	return &CategoryServiceImpl{categoryRepo}
}

func (ss *CategoryServiceImpl) GetAll(pagingOptions pagination.PagingOptions) (categoriesWithPaging pagination.Paging[models.Category], err error) {
	categoriesWithPaging, err = ss.categoryRepo.GetAll(pagingOptions)
	return
}

func (ss *CategoryServiceImpl) GetById(id uint) (category models.Category, err error) {
	category, err = ss.categoryRepo.GetById(id)
	return
}

func (ss *CategoryServiceImpl) Create(category *models.Category) (err error) {
	err = ss.categoryRepo.Create(category)
	return
}

func (ss *CategoryServiceImpl) Update(paramsForUpdate *models.Category) (updatedCategory models.Category, err error) {
	updatedCategory, err = ss.categoryRepo.Update(paramsForUpdate)
	return
}

func (ss *CategoryServiceImpl) HasTasks(id uint) (hasTasks bool, err error) {
	hasTasks, err = ss.categoryRepo.HasTasks(id)
	return
}

func (ss *CategoryServiceImpl) Delete(id uint) (err error) {
	hasTasks, err := ss.HasTasks(id)
	if err != nil {
		return
	}

	if hasTasks { //Если удаляемая категория имеет связи с задачами, возвращаем ошибку и удаление не производим
		err = common.ErrCategoryHasRelatedTasks
		return
	}
	err = ss.categoryRepo.Delete(id)
	return
}

func (ss *CategoryServiceImpl) DeleteForce(id uint) (err error) {
	err = ss.categoryRepo.Delete(id)
	return
}
