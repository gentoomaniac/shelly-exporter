package main

import (
	"io"
	"os"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog/log"

	gocli "github.com/gentoomaniac/shelly-exporter/pkg/cli"
	"github.com/gentoomaniac/shelly-exporter/pkg/config"
	"github.com/gentoomaniac/shelly-exporter/pkg/exporter"
	"github.com/gentoomaniac/shelly-exporter/pkg/logging"
)

var (
	version = "unknown"
	commit  = "unknown"
	binName = "unknown"
	builtBy = "unknown"
	date    = "unknown"
)

var cli struct {
	logging.LoggingConfig

	ConfigFile *os.File `help:"config file path" required:""`

	Version gocli.VersionFlag `short:"V" help:"Display version."`
}

func main() {
	ctx := kong.Parse(&cli, kong.UsageOnError(), kong.Vars{
		"version": version,
		"commit":  commit,
		"binName": binName,
		"builtBy": builtBy,
		"date":    date,
	})
	logging.Setup(&cli.LoggingConfig)

	defer cli.ConfigFile.Close()
	b, err := io.ReadAll(cli.ConfigFile)
	if err != nil {
		log.Fatal().Msg("failed reading config")
	}

	c, err := config.NewConfigFromContent(b)
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}

	e := exporter.New(c)
	e.Run()

	ctx.Exit(0)
}
