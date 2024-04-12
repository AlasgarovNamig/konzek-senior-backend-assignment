package utils

import "strings"

type Response struct {
	IsSuccess bool        `json:"isSuccess"`
	Status    string      `json:"status"`
	Message   string      `json:"message"`
	Error     interface{} `json:"error"`
	Data      interface{} `json:"data"`
}

func BuildSuccessResponse(message string, status string, data interface{}) Response {
	return Response{
		IsSuccess: true,
		Status:    status,
		Message:   message,
		Error:     nil,
		Data:      data,
	}

}

func BuildErrorResponse(message string, status string, err string) Response {
	return Response{
		IsSuccess: false,
		Message:   message,
		Status:    status,
		Error:     strings.Split(err, "\n"),
		Data:      new(struct{}),
	}
}
