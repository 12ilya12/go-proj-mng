package repos

import (
	"math"
	"strings"

	"github.com/12ilya12/go-proj-mng/common"
	"github.com/12ilya12/go-proj-mng/models"
	"github.com/12ilya12/go-proj-mng/pagination"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	DB *gorm.DB
}

func NewCategoryRepository(DB *gorm.DB) CategoryRepository {
	return CategoryRepository{DB}
}

func (sr *CategoryRepository) GetAll(pagingOptions pagination.PagingOptions) (categoriesWithPaging pagination.Paging[models.Category] /* categories []models.Category */, err error) {
	//Сортировка. По умолчанию по возрастанию идентификатора.
	var orderRule string
	if pagingOptions.OrderBy == "" {
		orderRule = "id"
	} else {
		var columnCount int64
		sr.DB.Select("column_name").Table("information_schema.columns").
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
	tx := sr.DB.Order(orderRule)

	//Пагинация
	if pagingOptions.PageSize > 0 {
		tx = tx.Limit(pagingOptions.PageSize)
	}
	if pagingOptions.Page > 0 {
		tx = tx.Offset((pagingOptions.Page - 1) * pagingOptions.PageSize)
	}

	err = tx.Find(&categoriesWithPaging.Items).Error

	//Собираем выходные данные пагинации
	sr.DB.Model(&models.Category{}).Count(&categoriesWithPaging.Pagination.TotalItems)
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

func (sr *CategoryRepository) GetById(id uint) (category models.Category, err error) {
	err = sr.DB.First(&category, id).Error
	return
}

func (sr *CategoryRepository) Create(category *models.Category) (err error) {
	err = sr.DB.Create(&category).Error
	return
}

func (sr *CategoryRepository) Update(paramsForUpdate *models.Category) (updatedCategory models.Category, err error) {
	updatedCategory = models.Category{}
	err = sr.DB.First(&updatedCategory, paramsForUpdate.ID).Error
	if err != nil {
		//Не найдена категория с заданным идентификатором, либо другая проблема с БД
		return
	}
	updatedCategory.Name = paramsForUpdate.Name
	err = sr.DB.Save(&updatedCategory).Error
	return
}

func (sr *CategoryRepository) HasTasks(id uint) (hasTasks bool, err error) {
	//Проверка наличия категории с заданным идентификатором
	err = sr.DB.First(&models.Category{}, id).Error
	if err != nil {
		//Не найдена категория с заданным идентификатором, либо другая проблема с БД
		return
	}

	var tasksWithCategoryCount int64
	sr.DB.Table("tasks").Where("category_id = ?", id).Count(&tasksWithCategoryCount)
	hasTasks = tasksWithCategoryCount > 0
	return
}

func (sr *CategoryRepository) Delete(id uint) (err error) {
	hasTasks, err := sr.HasTasks(id)
	if err != nil {
		return
	}

	if hasTasks { //Если удаляемая категория имеет связи с задачами, возвращаем ошибку и удаление не производим
		err = common.ErrCategoryHasRelatedTasks
		return
	}

	err = sr.DB.Delete(&models.Category{}, id).Error
	return
}

func (sr *CategoryRepository) DeleteForce(id uint) (err error) {
	//Проверка наличия категории с заданным идентификатором
	category := models.Category{}
	err = sr.DB.First(&category, id).Error
	if err != nil {
		//Не найдена категория с заданным идентификатором, либо другая проблема с БД
		return
	}

	//TODO: Проверить отработает ли хук для каскадного удаления задач, связанных с удаляемой категорией
	err = sr.DB.Delete(&category, id).Error
	return
}
