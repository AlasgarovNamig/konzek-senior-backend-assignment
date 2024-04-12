package utils

import (
	"product-categories-service/rpc/category"
)

func BuildSuccessResult(message string, status string) *category.Result {
	return &category.Result{
		IsSuccess:  true,
		StatusCode: status,
		Message:    message,
		Error:      "",
	}

}

func BuildErrorResult(status string, err string) *category.Result {
	return &category.Result{
		IsSuccess:  false,
		Message:    "",
		StatusCode: status,
		Error:      err,
	}
}
