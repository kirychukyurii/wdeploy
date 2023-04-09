package config

import (
	"fmt"
	"github.com/kirychukyurii/wdeploy/internal/pkg/file"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap/zapcore"
)

var Module = fx.Options(
	fx.Provide(NewConfig),
)

type Config struct {
	VarsFile  string
	HostsFile string
	LoggerConfig
	Variables
	Inventory
}

var varsConfigPath = "./config.yml"
var hostsConfigPath = "./config.yml"
var logLevel = "info"
var logFormat = "console"
var logDirectory = "./"

var defaultConfig = Config{
	VarsFile:  varsConfigPath,
	HostsFile: hostsConfigPath,
	LoggerConfig: LoggerConfig{
		LogLevel:     logLevel,
		LogFormat:    logFormat,
		LogDirectory: logDirectory,
	},
}

func NewConfig() Config {
	config := defaultConfig

	viper.SetConfigFile(varsConfigPath)
	if err := viper.ReadInConfig(); err != nil {
		panic(errors.Wrap(err, "failed to read config"))
	}

	if err := viper.Unmarshal(&config.Variables); err != nil {
		panic(errors.Wrap(err, "failed to marshal config"))
	}

	viper.SetConfigFile(hostsConfigPath)
	if err := viper.ReadInConfig(); err != nil {
		panic(errors.Wrap(err, "failed to read config"))
	}

	if err := viper.Unmarshal(&config.Inventory); err != nil {
		panic(errors.Wrap(err, "failed to marshal config"))
	}

	config.VarsFile = varsConfigPath
	config.HostsFile = hostsConfigPath
	config.LogLevel = logLevel
	config.LogFormat = logFormat
	config.LogDirectory = logDirectory

	return config
}

func SetVarConfigPath(path string) {
	if !file.IsFile(path) {
		panic(fmt.Sprintf("Path doesnt exists: %s", path))
	}

	fmt.Printf("loaded file: vars=%s\n", path)
	varsConfigPath = path
}

func SetHostsConfigPath(path string) {
	if !file.IsFile(path) {
		panic(fmt.Sprintf("Path doesnt exists: %s", path))
	}

	fmt.Printf("loaded file: hosts=%s\n", path)
	hostsConfigPath = path
}

func SetLoggerProperties(level string, format string, directory string) {
	_, err := zapcore.ParseLevel(level)
	if err != nil {
		level = logLevel
		fmt.Printf("%s. Setting to default: %s\n", err.Error(), level)
	}

	fmt.Printf("setting logger properties: log-level=%s, log-format=%s, log-directory=%s\n", level, format, directory)
	logLevel = level
	logFormat = format
	logDirectory = directory
}
