package service

import (
	"github.com/fzzp/gotk"
	"github.com/fzzp/hotel-booking-api/internal/db"
	"github.com/fzzp/hotel-booking-api/internal/dto"
	"github.com/fzzp/hotel-booking-api/internal/models"
)

type HotelService interface {
	CreateHotel(req dto.AddHotelRequest) (uint, *gotk.ApiError)
	GetAllHotels() ([]*models.Hotel, *gotk.ApiError)
	GetAllRoomTypes([]*models.RootType, *gotk.ApiError)
	GetRooms(hotelID uint, pageInt, pageSize int) (pkg, *gotk.ApiError)
}

var _ HotelService = (*hotelService)(nil)

type hotelService struct {
	hotel db.HotelRepo
}

// CreateHotel implements HotelService.
func (h *hotelService) CreateHotel(req dto.AddHotelRequest) (uint, *gotk.ApiError) {
	hotel := &models.Hotel{
		Name:    req.Name,
		Address: req.Address,
		Logo:    req.Logo,
	}
	id, err := h.hotel.InsertHotel(hotel)
	if err != nil {
		return 0, db.ConvertToApiError(err)
	}
	return id, nil
}

func NewHotelService(hotel db.HotelRepo) HotelService {
	hs := &hotelService{
		hotel: hotel,
	}

	return hs
}

// GetAllHotels implements HotelService.
func (h *hotelService) GetAllHotels() ([]*models.Hotel, *gotk.ApiError) {
	list, err := h.hotel.GetAllHotels()
	if err != nil {
		return nil, db.ConvertToApiError(err)
	}
	return list, nil
}

// GetAllRoomTypes implements HotelService.
func (h *hotelService) GetAllRoomTypes([]*models.RootType, *gotk.ApiError) {

}

// GetRooms implements HotelService.
func (h *hotelService) GetRooms(hotelID uint, pageInt, pageSize int) (pkg, *gotk.ApiError) {
	f := db.Filter{
		PageInt:    pageInt,
		PageSize:   pageSize,
		SortFields: []string{"id"},
		SortSafeFields: []string{
			"id", "-id",
			"price", "-price",
			"capacity", "-capacity",
			"room_type_id", "-room_type_id",
			"status", "-status",
			"updated_at", "-updated_at",
		},
	}
	list, metadata, err := h.hotel.GetRoomListByHotelID(hotelID, f)
	if err != nil {
		return pkg{}, db.ConvertToApiError(err)
	}

	return pkg{"list": list, "metadata": metadata}, nil
}
