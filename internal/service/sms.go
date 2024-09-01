package service

import (
	"errors"
	"strconv"
	"time"

	"github.com/fzzp/gotk"
	"github.com/fzzp/hotel-booking-api/internal/rdb"
	"github.com/fzzp/hotel-booking-api/pkg/errs"
)

type SmsService interface {
	SendSMSCode(phone string) (int64, *gotk.ApiError)
}

var _ SmsService = (*smsService)(nil)

type smsService struct {
}

func NewSmsService() smsService {
	return smsService{}
}

func (s smsService) SendSMSCode(phone string) (int64, *gotk.ApiError) {
	// 尝试获取获取
	smsCode, err := rRepo.GetSMSCode(phone)
	// 存在错误并且不是查询不到错误
	if err != nil && !errors.Is(err, errs.ErrNotFound) {
		return 0, err
	}

	if smsCode != nil && smsCode.SendNum > 5 {
		return 0, errs.ErrForbidden.AsMessage("今天发送次数已超过5次，你已被限制")
	}

	code := gotk.RandomInt(1000, 9999)
	// 有效时长5分钟
	sms := rdb.NewSMSModel(strconv.Itoa(int(code)), time.Now().Add(5*time.Minute))
	if smsCode != nil { // 累加当天发送次数
		sms.SendNum += smsCode.SendNum
	}
	err = rRepo.SaveSMSCode(phone, sms)
	return code, err
}
