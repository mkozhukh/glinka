package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()
	app.Name = "glinka"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "source, s",
			Usage: "load links structure from the `FILE`",
		},
		cli.StringFlag{
			Name:  "target, t",
			Usage: "save links structure to the `FILE`",
		},
		cli.StringFlag{
			Name:  "verbose",
			Usage: "enable verbose output",
		},
		cli.StringFlag{
			Name:  "quiet",
			Usage: "enable quiet output",
		},
		cli.StringFlag{
			Name:  "threads",
			Usage: "number of concurent requests",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "stats",
			Usage: "show links stats",
			Action: func(c *cli.Context) error {
				data, err := getData(c)
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}

				str := reportStats(data)
				fmt.Print(str)
				return nil
			},
		},
		{
			Name:  "errors",
			Usage: "show errors",
			Action: func(c *cli.Context) error {
				data, err := getData(c)
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}

				str := reportErrors(data)
				fmt.Print(str)
				return nil
			},
		},
	}

	app.Version = "0.2.0"
	app.Usage = "links analyzer"
	app.ArgsUsage = "domain"

	app.Action = func(c *cli.Context) error {
		_, err := getData(c)
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

func getData(c *cli.Context) (*LinksStore, error) {
	var data *LinksStore
	source := c.GlobalString("source")
	target := c.GlobalString("target")

	if source == "" {
		domain := c.Args().Get(0)
		if domain == "" {
			return nil, errors.New("Domain argument is mandatory. Try to use `glinka http://some.com`")
		}

		threads := c.GlobalInt("threads")
		if threads == 0 {
			threads = 3
		}
		verbose := c.GlobalBool("verbose")
		quiet := c.GlobalBool("quiet")
		s := spiderMaster{domain: domain, threads: threads, verbose: verbose, quiet: quiet}
		data = s.run(domain)
	} else {
		data = NewLinksStore()
		if err := data.load(source); err != nil {
			return data, err
		}
	}

	if target != "" {
		if err := data.save(target); err != nil {
			return data, err
		}
	}

	return data, nil
}
