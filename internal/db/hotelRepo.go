package db

import (
	"fmt"

	"github.com/fzzp/hotel-booking-api/internal/models"
)

type HotelRepo interface {
	// 酒店相关
	InsertHotel(hotel *models.Hotel) (id uint, err error)
	GetHotelById(uint) (*models.Hotel, error)
	UpdateHotel(id uint, hotel *models.Hotel) error
	GetAllHotels() ([]*models.Hotel, error)

	// 客房相关
	InserRoom(models.Room) (uint, error)
	UpdateRoom(models.Room) error
	UpdateRoomType(roomID, typeID uint) error
	UpdateRoomStatus(roomId uint, status string) error
	GetRoomListByHotelID(hotelID uint, f Filter) ([]*models.Room, Metadata, error)
}

var _ HotelRepo = (*hotelRepo)(nil)

type hotelRepo struct {
	DB Queryable
}

func NewHotelRepo(qb Queryable) *hotelRepo {
	return &hotelRepo{
		DB: qb,
	}
}

// GetRoomListByHotelID implements HotelRepo.
func (h *hotelRepo) GetRoomListByHotelID(hotelID uint, f Filter) ([]*models.Room, Metadata, error) {
	sqlCount := `select count(*) as total from rooms where hotel_id = ? and is_deleted=1`
	sqlRoom := `
		select
			r.id,
			r.hotel_id,
			r.room_no,
			r.images,
			r.price,
			r.capacity,
			r.status,
			r.room_type_id,
			r.description,
			r.created_at,
			r.updated_at,
			rt.rt_name,
			rt.rt_description
		from rooms r
		inner join room_types rt on r.room_type_id = rt.id
		where 
			hotel_id = ? and is_deleted=1
	`
	sum := struct {
		Total int `db:"total"`
	}{}
	err := h.DB.Get(&sum, sqlCount, hotelID)
	if err != nil {
		return nil, Metadata{}, err
	}

	list := make([]*models.Room, 0)

	sqlRoom += f.sortSQL()
	sqlRoom += f.limitSQL()

	fmt.Println("sqlRoom: ", sqlRoom)
	err = h.DB.Select(&list, sqlRoom, hotelID)
	if err != nil {
		return nil, Metadata{}, err
	}
	metadata := calculateMetadata(sum.Total, f.PageInt, f.PageSize)
	return list, metadata, nil
}

// InserRoom implements HotelRepo.
func (h *hotelRepo) InserRoom(models.Room) (uint, error) {
	panic("unimplemented")
}

// UpdateRoom implements HotelRepo.
func (h *hotelRepo) UpdateRoom(models.Room) error {
	panic("unimplemented")
}

// UpdateRoomStatus implements HotelRepo.
func (h *hotelRepo) UpdateRoomStatus(roomId uint, status string) error {
	panic("unimplemented")
}

// UpdateRoomType implements HotelRepo.
func (h *hotelRepo) UpdateRoomType(roomID uint, typeID uint) error {
	panic("unimplemented")
}

// GetHotelById 获取一个酒店基础信息
func (h *hotelRepo) GetHotelById(id uint) (*models.Hotel, error) {
	sql := `select id, name, address, logo, created_at, updated_at from hotels where id=? and is_deleted=1`
	var hotel models.Hotel
	err := h.DB.Get(&hotel, sql, id)
	if err != nil {
		return nil, err
	}
	return &hotel, nil
}

// InsertHotel 添加一个酒店
func (h *hotelRepo) InsertHotel(hotel *models.Hotel) (id uint, err error) {
	sql := `insert into hotels(name, address, logo) values(?,?,?);`
	return create(h.DB, sql, hotel.Name, hotel.Address, hotel.Logo)
}

// UpdateHotel 更新一个酒店
func (h *hotelRepo) UpdateHotel(id uint, hotel *models.Hotel) error {
	sql := `update hotel set 
		name=?, address=?, logo=?, updated_at=now() where id = ? and is_deleted = 1;`
	return update(h.DB, sql, hotel.Name, hotel.Address, hotel.Logo, hotel.ID)
}

// GetAllHotels 获取所有酒店
func (h *hotelRepo) GetAllHotels() ([]*models.Hotel, error) {
	sql := `select 
		id, name, address, logo, created_at, updated_at
	from hotels where is_deleted=1;`
	hotels := make([]*models.Hotel, 0)
	err := h.DB.Select(&hotels, sql)
	return hotels, err
}
