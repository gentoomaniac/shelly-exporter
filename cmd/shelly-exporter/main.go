package main

import (
	"github.com/alecthomas/kong"
	"github.com/rs/zerolog/log"

	gocli "github.com/gentoomaniac/shelly-exporter/pkg/cli"
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

	Foo struct{} `cmd:"" help:"FooBar command"`
	Run struct{} `cmd:"" help:"Run the application (default)." default:"1" hidden:""`

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

	switch ctx.Command() {
	case "foo":
		log.Info().Msg("foo command")
	default:
		exporter.DoSomething()
	}
	ctx.Exit(0)
}
