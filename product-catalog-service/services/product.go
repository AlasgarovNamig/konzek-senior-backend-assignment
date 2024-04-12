package services

import (
	"golang.org/x/net/context"
	"product-catalog-service/controllers"
	"product-catalog-service/rpc/product"
	//"product-categories-service/controllers"
)

type ProductCatalogService struct {
	controller controllers.IProductCatalogController
}

func NewProductCatalogService(controller controllers.IProductCatalogController) *ProductCatalogService {
	if controller == nil {
		panic("controller is nil")
	}
	return &ProductCatalogService{
		controller: controller,
	}
}
func (s *ProductCatalogService) SearchProduct(ctx context.Context, req *product.SearchRequest) (*product.SearchProductResponse, error) {
	return s.controller.SearchProduct(ctx, req), nil
}
func (s *ProductCatalogService) CreateProduct(ctx context.Context, req *product.CreateProductRequest) (*product.CreateProductResponse, error) {
	return s.controller.CreateProduct(ctx, req), nil
}
func (s *ProductCatalogService) GetProductById(ctx context.Context, req *product.GetProductByIdRequest) (*product.GetProductByIdResponse, error) {
	return s.controller.GetProductById(ctx, req), nil
}
