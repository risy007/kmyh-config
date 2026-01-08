package config

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type LogConfig struct {
	Level       string `mapstructure:"Level"`
	Format      string `mapstructure:"Format"`
	ToFile      bool   `mapstructure:"ToFile"`
	Directory   string `mapstructure:"Directory"`
	Development bool   `mapstructure:"Development"`
}

const TimeFormat = "2006-01-02 15:04:05"

func NewZapLogger(conf *AppConfig) *zap.Logger {
	var options []zap.Option
	var encoder zapcore.Encoder

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeTime:     localTimeEncoder,
	}

	if conf.Logger.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	//level := zap.NewAtomicLevelAt(toLevel(conf.Log.Level))
	//core := zapcore.NewCore(encoder, toWriter(conf), level)

	debugLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.DebugLevel
	})
	debugCore := zapcore.NewCore(encoder, toWriter(conf.Logger.Directory, "debug"), debugLevel)

	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.InfoLevel
	})
	infoCore := zapcore.NewCore(encoder, toWriter(conf.Logger.Directory, "info"), infoLevel)

	warnLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.WarnLevel
	})
	warnCore := zapcore.NewCore(encoder, toWriter(conf.Logger.Directory, "warn"), warnLevel)

	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	errorCore := zapcore.NewCore(encoder, toWriter(conf.Logger.Directory, "error"), errorLevel)

	stackLevel := zap.NewAtomicLevel()
	stackLevel.SetLevel(zap.ErrorLevel)

	options = append(options,
		//zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(stackLevel),
	)

	logger := zap.New(zapcore.NewTee(debugCore, infoCore, warnCore, errorCore), options...)

	return logger
}

func localTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(TimeFormat))
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

func toWriter(dir string, level string) zapcore.WriteSyncer {
	fp := ""
	sp := string(filepath.Separator)
	fp, _ = filepath.Abs(filepath.Dir(filepath.Join(".")))
	if dir != "" {
		fp += sp + dir + sp
	}
	switch level {
	case "debug":
		return zapcore.AddSync(&lumberjack.Logger{ // 文件切割
			Filename:   filepath.Join(fp, level) + ".log",
			MaxSize:    100,
			MaxAge:     7,
			MaxBackups: 14,
			LocalTime:  true,
			Compress:   true,
		})
	default:
		return zapcore.NewMultiWriteSyncer(
			zapcore.AddSync(os.Stdout),
			zapcore.AddSync(&lumberjack.Logger{ // 文件切割
				Filename:   filepath.Join(fp, level) + ".log",
				MaxSize:    100,
				MaxAge:     7,
				MaxBackups: 14,
				LocalTime:  true,
				Compress:   true,
			}),
		)
	}
}
