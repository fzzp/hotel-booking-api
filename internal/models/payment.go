package models

import "time"

type Payment struct {
	ID            uint       `db:"id" json:"id"`
	BookingID     uint       `db:"booking_id" json:"bookingId"`
	PaymentType   int8       `db:"payment_type" json:"paymentType"`
	PaymentTime   *time.Time `db:"payment_time" json:"paymentTime"`
	PaymentAmount uint       `db:"payment_amount" json:"paymentAmount"`
	IsDeleted     int8       `db:"is_deleted" json:"-"`
}
