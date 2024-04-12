package utils

var validRoles = map[string]bool{
	"product_read_all": true,
	"product_detail":   true,
	"product_create":   true,
	"category_detail":  true,
	"category_create":  true,
}

func ValidateUserRoles(roles []string) bool {
	for _, role := range roles {
		if _, exists := validRoles[role]; !exists {
			return false
		}
	}
	return true
}
