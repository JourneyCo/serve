package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func HandleCacheControl(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		headers := rw.Header()
		headers.Set("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate")
		headers.Set("Pragma", "no-cache")
		headers.Set("Expires", "0")
		next.ServeHTTP(rw, req)
	})
}

func JSONToCtx(int any, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		body, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("error reading body: %v", err)
			return
		}
		if err = json.Unmarshal(body, &int); err != nil {
			fmt.Printf("error unmarshalling body: %v", err)
			return
		}

		// Add JSON data to the context
		ctx := context.WithValue(r.Context(), "dto", &int)

		// Pass the updated context to the next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
