package logger

import (
	"fmt"

	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Logger() *zerolog.Logger {
	/* Create / Open the log file */
	logFile, err := os.OpenFile("./arc-dps-updater.log", os.O_APPEND|os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Error().Err(err).Msg("Error creating file")
	}
	fmt.Printf("The log file is allocated at %s\n", logFile.Name())

	/* Initalizing Loggers */
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stderr}
	fileWriter := zerolog.New(logFile)

	/* Initializing MultiLogger */
	multi := zerolog.MultiLevelWriter(consoleWriter, fileWriter)

	logger := zerolog.New(multi).With().Timestamp().Logger()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	/*
		debug := flag.Bool("debug", false, "sets log level to debug")
		flag.Parse()
		if *debug {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		}
	*/

	return &logger
}
