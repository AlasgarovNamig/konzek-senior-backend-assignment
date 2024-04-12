package interceptors

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"product-categories-service/config"
	"product-categories-service/utils"
	"strings"
)

func UnaryAuthHandler(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	utils.Log("INFO", fmt.Sprintf("Authenticating request for %s", info.FullMethod))

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok || len(md["authorization"]) == 0 {
		errMsg := "Unauthorized Client: no authorization metadata"
		utils.Log("ERROR", errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	authHeader := md["authorization"][0]
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == "" {
		errMsg := "Unauthorized Client: bearer token not found"
		utils.Log("ERROR", errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	decoded, err := base64.StdEncoding.DecodeString(config.Configuration.KeyCloak.RealmKonzekRS256PublicKey)
	if err != nil {
		utils.Log("ERROR", fmt.Sprintf("Failed to decode base64 string: %v", err))
		return nil, fmt.Errorf("Unauthorized Client due to public key decode error")
	}

	pubKey, err := parseRSAPublicKeyFromDER(decoded)
	if err != nil {
		utils.Log("ERROR", fmt.Sprintf("Failed to parse RSA public key: %v", err))
		return nil, fmt.Errorf("Unauthorized Client due to public key parse error")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return pubKey, nil
	})

	if err != nil || !token.Valid {
		utils.Log("ERROR", fmt.Sprintf("Token parse error or invalid token: %v", err))
		return nil, fmt.Errorf("Unauthorized Client due to token validation error")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("Unauthorized Client: claims could not be processed")
	}

	if info.FullMethod == "/category.ProductCategoryService/SearchCategories" {
		authorized, errMsg := processAuthorization(claims, "client_id", "category_detail")
		if !authorized {
			utils.Log("ERROR", errMsg)
			return nil, fmt.Errorf(errMsg)
		}
	} else if info.FullMethod == "/category.ProductCategoryService/CreateCategory" {
		authorized, errMsg := processAuthorization(claims, "client_id", "category_create")
		if !authorized {
			utils.Log("ERROR", errMsg)
			return nil, fmt.Errorf(errMsg)
		}
	} else {
		errMsg := "Attempted to access an unknown method"
		utils.Log("ERROR", errMsg)
		return nil, fmt.Errorf(errMsg)
	}
	return handler(ctx, req)
}

func processAuthorization(claims jwt.MapClaims, clientIdKey, roleName string) (bool, string) {
	utils.Log("INFO", fmt.Sprintf("Processing authorization for role: %s", roleName))
	if claims[clientIdKey] == nil { // This is a user
		utils.Log("INFO", "Detected user token")
		authorized := checkRoleForUser(claims, "roles", roleName)
		if !authorized {
			errMsg := fmt.Sprintf("Unauthorized User: missing %s role", roleName)
			utils.Log("ERROR", errMsg)
			return false, errMsg
		}
	} else {
		utils.Log("INFO", "Detected client token")
		authorized := checkRoleForClient(claims, "realm-management", roleName)
		if !authorized {
			errMsg := fmt.Sprintf("Unauthorized Client: missing %s role", roleName)
			utils.Log("ERROR", errMsg)
			return false, errMsg
		}
	}
	utils.Log("INFO", fmt.Sprintf("Authorization successful for role: %s", roleName))
	return true, ""
}

func parseRSAPublicKeyFromDER(derBytes []byte) (*rsa.PublicKey, error) {
	utils.Log("INFO", "Parsing RSA public key from DER")
	pubKey, err := x509.ParsePKIXPublicKey(derBytes)
	if err != nil {
		utils.Log("ERROR", fmt.Sprintf("Failed to parse RSA public key: %v", err))
		return nil, fmt.Errorf("Key parse error: %v", err)
	}

	switch pub := pubKey.(type) {
	case *rsa.PublicKey:
		utils.Log("INFO", "RSA public key successfully parsed")
		return pub, nil
	default:
		errMsg := "Key type is not RSA"
		utils.Log("ERROR", errMsg)
		return nil, fmt.Errorf(errMsg)
	}
}

func checkRoleForClient(claims jwt.MapClaims, resourceName, roleName string) bool {
	utils.Log("INFO", fmt.Sprintf("Checking client role: %s", roleName))
	if resourceAccess, ok := claims["resource_access"].(map[string]interface{}); ok {
		if resource, ok := resourceAccess[resourceName].(map[string]interface{}); ok {
			if roles, ok := resource["roles"].([]interface{}); ok {
				for _, r := range roles {
					if roleStr, ok := r.(string); ok && roleStr == roleName {
						utils.Log("INFO", fmt.Sprintf("Client has required role: %s", roleName))
						return true
					}
				}
			}
		}
	}
	return false
}

func checkRoleForUser(claims jwt.MapClaims, resourceName, roleName string) bool {
	utils.Log("INFO", fmt.Sprintf("Checking user role: %s", roleName))
	if realmAccess, ok := claims["realm_access"].(map[string]interface{}); ok {
		if roles, ok := realmAccess[resourceName].([]interface{}); ok {
			for _, r := range roles {
				if roleStr, ok := r.(string); ok && roleStr == roleName {
					utils.Log("INFO", fmt.Sprintf("User has required role: %s", roleName))
					return true
				}
			}
		}
	}
	return false
}
