package main

import (
	"log"
	"os"

	"gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()
	app.Name = "glinka"

	app.Version = "0.1.0"
	app.Usage = "links checker"
	app.ArgsUsage = "domain"
	app.Action = func(c *cli.Context) error {
		domain := c.Args().Get(0)
		if domain == "" {
			return cli.NewExitError("Domain argument is mandatory. Try to use `glinka http://some.com`", 1)
		}
		s := spiderMaster{domain: domain}
		e := s.run(domain)

		if e != nil {
			return cli.NewExitError(e.Error(), 1)
		}

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
