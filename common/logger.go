package common

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// / TODO : take a prefix to know where the logs come
// / from. E..g "[URCHIN]" and "[URCHIN-ADMIN]"
func SetupLogger(file string) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Info().Msg("Logger created")
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}

	errorFile, err := os.Create(file)
	if err != nil {
		log.Fatal().Err(err).Msg("No se pudo crear el archivo de errores")
	}

	filteredErrorWriter := zerolog.FilteredLevelWriter{
		Writer: zerolog.MultiLevelWriter(consoleWriter, errorFile), // El destino final es el archivo
		Level:  zerolog.ErrorLevel,                                 // El nivel mínimo que dejará pasar
	}

	logger := zerolog.New(&filteredErrorWriter).With().Timestamp().Logger().Level(zerolog.InfoLevel)
	logger.Info().Msgf("Logger created")
}
