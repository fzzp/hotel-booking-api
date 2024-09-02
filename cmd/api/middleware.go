package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/fzzp/gotk"
	"github.com/fzzp/gotk/token"
	"github.com/fzzp/hotel-booking-api/pkg/errs"
	"github.com/go-chi/cors"
)

const (
	tokenType      = "bearer"
	tokenHeaderKey = "Authorization"
)

// 不需要验证的路由
var unAuthPatterns = []string{}

// 中间件函数签名
type mwHandler func(next http.Handler) http.Handler

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
		fmt.Println("->>> 执行 RequiredAuth ")
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

// RequiredRole 指定角色才通过, 利用闭包传递adminRole
func (app *application) RequiredRole(adminRole int) mwHandler {
	return func(next http.Handler) http.Handler {
		fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("->>> 执行 RequiredRole ")
			var payload *token.Payload
			payload = GetByContext(r, tokenPayloadKey, payload)
			if payload == nil {
				slog.ErrorContext(r.Context(), "从上下文获取payload为nil")
				app.FAIL(w, r, errs.ErrServerError)
				return
			}
			userID, _ := strconv.Atoi(payload.UserText)
			if userID <= 0 {
				slog.ErrorContext(r.Context(), "从上下文 token payload 中 userID <= 0", slog.Any("payload", payload))
				app.FAIL(w, r, errs.ErrServerError)
				return
			}

			user, err := app.service.User.GetUserByID(uint(userID))
			if err != nil {
				app.FAIL(w, r, err)
				return
			}
			if user.Role != int8(adminRole) {
				app.FAIL(w, r, errs.ErrForbidden.AsMessage("没有权限"))
				return
			}
			// 存储user到上下文
			r = r.WithContext(context.WithValue(r.Context(), userInfoKey, user))
			next.ServeHTTP(w, r)
		})

		return fn
	}
}
