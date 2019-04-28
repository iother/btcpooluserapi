package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

var zlogger *zap.SugaredLogger

func init() {

	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.DebugLevel
	})

	var allCore []zapcore.Core

	consoleDebugging := zapcore.Lock(os.Stdout)

	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

	allCore = append(allCore, zapcore.NewCore(consoleEncoder, consoleDebugging, lowPriority))

	core := zapcore.NewTee(allCore...)

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	zlogger = logger.Sugar()
}

func InitLogger(isDebug bool) {

	fileName := "./logs/log"
	maxSize := 1 << 30
	maxBackups := 30
	maxAge := 1

	var compress = true
	// 打印错误级别的日志
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	// 打印所有级别的日志
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.DebugLevel
	})

	var allCore []zapcore.Core

	infoHook := lumberjack.Logger{
		Filename:   fileName + ".info",
		MaxSize:    maxSize, // megabytes
		MaxBackups: maxBackups,
		MaxAge:     maxAge,   //days
		Compress:   compress, // disabled by default
	}

	infoFileWriter := zapcore.AddSync(&infoHook)

	errHook := lumberjack.Logger{
		Filename:   fileName + ".err",
		MaxSize:    maxSize, // megabytes
		MaxBackups: maxBackups,
		MaxAge:     maxAge,   //days
		Compress:   compress, // disabled by default
	}

	errFileWriter := zapcore.AddSync(&errHook)

	consoleDebugging := zapcore.Lock(os.Stdout)

	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

	productionEncoderConfig := zap.NewProductionEncoderConfig()
	productionEncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	jsonConsoleEncoder := zapcore.NewConsoleEncoder(productionEncoderConfig)

	if isDebug {
		allCore = append(allCore, zapcore.NewCore(consoleEncoder, consoleDebugging, lowPriority))
	}

	allCore = append(allCore, zapcore.NewCore(jsonConsoleEncoder, errFileWriter, highPriority))

	allCore = append(allCore, zapcore.NewCore(consoleEncoder, infoFileWriter, lowPriority))

	core := zapcore.NewTee(allCore...)

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	zlogger = logger.Sugar()
}

func Debug(args ...interface{}) {
	zlogger.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	zlogger.Debugf(template, args...)
}

func Info(args ...interface{}) {
	zlogger.Info(args...)
}

func Infof(template string, args ...interface{}) {
	zlogger.Infof(template, args...)
}

func Warn(args ...interface{}) {
	zlogger.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	zlogger.Warnf(template, args...)
}

func Error(args ...interface{}) {
	zlogger.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	zlogger.Errorf(template, args...)
}

func Panic(args ...interface{}) {
	zlogger.Panic(args...)
}

func Panicf(template string, args ...interface{}) {
	zlogger.Panicf(template, args...)
}

func Fatal(args ...interface{}) {
	zlogger.Fatal(args...)
}

func Fatalf(template string, args ...interface{}) {
	zlogger.Fatalf(template, args...)
}
