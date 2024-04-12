package category

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"product-catalog-service/config"
	"product-catalog-service/integrations/keycloak"

	"time"
)

type CategoryClient struct {
	Category        ProductCategoryServiceClient
	IdentityManager keycloak.IIdentityManager
}

func NewCategoryClient(cfg *config.Config) (*CategoryClient, error) {
	var (
		conn *grpc.ClientConn
		err  error
	)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err = grpc.DialContext(ctx, fmt.Sprintf("%s:%s", cfg.CategoryServer.Service, cfg.CategoryServer.RPCPort), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return nil, err
	}

	client := &CategoryClient{
		Category:        NewProductCategoryServiceClient(conn),
		IdentityManager: keycloak.NewIdentityManager(),
	}

	return client, err
}

func (c *CategoryClient) CreateCategory(ctx context.Context, req *CreateCategoryRequest) (*CreateCategoryResponse, error) {
	md, err := c.addMetadata(ctx)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(metadata.NewOutgoingContext(ctx, md), 120*time.Second)
	defer cancel()
	return c.Category.CreateCategory(ctx, req)
}

func (c *CategoryClient) SearchCategories(ctx context.Context, req *SearchRequest) (*SearchCategoryResponse, error) {
	md, err := c.addMetadata(ctx)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(metadata.NewOutgoingContext(ctx, md), 120*time.Second)
	defer cancel()
	return c.Category.SearchCategories(ctx, req)
}

func (c *CategoryClient) addMetadata(ctx context.Context) (metadata.MD, error) {
	token, err := c.IdentityManager.LoginWithClient(ctx)
	if err != nil {
		return nil, err
	}
	return metadata.New(map[string]string{
		"authorization": "Bearer " + token.AccessToken,
	}), nil

}
