package middlewares

import (
	"auth-service/config"
	"auth-service/utils"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
)

type IJWTService interface {
	ValidateToken() gin.HandlerFunc
	AuthorizeByRole(expectedRole string) gin.HandlerFunc
}

type jwtService struct {
	realmRS256PublicKey string
}

func NewJWTService() IJWTService {
	return &jwtService{
		realmRS256PublicKey: config.Configuration.KeyCloak.RealmMasterRS256PublicKey,
	}
}

func (j *jwtService) ValidateToken() gin.HandlerFunc {
	return func(context *gin.Context) {
		decoded, err := base64.StdEncoding.DecodeString(j.realmRS256PublicKey)
		if err != nil {
			utils.Log("ERROR", fmt.Sprintf("Error decoding base64 string for the public key: %v\", err"))
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		pubKey, err := j.parseRSAPublicKeyFromDER(decoded)
		if err != nil {
			utils.Log("ERROR", fmt.Sprintf("Failed to parse RSA public key from decoded base64 string: %v", err))
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		tokenString, err := j.extractToken(context)
		if err != nil {
			utils.Log("ERROR", "Authorization token not found in the request headers.")
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				utils.Log("ERROR", fmt.Sprintf("Unexpected signing method in JWT: expected=%v, received=%v", jwt.SigningMethodRS256.Alg(), token.Header["alg"]))
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return pubKey, nil
		})

		if err != nil {
			utils.Log("ERROR", fmt.Sprintf("Failed to parse JWT: %v", err))
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		if token.Valid {
			utils.Log("INFO", fmt.Sprintf("Successfully validated token: %v", tokenString))
			context.Next()
		} else {
			utils.Log("ERROR", fmt.Sprintf("Failed to validate token: %v", err))
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
	}
}

func (j *jwtService) AuthorizeByRole(expectedRole string) gin.HandlerFunc {
	return func(context *gin.Context) {
		decoded, err := base64.StdEncoding.DecodeString(j.realmRS256PublicKey)
		if err != nil {
			utils.Log("ERROR", fmt.Sprintf("Error decoding base64 string for the public key: %v\", err"))
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		// RSA public key'e dönüştür
		pubKey, err := j.parseRSAPublicKeyFromDER(decoded)
		if err != nil {
			utils.Log("ERROR", fmt.Sprintf("Failed to parse RSA public key from decoded base64 string: %v", err))
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		tokenString, err := j.extractToken(context)
		if err != nil {
			utils.Log("ERROR", "Authorization token not found in the request headers.")
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				utils.Log("ERROR", fmt.Sprintf("Unexpected signing method in JWT: expected=%v, received=%v", jwt.SigningMethodRS256.Alg(), token.Header["alg"]))
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return pubKey, nil
		})

		if err != nil {
			utils.Log("ERROR", fmt.Sprintf("Failed to parse JWT: %v", err))
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			preferredUsername, ok := claims["preferred_username"].(string)
			if !ok {
				utils.Log("ERROR", fmt.Sprintf("Failed to parse JWT claims: %v", err))
				context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
				return
			}

			if preferredUsername == expectedRole {
				utils.Log("INFO", fmt.Sprintf("Successfully validated token: %v", tokenString))
				context.Next()
			} else {
				utils.Log("ERROR", fmt.Sprintf("Failed to validate role: %v", expectedRole))
				context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
				return
			}
		} else {
			utils.Log("ERROR", "invalid token")
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
	}
}
func (j *jwtService) extractToken(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {

		return "", fmt.Errorf("Token Not Found")
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	return strings.TrimSpace(tokenString), nil
}

func (j *jwtService) parseRSAPublicKeyFromDER(derBytes []byte) (*rsa.PublicKey, error) {
	pubKey, err := x509.ParsePKIXPublicKey(derBytes)
	if err != nil {
		return nil, err
	}

	switch pub := pubKey.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
		return nil, fmt.Errorf("key type is not RSA")
	}
}
