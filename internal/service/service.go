package service

import (
	"github.com/fzzp/hotel-booking-api/internal/db"
	"github.com/fzzp/hotel-booking-api/internal/rdb"
)

var rRepo *rdb.RedisRepo

type stringMap map[string]string
type pkg map[string]any

type DefaultService struct {
	SMS          SmsService
	User         UserService
	Hotel        HotelService
	Booking      BookingService
	Payment      PaymentService
	Notification NotificationService
}

func NewDefaultService(repo *db.Repository, cache *rdb.RedisRepo) *DefaultService {
	rRepo = cache // 保存局部变量，包内使用
	return &DefaultService{
		SMS:          NewSmsService(),
		User:         NewUserService(repo.UserRepo, repo.SessionRepo),
		Hotel:        NewHotelService(repo.HotelRepo),
		Booking:      NewBookingService(),
		Payment:      NewPaymentService(),
		Notification: NewNotificationService(),
	}
}
