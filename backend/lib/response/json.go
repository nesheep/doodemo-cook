package response

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func JSON(ctx context.Context, w http.ResponseWriter, body any, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		log.Printf("encode response error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)

		rsp := ErrorResponse{Message: http.StatusText(http.StatusInternalServerError)}
		if err := json.NewEncoder(w).Encode(rsp); err != nil {
			log.Printf("write error response error: %v", err)
		}
		return
	}

	w.WriteHeader(code)
	if _, err := fmt.Fprintf(w, "%s", bodyBytes); err != nil {
		log.Printf("write response error: %v", err)
	}
}
