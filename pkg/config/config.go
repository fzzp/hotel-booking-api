package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fzzp/gotk"
)

const (
	Development = "development"
	Production  = "production"
)

type PzzDuration struct {
	Duration time.Duration
}

// UnmarshalJSON 实现 json.Unmarshaler 接口
func (p *PzzDuration) UnmarshalJSON(input []byte) error {
	// 去除可能存在的双引号
	s := strings.Trim(string(input), `"`)
	d, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	p.Duration = d
	return nil
}

// MarshalJSON 实现 json.Marshaler 接口
func (p *PzzDuration) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.Duration.String())
}

// Config 整个应用配置
type Config struct {
	Port     int    `json:"port" validate:"required,min=1000,max=65535"`
	Mode     string `json:"mode" validate:"required,oneof=development production"`
	DBSource string // 非json配置，组装好给内部代码使用
	Database struct {
		DBHost      string      `json:"dbHost" validate:"required,ip"`
		DBPort      int         `json:"dbPort" validate:"required,min=1000,max=65535"`
		DBUser      string      `json:"dbUser" validate:"required"`
		DBPswd      string      `json:"dbPswd" validate:"required"`
		DBName      string      `json:"dbName" validate:"required"`
		MaxOpenConn int         `json:"maxOpenConn" validate:"required,min=5,max=100"`
		MaxIdleConn int         `json:"maxIdleConn" validate:"required,min=5,max=100"`
		MaxIdleTime PzzDuration `json:"maxIdleTime" validate:"required"`
	}
	Token struct {
		SecretKey      string      `json:"secretKey" validate:"required,min=32"`
		Issuer         string      `json:"issuer" validate:"required"`
		ATokenDuration PzzDuration `json:"aTokenDuration" validate:"required"`
		RTokenDuration PzzDuration `json:"rTokenDuration" validate:"required"`
	}
	Log struct {
		Level     string `json:"level" validate:"required,oneof=info debug warn error"`
		LogOutput string `json:"logOutput" validate:"required"`
	}
	Redis struct {
		Addr     string `json:"addr" validate:"required"`
		Password string `json:"password" validate:"-"`
		DB       int    `json:"db" validate:"required,min=0,max=16"`
	}
}

func (c Config) Println() {
	buf, _ := json.MarshalIndent(&c, "", " ")
	fmt.Println(string(buf))
}

// LoadConfig 加载配置
func LoadConfig(filename string) (Config, error) {
	var conf Config
	buf, err := os.ReadFile(filename)
	if err != nil {
		return conf, err
	}

	if err = json.Unmarshal(buf, &conf); err != nil {
		return conf, err
	}

	if err = gotk.CheckStruct(&conf); err != nil {
		return conf, err
	}

	conf.DBSource = fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.Database.DBUser,
		conf.Database.DBPswd,
		conf.Database.DBHost,
		conf.Database.DBPort,
		conf.Database.DBName,
	)

	return conf, nil
}
