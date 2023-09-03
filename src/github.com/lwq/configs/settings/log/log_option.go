package log

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

type LogOption struct {
	directory string
}

// 设置全局log
func (l *LogOption) SetGlobalLogger() {
	l.initGlobalLogger(time.Now().Format("2006-01-02"))
	go l.updateLogFileName()
}

func (l *LogOption) updateLogFileName() {
	now := time.Now()
	nextMidnight := now.Add(time.Hour*24 - time.Duration(now.Hour())*time.Hour - time.Duration(now.Minute())*time.Minute - time.Duration(now.Second())*time.Second)
	interval := nextMidnight.Sub(now)
	Logger.Info(fmt.Sprintf("[1]距离下次更新日志需要：%d", interval))
	timer := time.NewTimer(interval)
	defer timer.Stop()
	for range timer.C {
		l.initGlobalLogger(time.Now().Format("2006-01-02"))
		nextMidnight = now.Add(time.Hour*24 - time.Duration(now.Hour())*time.Hour - time.Duration(now.Minute())*time.Minute - time.Duration(now.Second())*time.Second)
		interval = nextMidnight.Sub(now)
		Logger.Info(fmt.Sprintf("[2]距离下次更新日志需要：%d", interval))
		timer.Reset(interval)
	}
}

func (l *LogOption) initGlobalLogger(date string) {
	if _, err := os.Stat(l.directory); os.IsNotExist(err) {
		// 目录不存在，创建它
		err := os.MkdirAll(l.directory, os.ModePerm)
		if err != nil {
			panic("Error creating directory:" + err.Error())
		}
	}
	cfg := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),              // 设置日志级别
		OutputPaths:      []string{l.directory + "info-" + date + ".log"},  // 设置日志文件路径
		ErrorOutputPaths: []string{l.directory + "error-" + date + ".log"}, // 设置错误日志文件路径
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder, // 也可以使用其他时间格式编码器
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
	}
	var err error
	Logger, err = cfg.Build()
	if err != nil {
		panic(fmt.Sprintf("init logger error:%s", err.Error()))
	}
	defer Logger.Sync() // 确保日志缓冲区中的所有日志都已写入文件
}

func NewLogOption() *LogOption {
	return &LogOption{
		directory: viper.GetString("log.directory"),
	}
}
