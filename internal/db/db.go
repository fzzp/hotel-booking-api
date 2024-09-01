package db

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type bindVar string

var (
	question bindVar = "?"
	// dollar   bindVar = "$"

	defaultTimeout = time.Second * 5
)

// 基础功能，提供给此包内使用

func create(db Queryable, querySQL string, args ...any) (uint, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	stmt, err := db.PrepareContext(ctx, querySQL)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, args...)
	if err != nil {
		return 0, err
	}

	_, err = res.RowsAffected()
	if err != nil {
		return 0, err
	}

	newId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	if newId <= 0 {
		return 0, ErrInsertFailed
	}

	return uint(newId), nil
}

func update(db Queryable, querySQL string, args ...any) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	stmt, err := db.PrepareContext(ctx, querySQL)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, args...)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows <= 0 {
		return ErrNotFound
	}
	return nil
}

// unqQuerySQL 唯一查询组合SQL
//
// fields 存储table唯一字段和对应的查询值，并且field对应的值不能是空字符串
func unqQuerySQL(fields map[string]string, bindvar bindVar) (sqlStr string, params []interface{}) {
	var list []string // sql 条件
	var index = 1
	for k, v := range fields {
		if v != "" {
			if bindvar == question {
				list = append(list, fmt.Sprintf("%s=?", k))
			} else {
				list = append(list, fmt.Sprintf("%s=$%d", k, index))
			}
			params = append(params, v)
			// 必须放条件内
			index++
		}
	}
	sql := " " + strings.Join(list, " and ") + " "
	return sql, params
}

// execTx 定义个执行事务公共的方法
func execTx(ctx context.Context, qb Queryable, fn func(*Repository) error) error {
	// qb => sqlx.DB/sqlx.Tx
	db, ok := qb.(*sqlx.DB)
	if !ok {
		return ErrNoEffectDB
	}
	// 开启事务
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	// 重新创建一个 Repository
	q := NewRepository(tx)
	if err = fn(q); err != nil {
		// 执行不成功，则 Rollback
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	// 执行成功 Commit
	return tx.Commit()
}
