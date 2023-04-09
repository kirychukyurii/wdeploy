package lib

import (
	"github.com/kirychukyurii/wdeploy/internal/config"
	"github.com/kirychukyurii/wdeploy/internal/constants"
	"github.com/kirychukyurii/wdeploy/internal/pkg/file"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"path/filepath"
	"strings"
	"time"
)

type Logger struct {
	Zap        *zap.SugaredLogger // by default
	DesugarZap *zap.Logger        // performance-sensitive code
}

func NewLogger(config config.Config) Logger {
	var options []zap.Option
	var encoder zapcore.Encoder
	/*
		config.LogFormat = "console"
		config.LogLevel = "debug"
	*/
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeTime:     localTimeEncoder,
	}

	if config.LogFormat == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	level := zap.NewAtomicLevelAt(toLevel(config.LogLevel))

	core := zapcore.NewCore(encoder, toWriter(config), level)

	stackLevel := zap.NewAtomicLevel()
	stackLevel.SetLevel(zap.WarnLevel)
	options = append(options,
		zap.AddCaller(),
		zap.AddStacktrace(stackLevel),
	)

	logger := zap.New(core, options...)
	return Logger{Zap: logger.Sugar(), DesugarZap: logger}
}

func localTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(constants.TimeFormat))
}

func toLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	case "dpanic":
		return zap.DPanicLevel
	case "panic":
		return zap.PanicLevel
	case "fatal":
		return zap.FatalLevel
	default:
		return zap.InfoLevel
	}
}

func toWriter(config config.Config) zapcore.WriteSyncer {
	fp := ""
	sp := string(filepath.Separator)

	fp, _ = filepath.Abs(filepath.Dir(filepath.Join(".")))
	fp += sp + "logs" + sp

	if config.LogDirectory != "" {
		if err := file.EnsureDirRW(config.LogDirectory); err != nil {
			fp = config.LogDirectory
		}
	}

	return zapcore.NewMultiWriteSyncer(
		//zapcore.AddSync(os.Stdout),

		// file rotation
		zapcore.AddSync(&lumberjack.Logger{
			Filename:   filepath.Join(fp, "tasker") + ".log",
			MaxSize:    100,
			MaxAge:     0,
			MaxBackups: 0,
			LocalTime:  true,
			Compress:   true,
		}),
	)
}
