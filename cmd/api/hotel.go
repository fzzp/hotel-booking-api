package main

import (
	"net/http"
	"strconv"

	"github.com/fzzp/hotel-booking-api/internal/dto"
	"github.com/fzzp/hotel-booking-api/pkg/errs"
	"github.com/go-chi/chi/v5"
)

func (app *application) AddHotelHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.AddHotelRequest
	if ok := app.ShouldBindJSON(w, r, &req); !ok {
		return
	}
	newID, err := app.service.Hotel.CreateHotel(req)
	if err != nil {
		app.FAIL(w, r, err)
		return
	}
	app.SUCC(w, r, newID)
}

func (app *application) GetHotels(w http.ResponseWriter, r *http.Request) {
	list, err := app.service.Hotel.GetAllHotels()
	if err != nil {
		app.FAIL(w, r, err)
		return
	}
	app.SUCC(w, r, list)
}

func (app *application) GetRoomsHandler(w http.ResponseWriter, r *http.Request) {
	idTxt := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idTxt)
	if id <= 0 {
		app.FAIL(w, r, errs.ErrBadRequest.AsMessage("无效的id"))
		return
	}
	pageIntTx := r.URL.Query().Get("pageInt")
	pageSizeTx := r.URL.Query().Get("pageSize")
	pageInt, _ := strconv.Atoi(pageIntTx)
	pageSize, _ := strconv.Atoi(pageSizeTx)

	resp, err := app.service.Hotel.GetRooms(uint(id), pageInt, pageSize)
	if err != nil {
		app.FAIL(w, r, err)
		return
	}

	app.SUCC(w, r, resp)

}
