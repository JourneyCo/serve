package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
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

func JSONToCtx(ifc any, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		body, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("error reading body: %v", err)
			return
		}

		dto := reflect.New(reflect.TypeOf(ifc)).Interface()

		if err = json.Unmarshal(body, &dto); err != nil {
			fmt.Printf("error unmarshalling body: %v", err)
			return
		}

		// Add JSON data to the context
		ctx = context.WithValue(ctx, "dto", &dto)
		ctx = context.WithValue(ctx, "body", body)

		// Pass the updated context to the next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// DisableCORS disables CORS.
func DisableCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		h.ServeHTTP(w, r)
	})
}
