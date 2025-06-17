package repos

import (
	"math"
	"strings"

	"github.com/12ilya12/go-proj-mng/common"
	"github.com/12ilya12/go-proj-mng/models"
	"github.com/12ilya12/go-proj-mng/pagination"
	"gorm.io/gorm"
)

type DependencyRepository struct {
	DB *gorm.DB
}

func NewDependencyRepository(DB *gorm.DB) DependencyRepository {
	return DependencyRepository{DB}
}

func (sr *DependencyRepository) Get(parentTaskId uint, pagingOptions pagination.PagingOptions) (dependenciesWithPaging pagination.Paging[models.Dependency], err error) {
	//Сортировка. По умолчанию по возрастанию идентификатора.
	var orderRule string
	if pagingOptions.OrderBy == "" {
		orderRule = "id"
	} else {
		var columnCount int64
		sr.DB.Select("column_name").Table("information_schema.columns").
			Where("table_name = ? AND column_name = ?", "dependencies", pagingOptions.OrderBy).Count(&columnCount)
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

	err = tx.Find(&dependenciesWithPaging.Items, "parentTaskId = ?", parentTaskId).Error

	//Собираем выходные данные пагинации
	sr.DB.Model(&models.Dependency{}).Where("parentTaskId = ?", parentTaskId).Count(&dependenciesWithPaging.Pagination.TotalItems)
	if pagingOptions.PageSize == 0 { //Если размер страницы не задан, показываем всё на одной странице
		dependenciesWithPaging.Pagination.TotalPages = 1
	} else { //Подсчитываем количество страниц
		dependenciesWithPaging.Pagination.TotalPages =
			int64(math.Ceil(float64(dependenciesWithPaging.Pagination.TotalItems) /
				float64(pagingOptions.PageSize)))
	}
	dependenciesWithPaging.Pagination.Options = pagingOptions

	return
}

func (sr *DependencyRepository) Create(
	dependency *models.Dependency,
	userInfo common.UserInfo,
	parentTaskUserId uint,
	childTaskUserId uint) (err error) {

	err = sr.DB.Create(&dependency).Error
	return
}

func (sr *DependencyRepository) Delete(parentTaskId uint, dependencyId uint) (err error) {
	tx := sr.DB.Delete(&models.Dependency{}, "id = ? AND parent_task_id = ?", dependencyId, parentTaskId)
	if tx.RowsAffected == 0 {
		err = gorm.ErrRecordNotFound
	}
	return
}
