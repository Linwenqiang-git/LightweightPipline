package log

import (
	"os"
	"time"

	kratoszap "github.com/go-kratos/kratos/contrib/log/zap/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogOption struct {
	logPath        string
	isOpenInternal bool
}

// 设置全局log
func (l *LogOption) SetGlobalLogger() {
	logfile := l.logPath + time.Now().Format("2006-01-02") + ".log"
	f, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return
	}
	writeSyncer := zapcore.AddSync(f)

	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)
	z := zap.New(core)

	logger := kratoszap.NewLogger(z)
	log.SetLogger(logger)
	if !l.isOpenInternal {
		go l.updateLogFileName()
	}
	l.isOpenInternal = true
}

func (l *LogOption) updateLogFileName() {
	// 计算距离下一个0点的时间间隔
	now := time.Now()
	nextMidnight := now.Add(time.Hour*24 - time.Duration(now.Hour())*time.Hour - time.Duration(now.Minute())*time.Minute - time.Duration(now.Second())*time.Second)
	interval := nextMidnight.Sub(now)

	// 创建计时器
	timer := time.NewTimer(interval)
	defer timer.Stop()

	// 启动循环
	for range timer.C {
		// 更新日志文件并计算下一个0点的时间间隔
		l.SetGlobalLogger()
		nextMidnight = now.Add(time.Hour*24 - time.Duration(now.Hour())*time.Hour - time.Duration(now.Minute())*time.Minute - time.Duration(now.Second())*time.Second)
		interval = nextMidnight.Sub(now)
		timer.Reset(interval)
	}
}

func NewLogOption() *LogOption {
	return &LogOption{
		logPath:        viper.GetString("log.logPath"),
		isOpenInternal: false,
	}
}
