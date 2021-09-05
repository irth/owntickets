package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/irth/owntickets/internal/owntickets"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"golang.org/x/crypto/bcrypt"
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
			{
				Name:  "hash-password",
				Usage: "Hash a password for the config",
				Action: func(c *cli.Context) error {
					fmt.Print("Password: ")
					scanner := bufio.NewScanner(os.Stdin)
					scanner.Scan()
					s := scanner.Text()
					if s == "" {
						logrus.Fatal("Password cannot be empty")
					}
					hashed, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
					if err != nil {
						logrus.WithError(err).Fatal("Couldn't hash the password")
					}
					fmt.Println(string(hashed))
					return err
				},
			},
		},
	}

	app.Run(os.Args)
}
