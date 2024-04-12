package keycloak

import (
	"github.com/Nerzal/gocloak/v13"
	"product-catalog-service/config"

	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

type IIdentityManager interface {
	LoginWithClient(ctx context.Context) (*gocloak.JWT, error)
}
type identityManager struct {
	baseUrl      string
	realm        string
	ClientId     string
	ClientSecret string
}

func NewIdentityManager() IIdentityManager {
	return &identityManager{
		baseUrl:      config.Configuration.KeyCloak.BaseUrl,
		realm:        config.Configuration.KeyCloak.Realm,
		ClientId:     config.Configuration.KeyCloak.KeyCloakClient.ClientId,
		ClientSecret: config.Configuration.KeyCloak.KeyCloakClient.ClientSecret,
	}
}

func (im *identityManager) LoginWithClient(ctx context.Context) (*gocloak.JWT, error) {
	client := gocloak.NewClient(im.baseUrl)
	token, err := client.LoginClient(ctx, im.ClientId, im.ClientSecret, im.realm)
	if err != nil {
		return nil, errors.Wrap(err, "unable to login the rest client")
	}
	return token, nil
}
