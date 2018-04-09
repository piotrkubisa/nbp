package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/karolgorecki/nbp/svc"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "nbp-cli"
	app.Usage = "NBP currencies"
	app.Commands = []cli.Command{
		{
			Name:    "avg",
			Aliases: []string{"a"},
			Usage:   `Get's average currency rates for all or selection of currencies`,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "date",
					Value: time.Now().Format("2006-01-02"),
					Usage: "",
				},
				cli.StringFlag{
					Name:  "code",
					Value: "*",
					Usage: "",
				},
			},
			Action: func(c *cli.Context) error {
				q, err := svc.Average(c.String("date"), c.String("code"))
				if err != nil {
					return err
				}

				return formatOutput(q)
			},
		},
		{
			Name:    "both",
			Aliases: []string{"c"},
			Usage:   `Get's buy, sell values for all or selection of currencies`,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "date",
					Value: time.Now().Format("2006-01-02"),
					Usage: "",
				},
				cli.StringFlag{
					Name:  "code",
					Value: "*",
					Usage: "",
				},
			},
			Action: func(c *cli.Context) error {
				q, err := svc.Both(c.String("date"), c.String("code"))
				if err != nil {
					return err
				}

				return formatOutput(q)
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func formatOutput(q svc.Query) error {
	b, err := json.MarshalIndent(q, "", "  ")
	if err != nil {
		return err
	}

	fmt.Printf("%s", b)
	return nil
}
