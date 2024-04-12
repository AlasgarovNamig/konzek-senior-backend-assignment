package interceptors

import (
	"context"
	"fmt"
	"github.com/mitchellh/copystructure"
	"google.golang.org/grpc"
	"product-catalog-service/utils"
	"runtime"
	"time"
)

var _ grpc.UnaryServerInterceptor = UnaryLogHandler

func UnaryLogHandler(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	reqCopy, copyErr := copystructure.Copy(req)
	if copyErr != nil {
		utils.Log("ERROR", fmt.Sprintf("Failed to copy request for logging: %v", copyErr))
		reqCopy = "request copy error; original request not modified"
	}
	start := time.Now()
	resp, err = handler(ctx, req)
	duration := time.Since(start)
	if err != nil {
		utils.Log("ERROR", fmt.Sprintf("RPC call to %s failed after %s with error: %v | Request details: %v", info.FullMethod, duration, err, reqCopy))
	} else {
		utils.Log("INFO", fmt.Sprintf("RPC call to %s completed in %s | Request details: %v", info.FullMethod, duration, reqCopy))
	}

	return resp, err
}
func LogPanicStackMultiLine(ctx context.Context, r interface{}) {
	_, file, line, ok := runtime.Caller(0)
	if ok {
		utils.Log("ERROR", fmt.Sprintf("Recovered from panic: %v in %s(%d)", r, file, line))
	}

	callers := []string{}
	for i := 0; true; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		fn := runtime.FuncForPC(pc)
		callers = append(callers, fmt.Sprintf("%d: %s(%d): %s", i, file, line, fn.Name()))
	}
	utils.Log("ERROR", "StackTrace:")
	for i := 0; len(callers) > i; i++ {
		utils.Log("ERROR", fmt.Sprintf("  %s", callers[i]))
	}
}
