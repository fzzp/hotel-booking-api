package models

import "time"

type Booking struct {
	ID          uint      `db:"id" json:"id"`
	UserID      uint      `db:"user_id" json:"userId"`
	CheckIn     time.Time `db:"check_in" json:"checkIn"`
	CheckOut    time.Time `db:"check_out" json:"checkOut"`
	Members     int       `db:"members" json:"members"`
	TotalAmount int       `db:"total_amount" json:"totalAmount"`
	CreatedAt   string    `db:"created_at" json:"createdAt"`
	UpdatedAt   string    `db:"updated_at" json:"updatedAt"`
	IsDeleted   int8      `db:"is_deleted" json:"-"`
}
