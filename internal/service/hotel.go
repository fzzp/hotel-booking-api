package service

import (
	"github.com/fzzp/gotk"
	"github.com/fzzp/hotel-booking-api/internal/models"
)

type HotelService interface {
	GetAllHotels() ([]*models.Hotel, *gotk.ApiError)
	GetAllRoomTypes([]*models.RootType, *gotk.ApiError)
	GetRooms(hotelID uint) ([]*models.Room, *gotk.ApiError)
}

var _ HotelService = (*hotelService)(nil)

type hotelService struct {
}

func NewHotelService() HotelService {
	hs := &hotelService{}

	return hs
}

// GetAllHotels implements HotelService.
func (h *hotelService) GetAllHotels() ([]*models.Hotel, *gotk.ApiError) {
	panic("unimplemented")
}

// GetAllRoomTypes implements HotelService.
func (h *hotelService) GetAllRoomTypes([]*models.RootType, *gotk.ApiError) {
	panic("unimplemented")
}

// GetRooms implements HotelService.
func (h *hotelService) GetRooms(hotelID uint) ([]*models.Room, *gotk.ApiError) {
	panic("unimplemented")
}
