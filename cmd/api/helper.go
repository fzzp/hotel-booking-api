package main

import (
	"log/slog"
	"net/http"

	"github.com/fzzp/gotk"
	"github.com/fzzp/hotel-booking-api/pkg/errs"
)

type pkg map[string]interface{}

func (app *application) ShouldBindJSON(w http.ResponseWriter, r *http.Request, dst interface{}) bool {
	err := gotk.ReadJSON(w, r, dst)
	if err != nil {
		app.FAIL(w, r, errs.ErrBadRequest.AsException(err))
		return false
	}

	slog.InfoContext(r.Context(), "arg", slog.Any("arg", dst))

	err = gotk.CheckStruct(dst)
	if err != nil {
		app.FAIL(w, r, errs.ErrBadRequest.AsException(err, err.Error()))
		return false
	}

	return true
}

// FAIL 请求失败
func (app *application) FAIL(w http.ResponseWriter, r *http.Request, err *gotk.ApiError) {
	slog.ErrorContext(
		r.Context(),
		r.Method+"-"+r.RequestURI,
		slog.String("err", err.Error()),
	)
	gotk.WriteJSON(w, r, err, nil)
}

// SUCC 请求成功
func (app *application) SUCC(w http.ResponseWriter, r *http.Request, data interface{}) {
	gotk.WriteJSON(w, r, errs.ErrOK, data)
}
