package logger

import (
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/fzzp/hotel-booking-api/pkg/config"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func InitLoagger(mode, level, logOutput string) {
	// 设置日志输出和分割
	logFile := &lumberjack.Logger{
		Filename:   logOutput,
		MaxSize:    100, // megabytes
		MaxBackups: 3,
		MaxAge:     28,    //days
		Compress:   false, // disabled by default
	}

	// 根据环境参数设置是否多输出（终端和文件输出）
	var mw io.Writer
	if mode == config.Development {
		mw = io.MultiWriter(os.Stdout, logFile)
	} else {
		mw = io.MultiWriter(logFile)
	}

	// 设置 logrus log风格、输出、日志等级
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000", // 包含毫秒数
	})
	logrus.SetOutput(mw)
	logrus.SetLevel(convertLogLevel(level))

	// 将logrus当作slog的扩展，组合使用
	sloger := slog.New(NewLogrusHandler(logrus.StandardLogger()))

	// 设为全局默认slog实例
	slog.SetDefault(sloger)
}

func convertLogLevel(level string) logrus.Level {
	var l logrus.Level

	switch strings.ToLower(level) {
	case "error":
		l = logrus.ErrorLevel
	case "warn":
		l = logrus.WarnLevel
	case "info":
		l = logrus.InfoLevel
	case "debug":
		l = logrus.DebugLevel
	default:
		l = logrus.InfoLevel
	}

	return l
}
