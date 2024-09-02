package dto

type RoomResponse struct {
	ID                  uint   `json:"id"`
	HotelID             uint   `djson:"hotelId"`
	RoomNo              string `json:"roomNo"`
	Images              string `json:"images"`
	Price               uint   `json:"price"`
	RoomTypeID          int8   `json:"roomTypeId"`
	Status              string `json:"status"`
	StatusAsText        string `json:"statusAsText"`
	Capacity            int8   `json:"capacity"`
	Description         string `json:"description"`
	CreatedAt           string `json:"createdAt"`
	UpdatedAt           string `json:"updatedAt"`
	RoomTypeName        string `json:"roomTypeName"`
	RoomTypeDescription string `json:"roomTypeDescription"`
}

type ListRoomResponse struct {
	List     *[]RoomResponse
	Metadata Metadata `json:"metadata"`
}
