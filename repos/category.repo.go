package repos

import (
	"math"
	"strings"

	"github.com/12ilya12/go-proj-mng/models"
	"github.com/12ilya12/go-proj-mng/pagination"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	GetAll(pagingOptions pagination.PagingOptions) (pagination.Paging[models.Category], error)
	GetById(id uint) (models.Category, error)
	Create(status *models.Category) error
	Update(paramsForUpdate *models.Category) (models.Category, error)
	HasTasks(statusId uint) (bool, error)
	Delete(id uint) error
}

type CategoryRepositoryImpl struct {
	DB *gorm.DB
}

func NewCategoryRepositoryImpl(DB *gorm.DB) CategoryRepository {
	return &CategoryRepositoryImpl{DB}
}

func (cr *CategoryRepositoryImpl) GetAll(pagingOptions pagination.PagingOptions) (categoriesWithPaging pagination.Paging[models.Category], err error) {
	//Сортировка. По умолчанию по возрастанию идентификатора.
	var orderRule string
	if pagingOptions.OrderBy == "" {
		orderRule = "id"
	} else {
		var columnCount int64
		cr.DB.Select("column_name").Table("information_schema.columns").
			Where("table_name = ? AND column_name = ?", "categories", pagingOptions.OrderBy).Count(&columnCount)
		if columnCount == 0 {
			//Колонка, по которой планировалось сортировать, отсутствует в таблице
			pagingOptions.OrderBy = ""
			orderRule = "id"
		} else {
			orderRule = pagingOptions.OrderBy
		}
	}
	if strings.ToLower(string(pagingOptions.Order)) == "desc" {
		orderRule += " desc"
	}
	tx := cr.DB.Order(orderRule)

	//Пагинация
	if pagingOptions.PageSize > 0 {
		tx = tx.Limit(pagingOptions.PageSize)
	}
	if pagingOptions.Page > 0 {
		tx = tx.Offset((pagingOptions.Page - 1) * pagingOptions.PageSize)
	}

	err = tx.Find(&categoriesWithPaging.Items).Error

	//Собираем выходные данные пагинации
	cr.DB.Model(&models.Category{}).Count(&categoriesWithPaging.Pagination.TotalItems)
	if pagingOptions.PageSize == 0 { //Если размер страницы не задан, показываем всё на одной странице
		categoriesWithPaging.Pagination.TotalPages = 1
	} else { //Подсчитываем количество страниц
		categoriesWithPaging.Pagination.TotalPages =
			int64(math.Ceil(float64(categoriesWithPaging.Pagination.TotalItems) /
				float64(pagingOptions.PageSize)))
	}
	categoriesWithPaging.Pagination.Options = pagingOptions

	return
}

func (cr *CategoryRepositoryImpl) GetById(id uint) (category models.Category, err error) {
	err = cr.DB.First(&category, id).Error
	return
}

func (cr *CategoryRepositoryImpl) Create(category *models.Category) (err error) {
	err = cr.DB.Create(&category).Error
	return
}

func (cr *CategoryRepositoryImpl) Update(paramsForUpdate *models.Category) (updatedCategory models.Category, err error) {
	updatedCategory = models.Category{}
	err = cr.DB.First(&updatedCategory, paramsForUpdate.ID).Error
	if err != nil {
		//Не найдена категория с заданным идентификатором, либо другая проблема с БД
		return
	}
	updatedCategory.Name = paramsForUpdate.Name
	err = cr.DB.Save(&updatedCategory).Error
	return
}

func (cr *CategoryRepositoryImpl) HasTasks(id uint) (hasTasks bool, err error) {
	//Проверка наличия категории с заданным идентификатором
	err = cr.DB.First(&models.Category{}, id).Error
	if err != nil {
		//Не найдена категория с заданным идентификатором, либо другая проблема с БД
		return
	}

	var tasksWithCategoryCount int64
	cr.DB.Table("tasks").Where("category_id = ?", id).Count(&tasksWithCategoryCount)
	hasTasks = tasksWithCategoryCount > 0
	return
}

func (cr *CategoryRepositoryImpl) Delete(id uint) (err error) {
	//Проверка наличия категории с заданным идентификатором
	category := models.Category{}
	err = cr.DB.First(&category, id).Error
	if err != nil {
		//Не найдена категория с заданным идентификатором, либо другая проблема с БД
		return
	}

	err = cr.DB.Delete(&category, id).Error
	return
}
