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

const TimeFormat = "2006-01-02 15:04:05"

// NewZapLogger 根据日志配置创建Zap日志记录器
// 支持多种输出格式（JSON/控制台）和分级日志输出
// 返回配置好的日志记录器实例
func NewZapLogger(logCfg LogConfig) *zap.Logger {
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

	if logCfg.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	debugLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.DebugLevel
	})
	debugCore := zapcore.NewCore(encoder, toWriter(logCfg.Directory, "debug"), debugLevel)

	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.InfoLevel
	})
	infoCore := zapcore.NewCore(encoder, toWriter(logCfg.Directory, "info"), infoLevel)

	warnLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.WarnLevel
	})
	warnCore := zapcore.NewCore(encoder, toWriter(logCfg.Directory, "warn"), warnLevel)

	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	errorCore := zapcore.NewCore(encoder, toWriter(logCfg.Directory, "error"), errorLevel)

	stackLevel := zap.NewAtomicLevel()
	stackLevel.SetLevel(zap.ErrorLevel)

	options = append(options,
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
