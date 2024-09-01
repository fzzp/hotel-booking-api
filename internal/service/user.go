package service

import (
	"database/sql"
	"errors"
	"strconv"
	"time"

	"github.com/fzzp/gotk"
	"github.com/fzzp/gotk/token"
	"github.com/fzzp/hotel-booking-api/internal/db"
	"github.com/fzzp/hotel-booking-api/internal/dto"
	"github.com/fzzp/hotel-booking-api/internal/models"
	"github.com/fzzp/hotel-booking-api/pkg/config"
	"github.com/fzzp/hotel-booking-api/pkg/errs"
	"github.com/fzzp/hotel-booking-api/util"
)

var (
	defaultPswd      = "87654321" // md5: 1ba4413ca86ad65f676579cdf83d6752
	defUserNamePre   = "用户_"
	smsExpireMessage = "验证码无效或已过期"
)

type UserService interface {
	CreateUser(input dto.CreateUserRequest) (uint, *gotk.ApiError)
	GetUserByID(uid uint) (*models.User, *gotk.ApiError)
	GetUserByPhoneNumber(phone string) (*models.User, *gotk.ApiError)
	UserLogin(req dto.LoginRequest, conf *config.Config, jwt token.Maker) (*dto.LoginResponse, *gotk.ApiError)
	UserRefreshToken(payload *token.Payload, conf *config.Config, jwt token.Maker) (*dto.RenewTokenResponse, *gotk.ApiError)
}

type userSerive struct {
	store db.UserRepo
	sess  db.SessionRepo
}

var _ UserService = (*userSerive)(nil)

func NewUserService(store db.UserRepo, sess db.SessionRepo) UserService {
	return userSerive{
		store: store,
		sess:  sess,
	}
}

func (s userSerive) CreateUser(input dto.CreateUserRequest) (uint, *gotk.ApiError) {
	// TODO: 真实短信验证码

	// 查询短信验证码
	sms, apiErr := rRepo.GetSMSCode(input.PhoneNumber)
	if apiErr != nil && errors.Is(apiErr, errs.ErrNotFound) {
		return 0, apiErr.AsMessage(smsExpireMessage)
	}
	if apiErr != nil {
		return 0, apiErr
	}

	if apiErr = sms.IsExpire(input.SMSCode); apiErr != nil {
		return 0, apiErr
	}

	// 查询手机号是否已经注册
	u, err := s.store.GetOneByUq(stringMap{"phone_number": input.PhoneNumber})
	if err != nil && !errors.Is(err, sql.ErrNoRows) { // 查询错误
		return 0, db.ConvertToApiError(err)
	}
	if err == nil { // 查询到用户数据，手机已经注册
		return 0, errs.ErrRecordExists
	}

	if u != nil && u.ID > 0 {
		return 0, errs.ErrRecordExists
	}

	// >>> 执行注册逻辑 <<<

	// 验证密码和用户名
	if input.PasswordHash == "" {
		input.PasswordHash = util.MD5(util.MD5(defaultPswd)) // 使用两次MD5转换，前端也同样需要执行2次MD5
	}
	if input.Username == "" {
		input.Username = defUserNamePre + gotk.RandomString(3)
	}
	// 加密
	pwsh, err := util.Hash(input.PasswordHash)
	if err != nil {
		return 0, errs.ErrBadRequest.AsException(err)
	}
	user := models.User{
		PhoneNumber:  input.PhoneNumber,
		PasswordHash: pwsh,
		Username:     input.Username,
		Avatar:       input.Avatar,
	}
	newID, err := s.store.InsertOne(&user)
	if err != nil {
		return 0, db.ConvertToApiError(err)
	}
	return newID, nil
}

func (s userSerive) GetUserByID(uid uint) (*models.User, *gotk.ApiError) {
	user, err := s.store.GetOneByUq(stringMap{"id": strconv.Itoa(int(uid))})
	if err != nil {
		return nil, db.ConvertToApiError(err)
	}
	return user, nil
}

func (s userSerive) GetUserByPhoneNumber(phone string) (*models.User, *gotk.ApiError) {
	user, err := s.store.GetOneByUq(stringMap{"phone_number": phone})
	if err != nil {
		return nil, db.ConvertToApiError(err)
	}
	return user, nil
}

// UserLogin 用户登陆逻辑
func (s userSerive) UserLogin(req dto.LoginRequest, conf *config.Config, jwt token.Maker) (*dto.LoginResponse, *gotk.ApiError) {
	// 判断是短信还是密码登陆
	if req.LoginType == "sms" { // 短信验证码登陆，先验证
		sms, apiErr := rRepo.GetSMSCode(req.PhoneNumber)
		if apiErr != nil && errors.Is(apiErr, errs.ErrNotFound) {
			return nil, apiErr.AsMessage(smsExpireMessage)
		}
		if apiErr != nil {
			return nil, apiErr
		}
		if apiErr = sms.IsExpire(req.SMSCode); apiErr != nil {
			return nil, apiErr
		}
	}

	// 获取用户信息
	user, err := s.store.GetOneByUq(stringMap{"phone_number": req.PhoneNumber})
	if err != nil {
		return nil, db.ConvertToApiError(err)
	}

	if req.LoginType == "psw" {
		if err := util.Matches(req.PasswordHash, user.PasswordHash); err != nil {
			return nil, errs.ErrBadRequest.AsMessage("密码不匹配")
		}
	}

	aPayload := token.NewPayload(strconv.Itoa(int(user.ID)), conf.Token.ATokenDuration.Duration)
	rPayload := token.NewPayload(strconv.Itoa(int(user.ID)), conf.Token.RTokenDuration.Duration)

	aToken, err := jwt.GenToken(aPayload)
	rToken, err2 := jwt.GenToken(rPayload)
	if err != nil || err2 != nil {
		return nil, errs.ErrServerError.AsException(err).AsException(err2)
	}

	// 存储refreshToken
	ss := models.Session{
		UserID:       user.ID,
		TokenID:      rPayload.ID,
		RefreshToken: rToken,
		ClientIP:     req.ClientIP,
		UserAgent:    req.UserAgent,
		ExpiresAt:    rPayload.ExpiredAt,
	}
	_, err = s.sess.InsertOne(&ss)
	if err != nil {
		return nil, db.ConvertToApiError(err)
	}

	resp := user.ToDto(aToken, rToken)
	return &resp, nil
}

func (s userSerive) UserRefreshToken(payload *token.Payload, conf *config.Config, jwt token.Maker) (*dto.RenewTokenResponse, *gotk.ApiError) {
	ss, err := s.sess.GetOneByUq(map[string]string{"token_id": payload.ID})
	if err != nil {
		return nil, db.ConvertToApiError(err)
	}

	if time.Now().After(ss.ExpiresAt) {
		return nil, errs.ErrUnauthorized.AsMessage("refresh已过期")
	}
	puid, _ := strconv.Atoi(payload.UserText)
	if ss.UserID != uint(puid) {
		return nil, errs.ErrForbidden.AsMessage("refreshToken已被篡改")
	}

	idStr := strconv.Itoa(int(ss.UserID))
	aPayload := token.NewPayload(idStr, conf.Token.ATokenDuration.Duration)
	aToken, err := jwt.GenToken(aPayload)
	if err != nil {
		return nil, errs.ErrServerError.AsException(err).AsException(err)
	}
	resp := dto.RenewTokenResponse{
		AccessToken: aToken,
		ExpiresAt:   aPayload.ExpiredAt,
	}
	return &resp, nil
}
