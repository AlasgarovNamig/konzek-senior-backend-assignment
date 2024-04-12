package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"product-catalog-service/domains"
	"product-catalog-service/integrations/cache"
	"product-catalog-service/integrations/category"
	"product-catalog-service/repositories"
	"product-catalog-service/rpc/product"
	"product-catalog-service/utils"
	"strconv"
	"time"
)

type ProductCatalogController struct {
	ProductRepository repositories.IProductRepository
	Cache             cache.Cache
	CategoryClient    CategoryClient
}
type IProductCatalogController interface {
	SearchProduct(ctx context.Context, req *product.SearchRequest) *product.SearchProductResponse
	GetProductById(ctx context.Context, req *product.GetProductByIdRequest) *product.GetProductByIdResponse
	CreateProduct(ctx context.Context, req *product.CreateProductRequest) *product.CreateProductResponse
}
type CategoryClient interface {
	CreateCategory(ctx context.Context, req *category.CreateCategoryRequest) (*category.CreateCategoryResponse, error)
	SearchCategories(ctx context.Context, req *category.SearchRequest) (*category.SearchCategoryResponse, error)
}

func NewProductController(
	productRepository repositories.IProductRepository,
	cache cache.Cache,
	categoryClient CategoryClient,
) *ProductCatalogController {
	if productRepository == nil {
		utils.Log("ERROR", "Attempted to initialize ProductCategoryController with a nil productRepository")
		panic("productRepository is nil")
	}
	if cache == nil {
		utils.Log("ERROR", "Attempted to initialize ProductCategoryController with a nil cache")
		panic("cache is nil")
	}
	if categoryClient == nil {
		utils.Log("ERROR", "Attempted to initialize ProductCategoryController with a nil categoryClient")
		panic("categoryClient is nil")
	}
	return &ProductCatalogController{
		ProductRepository: productRepository,
		Cache:             cache,
		CategoryClient:    categoryClient,
	}
}

