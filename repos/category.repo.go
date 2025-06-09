package repos

import (
	"math"

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
	//Собираем данные для ответа в ручке с пагинацией
	sr.DB.Find(&categoriesWithPaging.Items)
	categoriesWithPaging.Pagination.TotalItems = len(categoriesWithPaging.Items)
	if pagingOptions.PageSize == 0 { //Если размер страницы не задан, показываем всё на одной странице
		categoriesWithPaging.Pagination.TotalPages = 1
	} else { //Подсчитваем количество страниц
		categoriesWithPaging.Pagination.TotalPages =
			int(math.Ceil(float64(categoriesWithPaging.Pagination.TotalItems) / float64(pagingOptions.PageSize)))
	}

	//Значения по умолчанию для pagingOptions
	if pagingOptions.Order != "desc" {
		pagingOptions.Order = "asc"
	}
	if pagingOptions.Page <= 0 {
		pagingOptions.Page = 1
	}
	if pagingOptions.PageSize <= 0 {
		pagingOptions.PageSize = categoriesWithPaging.Pagination.TotalItems
	}
	if pagingOptions.OrderBy == "" {
		pagingOptions.OrderBy = "id"
	}
	categoriesWithPaging.Pagination.Options = pagingOptions

	//Добываем выборку с учетом параметров пагинации
	err = sr.DB.Order(pagingOptions.OrderBy + " " + string(pagingOptions.Order)).
		Limit(pagingOptions.PageSize).
		Offset((pagingOptions.Page - 1) * pagingOptions.PageSize).
		Find(&categoriesWithPaging.Items).Error

	return
}

func (sr *CategoryRepository) GetById(id int) (category models.Category, err error) {
	err = sr.DB.First(&category, id).Error
	return
}

func (sr *CategoryRepository) Create(category *models.Category) (err error) {
	err = sr.DB.Create(&category).Error
	return
}

func (sr *CategoryRepository) Update(updatedCategory *models.Category) (err error) {
	category := models.Category{}
	err = sr.DB.First(&category, updatedCategory.ID).Error
	if err != nil {
		//Не найдена категория с заданным идентификатором, либо другая проблема с БД
		return
	}
	category.Name = updatedCategory.Name
	err = sr.DB.Save(&category).Error
	return
}

func (sr *CategoryRepository) HasTasks(id int) (hasTasks bool, err error) {
	//Проверка наличия категории с заданным идентификатором
	err = sr.DB.First(&models.Category{}, id).Error
	if err != nil {
		//Не найдена категория с заданным идентификатором, либо другая проблема с БД
		return
	}

	tasksWithCategoryCount := int64(0)
	sr.DB.Table("tasks").Where("category_id = ?", id).Count(&tasksWithCategoryCount)
	hasTasks = tasksWithCategoryCount > 0
	return
}

func (sr *CategoryRepository) Delete(id int) (err error) {
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

func (sr *CategoryRepository) DeleteForce(id int) (err error) {
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
