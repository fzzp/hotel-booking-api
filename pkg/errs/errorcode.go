package errs

import (
	"net/http"

	"github.com/fzzp/gotk"
)

var (
	ErrOK                  = gotk.NewApiError(http.StatusOK, "200", "请求成功")
	ErrBadRequest          = gotk.NewApiError(http.StatusBadRequest, "400", "入参错误")
	ErrUnauthorized        = gotk.NewApiError(http.StatusUnauthorized, "401", "请先登陆")
	ErrForbidden           = gotk.NewApiError(http.StatusForbidden, "403", "禁止访问")
	ErrNotFound            = gotk.NewApiError(http.StatusNotFound, "404", "查无此记录")
	ErrMethodNotAllowed    = gotk.NewApiError(http.StatusMethodNotAllowed, "405", "请求方法不支持")
	ErrRecordExists        = gotk.NewApiError(http.StatusConflict, "409", "数据重复")
	ErrUnprocessableEntity = gotk.NewApiError(http.StatusUnprocessableEntity, "422", "请求无法处理")
	ErrTooManyRequests     = gotk.NewApiError(http.StatusTooManyRequests, "429", "请求繁忙")
	ErrServerError         = gotk.NewApiError(http.StatusInternalServerError, "500", "请求错误，请稍后重试！")
)
