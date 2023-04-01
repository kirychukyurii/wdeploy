package config

import (
	"fmt"
	"github.com/kirychukyurii/wdeploy/internal/pkg/file"
	"go.uber.org/fx"
	"go.uber.org/zap/zapcore"
)

var Module = fx.Options(
	fx.Provide(NewConfig),
)

type LoggerConfig struct {
	LogLevel     string
	LogFormat    string
	LogDirectory string
}

func NewConfig() LoggerConfig {
	props := new(LoggerConfig)

	props.LogLevel = "debug"
	props.LogFormat = "console"
	props.LogDirectory = "./"

	return *props
}

func SetConfigPath(path string) {
	if !file.IsFile(path) {
		panic(fmt.Sprintf("Path doesnt exists: %s", path))
		// logger.Zap.Panicf("Path doesnt exists: %s", path)
	}
}

func ParseLogLevel(logLevel string) {
	level, err := zapcore.ParseLevel(logLevel)
	if err != nil {
		logLevel = "info"
		fmt.Printf("%s: incorrect log level: %s. Setting to default: %s", err.Error(), level, logLevel)
		// logger.Zap.Warnf("Setting to default log-level: %s", logLevel)
	}
}

func SetLoggerProperties(logLevel string, logFormat string, logFile string) {
	props := new(LoggerConfig)

	props.LogLevel = logLevel
	props.LogFormat = logFormat
	props.LogDirectory = logFile
}
