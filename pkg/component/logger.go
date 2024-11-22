package component

import (
	"context"
	"github.com/sirupsen/logrus"
	"os"
	"rbd_proxy_dp/config"
)

type Logger struct{}

// NewLog -
func NewLog() *Logger {
	return &Logger{}
}

func (l *Logger) Start(ctx context.Context) error {
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(&logrus.TextFormatter{
		// 禁止重复时间戳
		DisableTimestamp: true,
		// 自定义时间格式
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true, // 确保显示完整时间
	})
	switch config.DefaultPublic().LogLevel {
	case "trace":
		logrus.SetLevel(logrus.TraceLevel)
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}
	return nil
}

func (l *Logger) CloseHandle() {}
