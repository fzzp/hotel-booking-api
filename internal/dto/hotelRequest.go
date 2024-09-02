package dto

type AddHotelRequest struct {
	Name    string `json:"name" validate:"required"`
	Address string `json:"address" validate:"required"`
	Logo    string `json:"logo" validate:"required"`
}
