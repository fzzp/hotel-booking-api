package main

import (
	"net/http"

	"github.com/fzzp/gotk"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	// 中间件执行顺序先进先执行；但是中间的defer还是按照后进后执行
	// 因此 AccessLogger 的defer函数要保证最后执行就放第一位
	mux.Use(app.AccessLogger)
	mux.Use(app.RecoverPanic)
	mux.Use(middleware.CleanPath)

	// 跨域中间件
	mux.Use(app.EnableCORS())

	v1 := chi.NewRouter()

	v1.Get("/healthz", app.HealthzHandler)
	v1.Post("/signup", app.SignUpHandler)
	v1.Post("/login", app.LoginHandler)
	v1.Post("/sms", app.SendSMSCode)
	v1.Post("/renewToken", app.RenewAccessToken)
	v1.Get("/hotels", app.GetHotels)
	v1.Get("/getRoomsByHotel/{id:[0-9]+}", app.GetRoomsHandler) // id?pageInt=1&pageSize=20

	// 需要登陆认证接口
	v1.Group(func(r chi.Router) {
		r.Use(app.RequiredAuth)
		r.Get("/profile", app.GetProfileHandler)
		r.Post("/updateUser", app.UpdateUserInfo)

	})

	// 需要管理员权限的接口
	v1.Group(func(r chi.Router) {
		r.Use(app.RequiredAuth)    // 先验证登陆
		r.Use(app.RequiredRole(1)) // 验证是否是管理员
		r.Post("/addHotel", app.AddHotelHandler)
	})

	// 附加/挂载到mux上,方便版本维护
	mux.Mount("/api/v1", gotk.SetVersionCtx(v1, "v1"))

	// 设置requesst_id，提供给其他中间使用
	return gotk.SetRequestIDCtx(mux)
}
