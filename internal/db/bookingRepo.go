package db

import "github.com/fzzp/hotel-booking-api/internal/models"

type BookingRepo interface {
	InsertOne(booking *models.Booking) (id uint, err error)
}

var _ BookingRepo = (*bookingRepo)(nil)

type bookingRepo struct {
	DB Queryable
}

func NewBookingRepo(qb Queryable) BookingRepo {
	return &bookingRepo{
		DB: qb,
	}
}

// InsertOne implements BookingRepo.
func (b *bookingRepo) InsertOne(booking *models.Booking) (id uint, err error) {
	panic("unimplemented")
}
