package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"gopkg.in/urfave/cli.v1"
)

type appConfig struct {
	source  string
	target  string
	threads uint
}

func main() {
	app := cli.NewApp()
	app.Name = "glinka"

	config := appConfig{}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "source, s",
			Usage:       "load links structure from the `FILE`",
			Destination: &config.source,
		},
		cli.StringFlag{
			Name:        "target, t",
			Usage:       "save links structure to the `FILE`",
			Destination: &config.target,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "stats",
			Usage: "show links stats",
			Action: func(c *cli.Context) error {
				data, err := getData(c, &config)
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}

				str := reportStats(data)
				fmt.Print(str)
				return nil
			},
		},
	}

	app.Version = "0.2.0"
	app.Usage = "links analyzer"
	app.ArgsUsage = "domain"

	app.Action = func(c *cli.Context) error {
		_, err := getData(c, &config)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func getData(c *cli.Context, config *appConfig) (*LinksStore, error) {
	var data *LinksStore
	if config.source == "" {
		domain := c.Args().Get(0)
		if domain == "" {
			return nil, errors.New("Domain argument is mandatory. Try to use `glinka http://some.com`")
		}

		s := spiderMaster{domain: domain}
		data = s.run(domain)
	} else {
		data = NewLinksStore()
		if err := data.load(config.source); err != nil {
			return data, err
		}
	}

	if config.target != "" {
		if err := data.save(config.target); err != nil {
			return data, err
		}
	}

	return data, nil
}
