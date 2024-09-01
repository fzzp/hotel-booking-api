package db

import "github.com/fzzp/hotel-booking-api/internal/models"

type HotelRepo interface {
	// 酒店相关
	InsertHotel(hotel *models.Hotel) (id uint, err error)
	GetHotelById(uint) (*models.Hotel, error)
	UpdateHotel(id uint, hotel *models.Hotel) error

	// 客房相关
	InserRoom(models.Room) (uint, error)
	UpdateRoom(models.Room) error
	UpdateRoomType(roomID, typeID uint) error
	UpdateRoomStatus(roomId uint, status string) error
	GetRoomListByHotelID(hotelID uint, f Filter) ([]*models.Room, error)
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
func (h *hotelRepo) GetRoomListByHotelID(hotelID uint, f Filter) ([]*models.Room, error) {
	panic("unimplemented")
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
