package utils

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

func InitLogger() {
	logFile := &lumberjack.Logger{
		Filename:   "./logs/app.log",
		MaxSize:    5, // megabytes
		MaxBackups: 3, 
		MaxAge:     7, // days
		Compress:   true, // compress rotated files
	}

	log.Logger = zerolog.New(logFile).With().Timestamp().Logger()
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}


