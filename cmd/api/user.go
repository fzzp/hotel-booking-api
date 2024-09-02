package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/fzzp/gotk"
	"github.com/fzzp/gotk/token"
	"github.com/fzzp/hotel-booking-api/internal/dto"
	"github.com/fzzp/hotel-booking-api/pkg/errs"
)

var (
	sendSmsFailText = "发送失败"
	sendSmsSuccText = "发送成功"
)

func (app *application) SendSMSCode(w http.ResponseWriter, r *http.Request) {
	var req dto.TelephoneRequest
	if ok := app.ShouldBindJSON(w, r, &req); !ok {
		return
	}
	code, err := app.service.SMS.SendSMSCode(req.PhoneNumber)
	if err != nil {
		app.FAIL(w, r, err.AsMessage(sendSmsFailText))
		return
	}
	// TODO: 因为是简单的模拟，直接返回
	app.SUCC(w, r, pkg{"code": code, "message": sendSmsSuccText})
}

func (app *application) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequest
	if ok := app.ShouldBindJSON(w, r, &req); !ok {
		return
	}
	newID, err := app.service.User.CreateUser(req)
	if err != nil {
		app.FAIL(w, r, err)
		return
	}
	app.SUCC(w, r, newID)
}

func (app *application) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if ok := app.ShouldBindJSON(w, r, &req); !ok {
		return
	}
	req.UserAgent = r.UserAgent()
	req.ClientIP = r.RemoteAddr
	resp, err := app.service.User.UserLogin(req, &app.conf, app.jwt)
	if err != nil {
		app.FAIL(w, r, err)
		return
	}
	app.SUCC(w, r, resp)
}

func (app *application) GetProfileHandler(w http.ResponseWriter, r *http.Request) {
	payload := GetByContext[*token.Payload](r, tokenPayloadKey, nil)
	if payload == nil {
		app.FAIL(w, r, errs.ErrUnauthorized)
		return
	}
	uid, _ := strconv.Atoi(payload.UserText)
	if uid <= 0 {
		app.FAIL(w, r, errs.ErrNotFound.AsMessage("用户id不存在"))
		return
	}
	user, err := app.service.User.GetUserByID(uint(uid))
	if err != nil {
		app.FAIL(w, r, err)
		return
	}
	app.SUCC(w, r, user)
}

func (app *application) UpdateUserInfo(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateUserRequest
	if ok := app.ShouldBindJSON(w, r, &req); !ok {
		return
	}
	var payload *token.Payload
	payload = GetByContext(r, tokenPayloadKey, payload)
	if payload == nil {
		app.FAIL(w, r, errs.ErrServerError.AsMessage("登陆用户获取payload失败"))
		return
	}

	err := app.service.User.UpdateUser(req, payload)
	if err != nil {
		app.FAIL(w, r, err)
		return
	}
	app.SUCC(w, r, "修改成功")
}

func (app *application) RenewAccessToken(w http.ResponseWriter, r *http.Request) {
	req := struct {
		RefreshToken string `json:"refreshToken" validate:"required"`
	}{}
	if ok := app.ShouldBindJSON(w, r, &req); !ok {
		return
	}
	p, err := app.jwt.ParseToken(req.RefreshToken)
	if err != nil && errors.Is(err, gotk.ErrInvalidInput) {
		app.FAIL(w, r, errs.ErrBadRequest.AsException(err).AsMessage(err.Error()))
		return
	}
	if err != nil {
		app.FAIL(w, r, errs.ErrServerError.AsException(err).AsMessage("令牌无效"))
		return
	}
	resp, apiErr := app.service.User.UserRefreshToken(p, &app.conf, app.jwt)
	if apiErr != nil {
		app.FAIL(w, r, apiErr)
		return
	}
	app.SUCC(w, r, resp)
}