func (c *ProductCatalogController) SearchProduct(ctx context.Context, req *product.SearchRequest) *product.SearchProductResponse {
	utils.Log("INFO", fmt.Sprintf("Searching for products with criteria: %+v", req))
	response := &product.SearchProductResponse{
		ProductList: []*product.Product{},
	}
	products, err := c.ProductRepository.Search(req)
	if err != nil {
		utils.Log("INFO", fmt.Sprintf("Failed to search products with criteria %+v: %v", req, err))
		response.Result = utils.BuildErrorResult("400", err.Error())
		return response
	}

	for _, productDomain := range *products {
		utils.Log("INFO", fmt.Sprintf("Found product: %+v", productDomain))

		// Category Get by Redis Cache
		productReps, err := c.getCategoryFromRedis(ctx, &productDomain)
		if err == nil {
			response.ProductList = append(response.ProductList, productReps)
			continue
		}

		// Get by gRPC Client
		productReps, err = c.getCategoryFromProductCategoryServiceWithGRPSClient(ctx, &productDomain)
		if err != nil {
			response.ProductList = []*product.Product{}
			response.Result = utils.BuildErrorResult("400", err.Error())
			return response
		}
		response.ProductList = append(response.ProductList, productReps)
	}

	utils.Log("INFO", fmt.Sprintf("Successfully found %d products matching criteria: %+v", len(*products), req))
	response.Result = utils.BuildSuccessResult("200", "product get successfully")
	return response
}
func (c *ProductCatalogController) GetProductById(ctx context.Context, req *product.GetProductByIdRequest) *product.GetProductByIdResponse {
	response := &product.GetProductByIdResponse{}
	parsedProductId, err := strconv.ParseUint(req.ProductId, 10, 64)
	if err != nil {
		response.Result = utils.BuildErrorResult("400", err.Error())
		return response
	}
	productDomain, err := c.ProductRepository.GetByID(uint(parsedProductId))
	if err != nil {
		response.Result = utils.BuildErrorResult("400", err.Error())
		return response
	}
	// Category Get by Redis Cache
	productReps, err := c.getCategoryFromRedis(ctx, productDomain)
	if err == nil {
		response.Result = utils.BuildSuccessResult("200", "product get successfully")
		response.Product = productReps
		return response
	}

	// Category Get by gRPC Client
	productReps, err = c.getCategoryFromProductCategoryServiceWithGRPSClient(ctx, productDomain)
	if err != nil {
		response.Result = utils.BuildErrorResult("400", err.Error())
		return response
	}
	response.Result = utils.BuildSuccessResult("200", "product get successfully")
	response.Product = productReps
	return response
}
func (c *ProductCatalogController) CreateProduct(ctx context.Context, req *product.CreateProductRequest) *product.CreateProductResponse {
	response := &product.CreateProductResponse{}
	// Check Category Exist by Redis Cache
	get, err := c.Cache.Get(ctx, fmt.Sprintf("%d", req.CategoryId))
	if err == nil && get != "" {
		err = c.ProductRepository.Create(&domains.Product{
			Name:       req.Name,
			Price:      req.Price,
			CategoryID: uint(req.CategoryId),
			CreatedAt:  time.Now(),
		})
		if err != nil {
			response.Result = utils.BuildErrorResult("400", err.Error())
			return response
		}
		response.Result = utils.BuildSuccessResult("200", "product successfully created")
		return response
	}

	// Check Category Exist by gRPC Client
	categorySearchReq := &category.SearchRequest{
		SearchFields: []*category.SearchField{{
			FieldName:      "id",
			SearchIntData:  req.CategoryId,
			SearchOperator: category.SearchOperator_EQUAL,
		}},
		Limit: 1,
		Page:  1,
	}
	categories, err := c.CategoryClient.SearchCategories(ctx, categorySearchReq)
	if err != nil || !categories.Result.IsSuccess || categories.CategoryDtoList == nil {
		response.Result = utils.BuildErrorResult("400", "invalid category by grpc")
		return response
	}

	err = c.ProductRepository.Create(&domains.Product{
		Name:       req.Name,
		Price:      req.Price,
		CategoryID: uint(req.CategoryId),
		CreatedAt:  time.Now(),
	})
	if err != nil {
		response.Result = utils.BuildErrorResult("400", err.Error())
		return response
	}
	response.Result = utils.BuildSuccessResult("200", "product successfully created")
	return response
}

func (c *ProductCatalogController) getCategoryFromRedis(ctx context.Context, productDomain *domains.Product) (*product.Product, error) {
	get, cacheErr := c.Cache.Get(ctx, fmt.Sprintf("%d", productDomain.CategoryID))
	if cacheErr == nil && get != "" {
		var categoryResp product.Category
		err := json.Unmarshal([]byte(get), &categoryResp)
		if err == nil {
			return &product.Product{
				Id:       int64(productDomain.ID),
				Name:     productDomain.Name,
				Price:    productDomain.Price,
				Category: &categoryResp,
			}, nil

		}
		return nil, fmt.Errorf("")
	}
	return nil, fmt.Errorf("")
}
func (c *ProductCatalogController) getCategoryFromProductCategoryServiceWithGRPSClient(ctx context.Context, productDomain *domains.Product) (*product.Product, error) {
	categorySearchReq := &category.SearchRequest{
		SearchFields: []*category.SearchField{{
			FieldName:      "id",
			SearchIntData:  int64(productDomain.CategoryID),
			SearchOperator: category.SearchOperator_EQUAL,
		}},
		Limit: 1,
		Page:  1,
	}
	categories, err := c.CategoryClient.SearchCategories(ctx, categorySearchReq)
	if err != nil || !categories.Result.IsSuccess || categories.CategoryDtoList == nil {
		return nil, fmt.Errorf("")
	}

	return &product.Product{
		Id:       int64(productDomain.ID),
		Name:     productDomain.Name,
		Price:    productDomain.Price,
		Category: utils.CategoryDtoToProductCatalogCategoryResponse(categories.CategoryDtoList[0].Id, categories.CategoryDtoList[0].Name),
	}, nil
}
