package logging

import (
	"io"
	"os"

	"log"

	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
)

type LoggingConfig struct {
	Verbosity int  `short:"v" help:"Increase verbosity." type:"counter"`
	Quiet     bool `short:"q" help:"Do not run upgrades." default:"false" negatable:""`
	Json      bool `help:"Log as json" default:"true" negatable:""`
	Debug     bool `help:"shortcut for -vvvv" default:"false" negatable:""`
}

func Setup(config *LoggingConfig) {
	if !config.Quiet {
		if config.Debug {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		} else {
			zerolog.SetGlobalLevel(zerolog.Level(int(zerolog.ErrorLevel) - config.Verbosity))
		}

		if !config.Json {
			zlog.Logger = zlog.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		}
	} else {
		zerolog.SetGlobalLevel(zerolog.Disabled)
		log.SetFlags(0)
		log.SetOutput(io.Discard)
	}
}
