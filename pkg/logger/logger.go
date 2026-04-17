package logger

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.Logger

	sugar *zap.SugaredLogger
	level zapcore.Level
}

func NewLogger(environment, logLevel string) *Logger {
	var zapCfg zap.Config

	if environment == "production" {
		zapCfg = zap.NewProductionConfig()

		zapCfg.EncoderConfig.TimeKey = "timestamp"
		zapCfg.EncoderConfig.MessageKey = "message"
		zapCfg.EncoderConfig.LevelKey = "severity"
		zapCfg.EncoderConfig.CallerKey = "caller"
		zapCfg.EncoderConfig.StacktraceKey = "stacktrace"
		zapCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		zapCfg.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	} else {
		zapCfg = zap.NewDevelopmentConfig()

		zapCfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		zapCfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("15:04:05.000")
		zapCfg.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	}

	level := getZapLevel(logLevel)
	zapCfg.Level = zap.NewAtomicLevelAt(level)

	zapCfg.OutputPaths = []string{"stdout"}
	zapCfg.ErrorOutputPaths = []string{"stdout"}

	logger, err := zapCfg.Build(
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
		zap.AddCallerSkip(1),
	)
	if err != nil {
		panic(fmt.Errorf("failed to build logger: %w", err))
	}

	return &Logger{
		Logger: logger,
		sugar:  logger.Sugar(),
		level:  level,
	}
}

func NewNop() *Logger {
	return &Logger{
		Logger: zap.NewNop(),
		sugar:  zap.NewNop().Sugar(),
	}
}

func (l *Logger) WithContext(ctx context.Context) *Logger {
	var (
		fields      []zap.Field
		sugarFields []any
	)

	if requestID := ctx.Value("request_id"); requestID != nil {
		if str, ok := requestID.(string); ok {
			fields = append(fields, zap.String("request_id", str))
			sugarFields = append(sugarFields, zap.String("request_id", str))
		}
	}

	if userID := ctx.Value("user_id"); userID != nil {
		if str, ok := userID.(string); ok {
			fields = append(fields, zap.String("user_id", str))
			sugarFields = append(sugarFields, zap.String("user_id", str))
		}
	}

	if shortCode := ctx.Value("short_code"); shortCode != nil {
		if str, ok := shortCode.(string); ok {
			fields = append(fields, zap.String("short_code", str))
			sugarFields = append(sugarFields, zap.String("short_code", str))
		}
	}

	if len(fields) > 0 {
		return &Logger{
			Logger: l.With(fields...),
			sugar:  l.sugar.With(sugarFields...),
			level:  l.level,
		}
	}

	return l
}

func (l *Logger) WithFields(fields map[string]any) *Logger {
	zapFields := make([]zap.Field, 0, len(fields))

	zapSugarFields := make([]any, 0, len(fields))

	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
		zapSugarFields = append(zapSugarFields, zap.Any(k, v))
	}

	return &Logger{
		Logger: l.With(zapFields...),
		sugar:  l.sugar.With(zapSugarFields...),
		level:  l.level,
	}
}

// Debug логирует debug сообщение.
func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.Logger.Debug(msg, fields...)
}

// Debugw логирует debug сообщение с ключ-значение полями (sugar).
func (l *Logger) Debugw(msg string, keysAndValues ...any) {
	l.sugar.Debugw(msg, keysAndValues...)
}

// Info логирует info сообщение.
func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.Logger.Info(msg, fields...)
}

// Infow логирует info сообщение с ключ-значение полями (sugar).
func (l *Logger) Infow(msg string, keysAndValues ...any) {
	l.sugar.Infow(msg, keysAndValues...)
}

// Warn логирует warning сообщение.
func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.Logger.Warn(msg, fields...)
}

// Warnw логирует warning сообщение с ключ-значение полями (sugar).
func (l *Logger) Warnw(msg string, keysAndValues ...any) {
	l.sugar.Warnw(msg, keysAndValues...)
}

// Errorw логирует ошибку с ключ-значение полями (sugar).
func (l *Logger) Errorw(msg string, keysAndValues ...any) {
	l.sugar.Errorw(msg, keysAndValues...)
}

// Fatalw логирует фатальную ошибку с ключ-значение полями (sugar) и завершает программу.
func (l *Logger) Fatalw(msg string, keysAndValues ...any) {
	l.sugar.Fatalw(msg, keysAndValues...)
}

// Error логирует ошибку.
func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.Logger.Error(msg, fields...)
}

// Fatal логирует фатальную ошибку и завершает программу.
func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.Logger.Fatal(msg, fields...)
}

// Sync сбрасывает буферы логгера.
func (l *Logger) Sync() {
	err := l.Logger.Sync()
	if err != nil {
		panic(err)
	}
}

// getZapLevel преобразует строковый уровень в zapcore. Level.
func getZapLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

// RequestLog структура для логирования HTTP запросов.
type RequestLog struct {
	Method     string
	Path       string
	StatusCode int
	Duration   time.Duration
	ClientIP   string
	UserAgent  string
	RequestID  string
}

// LogHTTPRequest логирует HTTP запрос с правильными уровнями.
func (l *Logger) LogHTTPRequest(req RequestLog) {
	fields := []zap.Field{
		zap.String("method", req.Method),
		zap.String("path", req.Path),
		zap.Int("status_code", req.StatusCode),
		zap.Duration("duration", req.Duration),
		zap.Int64("duration_ms", req.Duration.Milliseconds()),
		zap.String("client_ip", req.ClientIP),
		zap.String("user_agent", req.UserAgent),
	}

	if req.RequestID != "" {
		fields = append(fields, zap.String("request_id", req.RequestID))
	}

	switch {
	case req.StatusCode >= 500:
		l.Error("HTTP request failed", fields...)
	case req.StatusCode >= 400:
		l.Warn("HTTP request failed", fields...)
	default:
		l.Info("HTTP request succeeded", fields...)
	}
}
