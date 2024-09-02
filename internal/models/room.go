package models

import (
	"errors"

	"github.com/fzzp/hotel-booking-api/internal/dto"
)

type RoomStatus string

const (
	// NOTE: 如果修改这块类型, 记得更改 CheckRoomStatus 方法
	Available RoomStatus = "available" // 可用
	Occupied  RoomStatus = "occupied"  // 被占用
	Maintain  RoomStatus = "maintain"  // 维护/清洁
)

var (
	ErrRoomStatusNotSupported = errors.New("不支持该客房状态")
)

type Room struct {
	ID          uint   `db:"id" json:"id"`
	HotelID     uint   `db:"hotel_id" json:"hotelId"`
	RoomNo      string `db:"room_no" json:"roomNo"`
	Images      string `db:"images" json:"images"`
	Price       uint   `db:"price" json:"price"`
	RoomTypeID  int8   `db:"room_type_id" json:"roomTypeId"`
	Status      string `db:"status" json:"status"`
	Capacity    int8   `db:"capacity" json:"capacity"`
	Description string `db:"description" json:"description"`
	CreatedAt   string `db:"created_at" json:"createdAt"`
	UpdatedAt   string `db:"updated_at" json:"updatedAt"`
	IsDeleted   int8   `db:"is_deleted" json:"-"`
	RootType
}

func (r *Room) ToDto() dto.RoomResponse {
	var statusAsText = "维护"
	switch r.Status {
	case string(Available):
		statusAsText = "可用"
	case string(Occupied):
		statusAsText = "占用"
	}
	resp := dto.RoomResponse{
		ID:                  r.ID,
		HotelID:             r.HotelID,
		RoomNo:              r.RoomNo,
		Images:              r.Images,
		Price:               r.Price,
		RoomTypeID:          r.RoomTypeID,
		Status:              r.Status,
		StatusAsText:        statusAsText,
		Capacity:            r.Capacity,
		Description:         r.Description,
		CreatedAt:           r.CreatedAt,
		UpdatedAt:           r.UpdatedAt,
		RoomTypeName:        r.RoomTypeName,
		RoomTypeDescription: r.RoomTypeDescription,
	}

	return resp
}

type RootType struct {
	ID                  uint   `db:"id" json:"id"`
	RoomTypeName        string `db:"rt_name" json:"roomTypeName"`
	RoomTypeDescription string `db:"rt_description" json:"roomTypeDescription"`
}

func (r *Room) CheckRoomStatus(s string) error {
	var status = [...]RoomStatus{"available", "occupied", "maintain"}
	for _, v := range status {
		if v == RoomStatus(s) {
			return nil
		}
	}

	return ErrRoomStatusNotSupported
}
