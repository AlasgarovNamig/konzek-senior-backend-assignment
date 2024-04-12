package keycloak

import (
	"auth-service/config"
	"fmt"
	"strings"

	"github.com/Nerzal/gocloak/v13"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

type IIdentityManager interface {
	CreateUser(ctx *gin.Context, user gocloak.User, password string, role []string) (*gocloak.User, error)
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

func (im *identityManager) loginRestApiClient(ctx context.Context) (*gocloak.JWT, error) {
	client := gocloak.NewClient(im.baseUrl)

	token, err := client.LoginClient(ctx, im.ClientId, im.ClientSecret, im.realm)
	if err != nil {
		return nil, errors.Wrap(err, "unable to login the rest client")
	}
	return token, nil
}

func (im *identityManager) assignRolesToUser(ctx *gin.Context, client *gocloak.GoCloak, accessToken string, realm string, userID string, roles []string) error {
	for _, roleName := range roles {
		roleNameLowerCase := strings.ToLower(roleName)
		roleKeycloak, err := client.GetRealmRole(ctx, accessToken, realm, roleNameLowerCase)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("unable to get role by name: '%v'", roleNameLowerCase))
		}
		if err = client.AddRealmRoleToUser(ctx, accessToken, realm, userID, []gocloak.Role{*roleKeycloak}); err != nil {
			return errors.Wrap(err, "unable to add a realm role to user")
		}
	}
	return nil
}

func (im *identityManager) CreateUser(ctx *gin.Context, user gocloak.User, password string, roles []string) (*gocloak.User, error) {
	token, err := im.loginRestApiClient(ctx)
	if err != nil {
		return nil, err
	}

	client := gocloak.NewClient(im.baseUrl)
	userId, err := client.CreateUser(ctx, token.AccessToken, im.realm, user)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create the user")
	}

	if err = client.SetPassword(ctx, token.AccessToken, userId, im.realm, password, false); err != nil {
		return nil, errors.Wrap(err, "unable to set the password for the user")
	}

	if err = im.assignRolesToUser(ctx, client, token.AccessToken, im.realm, userId, roles); err != nil {
		return nil, err
	}

	userKeycloak, err := client.GetUserByID(ctx, token.AccessToken, im.realm, userId)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get recently created user")
	}

	return userKeycloak, nil
}
