package repositories

import (
	"fmt"
	"gorm.io/gorm"
	"product-categories-service/domains"
	"product-categories-service/rpc/category"
	"product-categories-service/utils"
)

type ICategoryRepository interface {
	GetAll() (*[]domains.Category, error)
	Search(req *category.SearchRequest) (*[]domains.Category, error)
	Create(entity *domains.Category) error
}

type categoryRepository struct {
	conn *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) ICategoryRepository {
	if db == nil {
		utils.Log("ERROR", "Database connection cannot be nil in NewCategoryRepository")
		panic("NewCategoryRepository database connection is nil ")
	}
	utils.Log("INFO", "CategoryRepository initialized successfully")
	return &categoryRepository{
		conn: db,
	}
}

func (r *categoryRepository) GetAll() (*[]domains.Category, error) {
	utils.Log("INFO", "Retrieving all categories from the database")
	var categories []domains.Category
	result := r.conn.Find(&categories)
	if result.Error != nil {
		utils.Log("ERROR", fmt.Sprintf("Failed to retrieve all categories: %v", result.Error))
		return nil, result.Error
	}
	utils.Log("INFO", fmt.Sprintf("Successfully retrieved %d categories", len(categories)))
	return &categories, nil
}

func (r *categoryRepository) Search(req *category.SearchRequest) (*[]domains.Category, error) {
	utils.Log("INFO", fmt.Sprintf("Searching for categories with criteria: %+v", req))
	categories := &[]domains.Category{}
	db := r.conn.Model(&domains.Category{})
	var baseQuery *gorm.DB

	for i, field := range req.SearchFields {
		if i == 0 {
			condition, value, err := getCriteriaForProduct(field)
			if err != nil {
				return nil, err
			}
			baseQuery = db.Where(condition, value)
		} else {
			condition, value, err := getCriteriaForProduct(field)
			if err != nil {
				return nil, err
			}
			if field.MatchType == category.MatchType_AND {
				baseQuery = baseQuery.Where(condition, value)
			} else if field.MatchType == category.MatchType_OR {
				baseQuery = baseQuery.Or(condition, value)
			} else {
				utils.Log("ERROR", fmt.Sprintf("Attempted to search using an unauthorized match type: %s", field.MatchType))
				return nil, fmt.Errorf("search on match type '%s' is not allowed", field.MatchType)
			}
		}
	}

	if baseQuery == nil {
		baseQuery = db
	}

	offset := (req.Page - 1) * req.Limit
	baseQuery = baseQuery.Offset(int(offset)).Limit(int(req.Limit))

	if err := baseQuery.Find(categories).Error; err != nil {
		utils.Log("ERROR", fmt.Sprintf("Failed to search categories: %v", err))
		return nil, err
	}

	utils.Log("INFO", fmt.Sprintf("Found %d categories matching the search criteria", len(*categories)))
	return categories, nil
}

func (r *categoryRepository) Create(entity *domains.Category) error {
	utils.Log("INFO", fmt.Sprintf("Creating a new category: %s", entity.Name))
	result := r.conn.Create(entity)
	if result.Error != nil {
		utils.Log("ERROR", fmt.Sprintf("Failed to create a new category '%s': %v", entity.Name, result.Error))
		return result.Error
	}
	utils.Log("INFO", fmt.Sprintf("Successfully created a new category: %s", entity.Name))
	return nil
}

func getCriteriaForProduct(field *category.SearchField) (string, interface{}, error) {
	switch {
	case field.FieldName == "id" && field.SearchIntData != 0:
		return fmt.Sprintf("%s = ?", field.FieldName), field.SearchIntData, nil
	case field.FieldName == "name" && field.SearchStringData != "":
		return fmt.Sprintf("%s = ?", field.FieldName), field.SearchStringData, nil
	case field.FieldName == "category_id" && field.SearchIntData != 0:
		return fmt.Sprintf("%s = ?", field.FieldName), field.SearchIntData, nil
	case field.FieldName == "price" && field.SearchDecimalData != 0:
		return fmt.Sprintf("%s = ?", field.FieldName), field.SearchDecimalData, nil
	default:
		utils.Log("ERROR", fmt.Sprintf("Attempted to search using an unauthorized field: %s", field.FieldName))
		return "", nil, fmt.Errorf("search on field '%s' is not allowed", field.FieldName)
	}
}
