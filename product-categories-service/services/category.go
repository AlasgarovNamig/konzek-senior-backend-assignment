package services

import (
	"fmt"
	"golang.org/x/net/context"
	"product-categories-service/controllers"
	"product-categories-service/rpc/category"
	"product-categories-service/utils"
)

type ProductCategoryService struct {
	controller controllers.IProductCategoryController
}

func NewProductCategoryService(controller controllers.IProductCategoryController) *ProductCategoryService {
	if controller == nil {
		utils.Log("ERROR", "Attempted to initialize ProductCategoryService with a nil controller")
		panic("ProductCategoryService is nil")
	}
	utils.Log("INFO", "ProductCategoryService initialized successfully")
	return &ProductCategoryService{
		controller: controller,
	}
}

func (s *ProductCategoryService) CreateCategory(ctx context.Context, req *category.CreateCategoryRequest) (*category.CreateCategoryResponse, error) {
	utils.Log("INFO", fmt.Sprintf("Creating new category: %s", req.GetName()))
	return s.controller.CreateCategory(ctx, req), nil
}

func (s *ProductCategoryService) SearchCategories(ctx context.Context, req *category.SearchRequest) (*category.SearchCategoryResponse, error) {
	utils.Log("INFO", fmt.Sprintf("Searching for categories"))
	return s.controller.SearchCategories(ctx, req), nil
}
