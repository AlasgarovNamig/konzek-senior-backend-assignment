package utils

import (
	"product-catalog-service/rpc/product"
)

func BuildSuccessResult(message string, status string) *product.Result {
	return &product.Result{
		IsSuccess:  true,
		StatusCode: status,
		Message:    message,
		Error:      "",
	}

}

func BuildErrorResult(status string, err string) *product.Result {
	return &product.Result{
		IsSuccess:  false,
		Message:    "",
		StatusCode: status,
		Error:      err,
	}
}
