package main

import (
	"os"

	"github.com/irth/owntickets/internal/owntickets"
	"github.com/urfave/cli/v2"
)

func main() {
	ownTickets := owntickets.OwnTickets{}

	app := &cli.App{
		Name:  "owntickets",
		Usage: "personal ticketing system",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Load configuration from `FILE`",
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "serve",
				Usage: "Start the server",
				Action: func(c *cli.Context) error {
					configFile := c.String("config")
					if configFile != "" {
						ownTickets.Config.LoadFromFile(configFile)
					}
					ownTickets.Config.LoadFromEnv()
					return ownTickets.Run()
				},
			},
		},
	}

	app.Run(os.Args)
}
