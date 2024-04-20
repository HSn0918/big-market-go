package xcode

import (
	"context"
	"net/http"

	"github.com/hsn0918/BigMarket/pkg/xcode/types"
)

func ErrHandler(err error) (int, any) {
	code := CodeFromError(err)

	return http.StatusOK, types.Status{
		Code:    int32(code.Code()),
		Message: code.Message(),
	}
}
func OkHandler(ctx context.Context, data any) any {
	resp := types.Response{
		Code:    200, // 或者其他业务成功的代码
		Message: "操作成功",
		Data:    data,
	}
	return resp
}
