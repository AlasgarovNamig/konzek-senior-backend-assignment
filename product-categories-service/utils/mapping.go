package utils

import (
	"product-categories-service/rpc/category"
	"time"
)

func CategoryDomainListToCategoryDtoList(id uint, name string, createdAt time.Time) *category.CategoryDto {
	Log("INFO", "Converting domain category to category DTO...")
	return &category.CategoryDto{
		Id:        int64(id),
		Name:      name,
		CreatedAt: createdAt.Format("2006-01-02 15:04:05"),
	}
}
