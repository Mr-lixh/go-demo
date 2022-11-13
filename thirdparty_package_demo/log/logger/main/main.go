package main

import (
	"flag"
	"github.com/Mr-lixh/go-demo/thirdparty_package_demo/log/logger"
	"go.uber.org/zap/zapcore"
)

func main() {
	var logOptions = &logger.Options{}

	logOptions.BindFlags(flag.CommandLine)

	flag.Parse()
	logger.Init(logOptions, func(config *zapcore.EncoderConfig) {
		config.TimeKey = "customer_time_key"
	})

	logger.G().Infow("main", "name", "xhli18")
}
