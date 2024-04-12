package utils

import (
	"product-catalog-service/rpc/product"
)

func CategoryDtoToProductCatalogCategoryResponse(id int64, name string) *product.Category {
	Log("INFO", "Converting domain category list to DTO list...")
	return &product.Category{
		Id:   id,
		Name: name,
	}
}
