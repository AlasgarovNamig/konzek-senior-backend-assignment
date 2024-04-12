package api

import (
	"google.golang.org/grpc/status"
	"strings"
)

type Error struct {
	IsSuccess bool   `json:"isSuccess"`
	Status    int    `json:"status"`
	Error     string `json:"error"`
}

func GetGRPCErrorResponse(err error) *Error {
	if st, ok := status.FromError(err); ok {
		fullMessage := st.Message()
		parts := strings.SplitN(fullMessage, "desc = ", 2)
		if len(parts) == 2 {
			return &Error{
				IsSuccess: false,
				Status:    400,
				Error:     parts[1],
			}
		}
	}
	return &Error{
		IsSuccess: false,
		Status:    400,
		Error:     err.Error(),
	}

}
