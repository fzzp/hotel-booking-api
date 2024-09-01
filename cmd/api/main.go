package main

import (
	"log"
	"log/slog"
	"sync"

	"github.com/fzzp/gotk"
	"github.com/fzzp/gotk/token"
	"github.com/fzzp/hotel-booking-api/internal/db"
	"github.com/fzzp/hotel-booking-api/internal/rdb"
	"github.com/fzzp/hotel-booking-api/internal/service"
	"github.com/fzzp/hotel-booking-api/pkg/config"
	"github.com/fzzp/hotel-booking-api/pkg/logger"
)

type application struct {
	conf    config.Config
	wg      sync.WaitGroup
	jwt     token.Maker
	service *service.DefaultService
}

func main() {
	// 优先启动验证器，解析json配置需要使用其验证
	gotk.InitValidation("zh")
	// 加载应用配置
	conf, err := config.LoadConfig("config.local.json")
	if err != nil {
		log.Fatal(err)
	}

	// 打印一下看看
	// conf.Println()

	// 初始化日志输出
	logger.InitLoagger(conf.Mode, conf.Log.Level, conf.Log.LogOutput)

	// 创建jwt管理者
	jwt, err := token.NewJWTMaker(conf.Token.SecretKey, conf.Token.Issuer)
	if err != nil {
		log.Fatal(err)
	}

	// 声明 application 赋值
	app := application{
		conf: conf,
		jwt:  jwt,
	}

	// 创建数据库链接，配置基本参数
	conn := db.NewSQLxDb("mysql", conf.DBSource)
	conn.DB.SetMaxOpenConns(conf.Database.MaxOpenConn)
	conn.DB.SetMaxIdleConns(conf.Database.MaxIdleConn)
	conn.DB.SetConnMaxIdleTime(conf.Database.MaxIdleTime.Duration)

	// 创建redis client
	rClient, err := rdb.NewRedisClient(conf.Redis.Addr, conf.Redis.Password, conf.Redis.DB)
	if err != nil {
		log.Fatal(err)
	}

	// 创建服务层
	dbRepo := db.NewRepository(conn)
	rRepo := rdb.NewRedisRepo(rClient)
	app.service = service.NewDefaultService(dbRepo, rRepo)

	if err = app.serve(); err != nil {
		slog.Error("stop server", slog.String("err", err.Error()))
	}

	// start := time.Now()
	// for i := 0; i < 100000; i++ {
	// 	slog.Info("init app", slog.Any("conf", app.conf.DBSource))
	// }
	// fmt.Println(time.Since(start)) // 2.8s 3s 2.8 3.1 2.9
}
