package main

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog/log"

	gocli "github.com/gentoomaniac/shelly-exporter/pkg/cli"
	"github.com/gentoomaniac/shelly-exporter/pkg/logging"
	"github.com/gentoomaniac/shelly-exporter/pkg/shelly"
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

	p := shelly.NewPlugS(net.ParseIP("10.1.3.117"), "", "")
	err := p.Update()
	if err != nil {
		log.Error().Err(err).Msg("")
	}

	blob, err := json.Marshal(p.Status)
	fmt.Println(string(blob))

	ctx.Exit(0)
}
