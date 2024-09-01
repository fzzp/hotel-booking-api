package rdb

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/fzzp/gotk"
	"github.com/fzzp/hotel-booking-api/pkg/errs"
	"github.com/redis/go-redis/v9"
)

var (
	MaxSendNum   = 5
	dsTimeLayout = "2006-01-02 15:04:05"
)

type SMSModel struct {
	IsUsed   bool      `json:"isUsed"`    // 是否被使用
	SendNum  int       `json:"sendNum"`   // 当天发送次数
	Code     string    `json:"code"`      // 验证码
	ExpireAt time.Time `json:"expiredAt"` // 有效期
}

func NewSMSModel(code string, expireAt time.Time) *SMSModel {
	return &SMSModel{
		IsUsed:   false,
		SendNum:  1,
		Code:     code,
		ExpireAt: expireAt,
	}
}

func (sms *SMSModel) IsExpire(code string) *gotk.ApiError {
	if sms.Code != code {
		return errs.ErrBadRequest.AsMessage("验证不匹配")
	}
	if sms.IsUsed {
		return errs.ErrBadRequest.AsMessage("验证已失效")
	}
	if time.Now().After(sms.ExpireAt) {
		return errs.ErrBadRequest.AsMessage("验证已过期")
	}
	return nil
}

func (r *RedisRepo) SaveSMSCode(telephone string, sms *SMSModel) *gotk.ApiError {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	buf, err := json.Marshal(sms)
	if err != nil {
		return errs.ErrServerError.AsException(err)
	}

	dur, err := diffLastTime()
	if err != nil {
		return errs.ErrServerError.AsException(err)
	}

	_, err = r.Client.SetEx(ctx, smsCodeKey(telephone), string(buf), dur).Result()
	if err != nil {
		return errs.ErrServerError.AsException(err)
	}
	return nil
}

func (r *RedisRepo) GetSMSCode(telephone string) (*SMSModel, *gotk.ApiError) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	// // 多个命令一次提交
	// tx := r.Client.TxPipeline()

	// val, err := tx.Get(ctx, smsCodeKey(telephone)).Result()
	// if errors.Is(err, redis.Nil) {
	// 	tx.Discard() // 丢弃后面的命令，不用执行
	// 	return nil, errs.ErrNotFound
	// } else if err != nil {
	// 	tx.Discard() // 丢弃后面的命令，不用执行
	// 	return nil, errs.ErrServerError.AsException(err)
	// }

	// var sms SMSModel
	// err = json.Unmarshal([]byte(val), &sms)
	// if err != nil {
	// 	tx.Discard() // 丢弃后面的命令，不用执行
	// 	return nil, errs.ErrServerError.AsException(err)
	// }

	// sms.IsUsed = true
	// dataBytes, err := json.Marshal(sms)
	// if err != nil {
	// 	tx.Discard() // 丢弃后面的命令，不用执行
	// 	return nil, errs.ErrServerError.AsException(err)
	// }
	// dur, err := diffLastTime()
	// if err != nil {
	// 	tx.Discard() // 丢弃后面的命令，不用执行
	// 	return nil, errs.ErrServerError.AsException(err)
	// }
	// _, err = r.Client.SetNX(ctx, smsCodeKey(telephone), string(dataBytes), dur).Result()
	// if err != nil {
	// 	return nil, errs.ErrServerError.AsException(err)
	// }

	// if _, err := tx.Exec(ctx); err != nil {
	// 	return nil, errs.ErrServerError.AsException(err)
	// }

	// return &sms, nil

	// ====================================

	val, err := r.Client.Get(ctx, smsCodeKey(telephone)).Result()
	if errors.Is(err, redis.Nil) {
		return nil, errs.ErrNotFound
	} else if err != nil {
		return nil, errs.ErrServerError.AsException(err)
	}
	var sms SMSModel
	err = json.Unmarshal([]byte(val), &sms)
	if err != nil {
		return nil, errs.ErrServerError.AsException(err)
	}

	return &sms, nil
}

// diffLastTime 计算当前时间到当天23:59:59还有多久
func diffLastTime() (time.Duration, error) {
	lastTimeString := time.Now().Format("2006-01-02") + " 23:59:59"
	lastTime, err := time.Parse(dsTimeLayout, lastTimeString)
	if err != nil {
		return 0, err
	}
	return time.Until(lastTime), nil
}
