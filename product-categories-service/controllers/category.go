package controllers

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"product-categories-service/domains"
	"product-categories-service/integrations/cache"
	"product-categories-service/integrations/database/postgres"
	"product-categories-service/repositories"
	"product-categories-service/rpc/category"
	"product-categories-service/utils"
	"time"
)

type IProductCategoryController interface {
	CreateCategory(context.Context, *category.CreateCategoryRequest) *category.CreateCategoryResponse
	SearchCategories(ctx context.Context, req *category.SearchRequest) *category.SearchCategoryResponse
}

type ProductCategoryController struct {
	CategoryRepository repositories.ICategoryRepository
	Cache              cache.Cache
}

func NewCategoryController(categoryRepository repositories.ICategoryRepository, cache cache.Cache) *ProductCategoryController {
	if categoryRepository == nil {
		utils.Log("ERROR", "Attempted to initialize ProductCategoryController with a nil categoryRepository")
		panic("ProductCategoryController is nil")
	}
	if cache == nil {
		utils.Log("ERROR", "Attempted to initialize ProductCategoryController with a nil cache")
		panic("cache is nil")
	}
	utils.Log("INFO", "ProductCategoryController initialized successfully")
	return &ProductCategoryController{
		CategoryRepository: categoryRepository,
		Cache:              cache,
	}
}

func (c *ProductCategoryController) CreateCategory(ctx context.Context, req *category.CreateCategoryRequest) *category.CreateCategoryResponse {
	utils.Log("INFO", fmt.Sprintf("Attempting to create a new category: %s", req.Name))
	transaction := postgres.Connection.Begin()

	newCategory := &domains.Category{
		Name:      req.Name,
		CreatedAt: time.Now(),
	}

	if err := transaction.Create(newCategory).Error; err != nil {
		utils.Log("ERROR", fmt.Sprintf("Failed to create a new category '%s': %v", req.Name, err))
		transaction.Rollback()
		return &category.CreateCategoryResponse{
			Result: utils.BuildErrorResult("400", err.Error()),
		}
	}

	categoryJSON, err := json.Marshal(newCategory)
	if err != nil {
		utils.Log("ERROR", fmt.Sprintf("Failed to marshal new category '%s' to JSON: %v", newCategory.Name, err))
		transaction.Rollback()
		return &category.CreateCategoryResponse{
			Result: utils.BuildErrorResult("400", err.Error()),
		}
	}

	err = c.Cache.Set(ctx, fmt.Sprintf("%d", newCategory.Id), categoryJSON, 0)
	if err != nil {
		utils.Log("ERROR", fmt.Sprintf("Failed to set new category '%s' in Redis: %v", newCategory.Name, err))
		transaction.Rollback()
		return &category.CreateCategoryResponse{
			Result: utils.BuildErrorResult("400", err.Error()),
		}
	}

	if err := transaction.Commit().Error; err != nil {
		utils.Log("ERROR", fmt.Sprintf("Failed to commit new category '%s' creation transaction: %v", newCategory.Name, err))
		return &category.CreateCategoryResponse{
			Result: utils.BuildErrorResult("400", err.Error()),
		}
	}

	utils.Log("INFO", fmt.Sprintf("Successfully created new category: %s", newCategory.Name))
	return &category.CreateCategoryResponse{
		Result: utils.BuildSuccessResult("201", "Category Successfully Created"),
	}
}

func (c *ProductCategoryController) SearchCategories(ctx context.Context, req *category.SearchRequest) *category.SearchCategoryResponse {
	utils.Log("INFO", fmt.Sprintf("Searching for categories with criteria: %+v", req))
	categories, err := c.CategoryRepository.Search(req)
	if err != nil {
		utils.Log("ERROR", fmt.Sprintf("Failed to search categories with criteria %+v: %v", req, err))
		return &category.SearchCategoryResponse{
			Result: utils.BuildErrorResult("400", err.Error()),
		}
	}
	var categoryDtoList []*category.CategoryDto

	for _, categoryDomain := range *categories {
		categoryDto := utils.CategoryDomainListToCategoryDtoList(categoryDomain.Id, categoryDomain.Name, categoryDomain.CreatedAt)
		categoryDtoList = append(categoryDtoList, categoryDto)
		utils.Log("INFO", fmt.Sprintf("Converted category '%s' to DTO format.", categoryDomain.Name))
	}
	utils.Log("INFO", fmt.Sprintf("Successfully found %d categories matching criteria: %+v", len(*categories), req))
	return &category.SearchCategoryResponse{
		CategoryDtoList: categoryDtoList,
		Result:          utils.BuildSuccessResult("201", "Category Successfully Get"),
	}
}
