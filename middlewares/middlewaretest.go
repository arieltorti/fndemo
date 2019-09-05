package main

import (
	"context"
	"strings"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/fnproject/fn/api/server"

	// The trace package is imported in several places by different dependencies and if we don't import explicity here it is
	// initialized every time it is imported and that creates a panic at run time as we register multiple time the handler for
	// /debug/requests. For example see: https://github.com/GoogleCloudPlatform/google-cloud-go/issues/663 and https://github.com/bradleyfalzon/gopherci/issues/101
	//_ "golang.org/x/net/trace"

	// EXTENSIONS: Add extension imports here or use `fn build-server`. Learn more: https://github.com/fnproject/fn/blob/master/docs/operating/extending.md
	_ "github.com/fnproject/fn/api/server/defaultexts"
)

func main() {
	ctx := context.Background()
	//registerViews()

	funcServer := server.NewFromEnv(ctx)

	// They run in the order you register them
	funcServer.AddMiddleware(&requiresFooMiddleware{})
	funcServer.AddMiddleware(&requiresBarMiddleware{})
	funcServer.AddMiddleware(&passMiddleware{})
	funcServer.Start(ctx)
}

type requiresFooMiddleware struct {}
func (h *requiresFooMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Foo middleware started")
		fmt.Println(r.URL.Path)

		tokenHeader := r.Header.Get("foo")
		if len(tokenHeader) < 1 && !isPublicPath(r.URL) {
			w.WriteHeader(http.StatusUnauthorized)
			m := map[string]string{"error": "foo header required"}
			json.NewEncoder(w).Encode(m)
			return
		}

		fmt.Println("Foo middleware finished")
		next.ServeHTTP(w, r)
	})
}

type requiresBarMiddleware struct {}
func (h *requiresBarMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Bar middleware started")

		tokenHeader := r.Header.Get("bar")
		if len(tokenHeader) < 1 && !isPublicPath(r.URL) {
			w.WriteHeader(http.StatusUnauthorized)
			m := map[string]string{"error": "bar header required"}
			json.NewEncoder(w).Encode(m)
			return
		}

		fmt.Println("Bar middleware finished")
		next.ServeHTTP(w, r)
	})
}

type passMiddleware struct {}
func (h *passMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Pass middleware started")

		r.Header.Set("Pass", "Middleware")
		r = r.WithContext(context.WithValue(r.Context(), "passMiddleware", "Passed."))

		fmt.Println("Pass middleware finished")
		next.ServeHTTP(w, r)
	})
}

func isPublicPath(url *url.URL) bool {
	return strings.HasPrefix(url.Path,"/v2/")
}
