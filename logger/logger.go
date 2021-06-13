package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

func Logger() *zerolog.Logger {
	/* Create / Open the log file */
	logFile, err := os.OpenFile("./arc-dps-updater.log", os.O_APPEND|os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Error().Err(err).Msg("Error creating file")
	}

	/* Initalizing Loggers */
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stderr}
	fileWriter := zerolog.ConsoleWriter{Out: logFile, NoColor: true, TimeFormat: time.RFC1123}

	/* Initializing MultiLogger */
	multi := zerolog.MultiLevelWriter(consoleWriter, fileWriter)

	logger := zerolog.New(multi).With().Timestamp().Caller().Logger()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	return &logger
}
