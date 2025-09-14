package util

import (
	"log/slog"

	"gopkg.in/natefinch/lumberjack.v2"
)

func InitializeLog() {
	r := &lumberjack.Logger{
		Filename:   "./foo.log",
		LocalTime:  true,
		MaxSize:    1,
		MaxAge:     3,
		MaxBackups: 5,
		Compress:   true,
	}
	logger := slog.New(slog.NewJSONHandler(r, nil))

	slog.SetDefault(logger)
}
