package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/netip"
	"net/url"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog/log"

	gocli "github.com/gentoomaniac/shelly-exporter/pkg/cli"
	"github.com/gentoomaniac/shelly-exporter/pkg/logging"
	"github.com/gentoomaniac/shelly-exporter/pkg/shelly"
	"github.com/gentoomaniac/shelly-exporter/pkg/shelly/auth"
	"github.com/gentoomaniac/shelly-exporter/pkg/shelly/request"
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

	Address netip.Addr `short:"a" help:"target address of the device" required:""`
	Rpc     string     `short:"r" help:"rpc method to call" required:""`

	User     string `short:"u" help:"authentication username" default:"admin"`
	Password string `short:"p" help:"authentication password"`

	Pretty      bool `help:"pretty print json response" default:"false"`
	ShowDevinfo bool `help:"show status endpoints data" default:"false"`

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

	devInfo, err := shelly.GetDeviceInfo(&cli.Address)
	if err != nil {
		log.Fatal().Err(err).Msg("failed getting device info")
	}
	if cli.ShowDevinfo {
		if cli.Pretty {
			prettyJSON, err := json.MarshalIndent(devInfo, "", "  ")
			if err != nil {
				log.Fatal().Err(err).Msg("failed marshalling devInfo")
			}
			fmt.Println(string(prettyJSON))
		} else {
			json, err := json.Marshal(devInfo)
			if err != nil {
				log.Fatal().Err(err).Msg("failed marshalling devInfo")
			}
			fmt.Println(string(json))
		}
	}

	var authInfo *auth.Auth
	if devInfo.Auth || devInfo.AuthEn {
		authInfo = &auth.Auth{User: cli.User, Password: cli.Password}
	}

	baseUrl, err := url.Parse("http://" + cli.Address.String())
	if err != nil {
		log.Fatal().Err(err).Msg("failed parsing base URL")
	}

	var resp []byte
	if devInfo.Gen == 0 {
		endpoint := baseUrl.JoinPath(cli.Rpc)

		resp, err = request.Request(endpoint, authInfo)
		if err != nil {
			log.Fatal().Err(err).Msg("failed sending request")
		}
	} else {
		endpoint := baseUrl.JoinPath("rpc", cli.Rpc)

		resp, err = request.DigestAuthedRequest(endpoint, authInfo, map[string]string{"id": "0"})
		if err != nil {
			log.Fatal().Err(err).Msg("failed sending request")
		}
	}

	if cli.Pretty {
		var prettyJSON bytes.Buffer
		err := json.Indent(&prettyJSON, resp, "", "  ")
		if err != nil {
			log.Fatal().Err(err).Msg("failed prettyfying json response")
		}
		fmt.Println(prettyJSON.String())
	} else {
		fmt.Println(string(resp))
	}

	ctx.Exit(0)
}
