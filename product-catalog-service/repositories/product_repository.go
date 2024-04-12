package repositories

import (
	"fmt"
	"gorm.io/gorm"
	"product-catalog-service/domains"
	"product-catalog-service/rpc/product"
	"product-catalog-service/utils"
)

type IProductRepository interface {
	GetByID(id uint) (*domains.Product, error)
	Create(entity *domains.Product) error
	Search(req *product.SearchRequest) (*[]domains.Product, error)
}

type productRepository struct {
	conn *gorm.DB
}

func NewProductRepository(db *gorm.DB) IProductRepository {
	return &productRepository{
		conn: db,
	}
}

func (r *productRepository) GetByID(id uint) (*domains.Product, error) {
	var domainProduct domains.Product
	err := r.conn.First(&domainProduct, id).Error
	if err != nil {
		return nil, err
	}
	return &domainProduct, nil
}

func (r *productRepository) Create(entity *domains.Product) error {
	result := r.conn.Create(entity)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *productRepository) Search(req *product.SearchRequest) (*[]domains.Product, error) {
	utils.Log("INFO", fmt.Sprintf("Searching for products with criteria: %+v", req))
	products := &[]domains.Product{}
	db := r.conn.Model(&domains.Product{})
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
			if field.MatchType == product.MatchType_AND {
				baseQuery = baseQuery.Where(condition, value)
			} else if field.MatchType == product.MatchType_OR {
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
	if err := baseQuery.Find(products).Error; err != nil {
		utils.Log("ERROR", fmt.Sprintf("Failed to search products: %v", err))
		return nil, err
	}

	utils.Log("INFO", fmt.Sprintf("Found %d products matching the search criteria", len(*products)))
	return products, nil

}

func getCriteriaForProduct(field *product.SearchField) (string, interface{}, error) {
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
