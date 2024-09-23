package logger

import (
	"flag"

	"github.com/rs/zerolog"
)

// Setup initializes the logging configuration using the zerolog package.
// It sets the log level based on a command-line flag.
func Setup() {
	// Set the time field format to Unix time.
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// Define a command-line flag named "debug" to set the log level to debug.
	debug := flag.Bool("debug", false, "sets log level to debug")

	// Parse the command-line flags.
	flag.Parse()

	// Set the default global log level to Info.
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// If the debug flag is set, change the global log level to Debug.
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}
