package service

type BookingService interface {
}

var _ BookingService = (*bookingService)(nil)

type bookingService struct {
}

func NewBookingService() BookingService {
	bs := &bookingService{}

	return &bs
}
