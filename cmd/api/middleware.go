package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/fzzp/gotk"
	"github.com/fzzp/hotel-booking-api/pkg/errs"
	"github.com/go-chi/cors"
)

const (
	tokenType      = "bearer"
	tokenHeaderKey = "Authorization"
)

// 不需要验证的路由
var unAuthPatterns = []string{}

func (app *application) EnableCORS() func(http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
}

func (app *application) AccessLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			slog.InfoContext(
				r.Context(),
				r.Method+" "+r.URL.Path+" "+fmt.Sprint(time.Since(start)),
				slog.String("ip", r.RemoteAddr),
				slog.String("userAgent", r.UserAgent()),
				slog.String("status", w.Header().Get("Status")),
				slog.Any("query", r.URL.Query()),
			)
		}()
		next.ServeHTTP(w, r)
	})
}

func (app *application) RecoverPanic(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			f := recover()
			if f != nil {
				err, ok := f.(*gotk.ApiError)
				if ok {
					app.FAIL(w, r, err)
					return
				}
				w.Header().Set("Connection", "Close")
				app.FAIL(w, r, errs.ErrServerError.AsException(fmt.Errorf("%v", f)))
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func (app *application) RequiredAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, pattern := range unAuthPatterns {
			if pattern == r.URL.Path {
				next.ServeHTTP(w, r)
				return
			}
		}

		tokenStr := r.Header.Get(tokenHeaderKey)
		fields := strings.Fields(tokenStr)
		if len(fields) != 2 {
			app.FAIL(w, r, errs.ErrUnauthorized.AsMessage("令牌无效"))
			return
		}

		if strings.ToLower(fields[0]) != tokenType {
			app.FAIL(w, r, errs.ErrUnauthorized.AsMessage("不支持该令牌类型"))
			return
		}

		payload, err := app.jwt.ParseToken(fields[1])
		if err != nil {
			app.FAIL(w, r, errs.ErrUnauthorized.AsException(err))
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), tokenPayloadKey, payload))

		next.ServeHTTP(w, r)
	})
}
