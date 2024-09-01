package db

import (
	"database/sql"
	"errors"

	"github.com/fzzp/gotk"
	"github.com/fzzp/hotel-booking-api/pkg/errs"
	"github.com/go-sql-driver/mysql"
)

// NOTE: 增加错误时，在 ConvertToApiError 添加对应的处理
var (
	ErrNotFound     = errors.New("记录不存在")
	ErrInsertFailed = errors.New("插入数据失败")
	ErrNoEffectDB   = errors.New("qb not is *sql.DB")
)

// ConvertToApiError 将db错误转换为 *gotk.ApiError
func ConvertToApiError(err error) *gotk.ApiError {
	// MySQL 错误码参照：
	// https://dev.mysql.com/doc/mysql-errors/8.0/en/server-error-reference.html

	if errors.Is(err, ErrNotFound) {
		return errs.ErrNotFound.AsException(err)
	}
	if errors.Is(err, ErrInsertFailed) || errors.Is(err, ErrNoEffectDB) {
		return errs.ErrServerError.AsException(err)
	}

	if err == sql.ErrNoRows {
		return errs.ErrNotFound.AsException(err)
	}

	dbErr, ok := err.(*mysql.MySQLError)
	if ok {
		if dbErr.Number == 1062 {
			return errs.ErrRecordExists.AsException(err)
		}
	}

	return errs.ErrServerError.AsException(err)
}
