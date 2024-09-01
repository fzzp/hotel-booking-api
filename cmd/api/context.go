package main

import (
	"net/http"
)

type ctxKey string

var (
	tokenPayloadKey = ctxKey("auth_payload")
)

func GetByContext[T any](r *http.Request, key ctxKey, defaultVal T) T {
	val, exist := r.Context().Value(key).(T)
	if !exist {
		return defaultVal
	}
	return val
}
