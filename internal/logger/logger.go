package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	DebugLevel = "debug"
	InfoLevel  = "info"
	WarnLevel  = "warn"
	ErrorLevel = "error"
	FatalLevel = "fatal"
)

var (
	log *zap.SugaredLogger
)

func Init(level string) {
	var zapLevel zapcore.Level
	switch level {
	case DebugLevel:
		zapLevel = zapcore.DebugLevel
	case InfoLevel:
		zapLevel = zapcore.InfoLevel
	case WarnLevel:
		zapLevel = zapcore.WarnLevel
	case ErrorLevel:
		zapLevel = zapcore.ErrorLevel
	case FatalLevel:
		zapLevel = zapcore.FatalLevel
	default:
		zapLevel = zapcore.InfoLevel
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		zapLevel,
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	log = logger.Sugar()
}

func Debug(args ...interface{}) {
	if log == nil {
		initDefaultLogger()
	}
	log.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	if log == nil {
		initDefaultLogger()
	}
	log.Debugf(format, args...)
}

func Info(args ...interface{}) {
	if log == nil {
		initDefaultLogger()
	}
	log.Info(args...)
}

func Infof(format string, args ...interface{}) {
	if log == nil {
		initDefaultLogger()
	}
	log.Infof(format, args...)
}

func Warn(args ...interface{}) {
	if log == nil {
		initDefaultLogger()
	}
	log.Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	if log == nil {
		initDefaultLogger()
	}
	log.Warnf(format, args...)
}

func Error(args ...interface{}) {
	if log == nil {
		initDefaultLogger()
	}
	log.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	if log == nil {
		initDefaultLogger()
	}
	log.Errorf(format, args...)
}

func Fatal(args ...interface{}) {
	if log == nil {
		initDefaultLogger()
	}
	log.Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	if log == nil {
		initDefaultLogger()
	}
	log.Fatalf(format, args...)
}

func With(args ...interface{}) *zap.SugaredLogger {
	if log == nil {
		initDefaultLogger()
	}
	return log.With(args...)
}

func Sync() error {
	if log == nil {
		return nil
	}
	return log.Sync()
}

func initDefaultLogger() {
	fmt.Println("Warning: Using default logger. Call logger.Init() to configure.")
	Init(InfoLevel)
}
