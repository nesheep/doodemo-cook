package response

import (
	"context"
	"net/http"
)

type ErrorResponse struct {
	Message string `json:"message"`
	Detail  string `json:"detail"`
}

func FromStatusCode(ctx context.Context, w http.ResponseWriter, code int, detail string) {
	rsp := ErrorResponse{
		Message: http.StatusText(code),
		Detail:  detail,
	}

	JSON(ctx, w, rsp, code)
}
