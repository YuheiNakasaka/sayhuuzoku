package main

import (
	"log"
	"os"

	"github.com/YuheiNakasaka/sayhuuzoku/scraping"
	"github.com/YuheiNakasaka/sayhuuzoku/wakati"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "sayhuuzoku"
	app.Version = "0.0.1"

	app.Commands = []cli.Command{
		{
			Name:    "scraping",
			Aliases: []string{"s"},
			Usage:   "Fetch shop name from http://fujoho.jp/index.php?p=shop_list",
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "max-page, mp",
					Usage: "max page of scraping site",
				},
			},
			Action: func(c *cli.Context) error {
				scraping.Start(c.Int("max-page"))
				return nil
			},
		},
		{
			Name:    "wakati",
			Aliases: []string{"w"},
			Usage:   "Create wakati file from shoplist file",
			Action: func(c *cli.Context) error {
				wakati.Start()
				return nil
			},
		},
		{
			Name:    "generate",
			Aliases: []string{"g"},
			Usage:   "Generate shop name like huuzoku",
			Action: func(c *cli.Context) error {
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
