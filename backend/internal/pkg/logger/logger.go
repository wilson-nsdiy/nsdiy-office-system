package logger

import (
	"context"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var globalLogger *zap.Logger

type Options struct {
	Level       string
	Format      string
	Output      OutputOptions
	Rotation    RotationOptions
	Caller      bool
	ServiceName string
}

type OutputOptions struct {
	ToStdout bool
	ToFile   bool
	FilePath string
}

type RotationOptions struct {
	MaxSizeMB  int
	MaxBackups int
	MaxAgeDays int
	Compress   bool
}

func Init(opts Options) error {
	level := parseLevel(opts.Level)

	encoderCfg := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	var enc zapcore.Encoder
	if opts.Format == "console" {
		enc = zapcore.NewConsoleEncoder(encoderCfg)
	} else {
		enc = zapcore.NewJSONEncoder(encoderCfg)
	}

	cores := make([]zapcore.Core, 0, 2)

	if opts.Output.ToStdout {
		infoPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= level && lvl < zapcore.WarnLevel
		})
		errPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= level && lvl >= zapcore.WarnLevel
		})
		cores = append(cores, zapcore.NewCore(enc, zapcore.Lock(os.Stdout), infoPriority))
		cores = append(cores, zapcore.NewCore(enc, zapcore.Lock(os.Stderr), errPriority))
	}

	if opts.Output.ToFile {
		filePath := opts.Output.FilePath
		if filePath == "" {
			filePath = "logs/oa-nsdiy.log"
		}

		dir := strings.TrimSuffix(filePath, "/")
		if idx := strings.LastIndex(filePath, "/"); idx > 0 {
			dir = filePath[:idx]
		}
		if idx := strings.LastIndex(filePath, "\\"); idx > 0 {
			dir = filePath[:idx]
		}
		if dir != "" {
			_ = os.MkdirAll(dir, 0o755)
		}

		lj := &lumberjack.Logger{
			Filename:   filePath,
			MaxSize:    opts.Rotation.MaxSizeMB,
			MaxBackups: opts.Rotation.MaxBackups,
			MaxAge:     opts.Rotation.MaxAgeDays,
			Compress:   opts.Rotation.Compress,
		}
		cores = append(cores, zapcore.NewCore(enc, zapcore.AddSync(lj), level))
	}

	if len(cores) == 0 {
		cores = append(cores, zapcore.NewCore(enc, zapcore.Lock(os.Stdout), level))
	}

	core := zapcore.NewTee(cores...)

	var zapOpts []zap.Option
	if opts.Caller {
		zapOpts = append(zapOpts, zap.AddCaller())
	}

	logger := zap.New(core, zapOpts...)

	if opts.ServiceName != "" {
		logger = logger.With(zap.String("service", opts.ServiceName))
	}

	globalLogger = logger
	return nil
}

func L() *zap.Logger {
	if globalLogger != nil {
		return globalLogger
	}
	return zap.NewNop()
}

func S() *zap.SugaredLogger {
	return L().Sugar()
}

func With(fields ...zap.Field) *zap.Logger {
	return L().With(fields...)
}

func Sync() {
	if globalLogger != nil {
		_ = globalLogger.Sync()
	}
}

func parseLevel(level string) zapcore.Level {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "debug":
		return zapcore.DebugLevel
	case "warn", "warning":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

// --- Context-aware logger ---

type loggerCtxKey struct{}

// IntoContext stores a request-scoped logger in the context.
func IntoContext(ctx context.Context, l *zap.Logger) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if l == nil {
		l = L()
	}
	return context.WithValue(ctx, loggerCtxKey{}, l)
}

// FromContext retrieves the request-scoped logger from context.
// Falls back to the global logger when no logger is stored.
func FromContext(ctx context.Context) *zap.Logger {
	if ctx == nil {
		return L()
	}
	if l, ok := ctx.Value(loggerCtxKey{}).(*zap.Logger); ok && l != nil {
		return l
	}
	return L()
}
