package main

import (
	"github.com/spacemeshos/poet-core-api"
	"github.com/urfave/cli"
	"log"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "pccli"
	app.Usage = "control plane for poet core"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "rpcserver",
			Value: poet_core_api.DefaultRPCHostPort,
			Usage: "host:port of poet core service",
		},
	}
	app.Commands = []cli.Command{
		computeCommand,
		cleanCommand,
		getNIPCommand,
		getProofCommand,
		verifyCommand,
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
