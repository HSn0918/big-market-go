package xcode

import (
	"github.com/hsn0918/BigMarket/pkg/xcode/types"
	"net/http"
)

func ErrHandler(err error) (int, any) {
	code := CodeFromError(err)

	return http.StatusOK, types.Status{
		Code:    int32(code.Code()),
		Message: code.Message(),
	}
}
