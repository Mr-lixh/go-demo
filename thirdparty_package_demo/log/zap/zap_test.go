package zap

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"k8s.io/apimachinery/pkg/util/uuid"
	"os"
	"testing"
)

// TestBase basic samples of zap.
func TestBase(t *testing.T) {
	var logger *zap.Logger
	// Create a logger use default production config.
	logger, _ = zap.NewProduction()
	logger.Info("this is a info message")
	logger.Error("this is a error message")

	// Create a sugared logger.
	sugared := logger.Sugar()
	sugared.Infof("this is a info message with args: name=%s", "test")
	sugared.Infow("this is a info message with key value", "name", "test", "age", 21)
}

func TestLogWithCustomerConfig(t *testing.T) {
	// customer writer
	writer := zapcore.AddSync(os.Stdout)

	// customer encoder
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05") // config time layout
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder                       // config level capital
	// customer encoder keys
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.MessageKey = "message"
	encoderConfig.LevelKey = "err_level"

	encoder := zapcore.NewJSONEncoder(encoderConfig)
	// encoder := zapcore.NewConsoleEncoder(encoderConfig)   // new console encoder

	// new core with customer encoder、writer and logger level
	core := zapcore.NewCore(encoder, writer, zapcore.InfoLevel)

	// new logger instance with core and options: AddCaller、Fields
	logger := zap.New(core, zap.AddCaller(), zap.Fields(zapcore.Field{
		Key:    "uid",
		Type:   zapcore.StringType,
		String: string(uuid.NewUUID()),
	}))

	logger.Info("this is a info message")
	logger.Error("this is a error message")

	sugared := logger.Sugar()
	sugared.Infof("this is a info message with args: name=%s", "test")
	sugared.Infow("this is a info message with key value", "name", "test", "age", 21)
}
