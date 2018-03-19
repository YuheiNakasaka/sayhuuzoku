package main

import (
	"fmt"
	"log"
	"os"

	"github.com/YuheiNakasaka/sayhuuzoku/db"
	"github.com/YuheiNakasaka/sayhuuzoku/generator"
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
			Name:    "init",
			Aliases: []string{"i"},
			Usage:   "Init database",
			Action: func(c *cli.Context) error {
				mydb := db.MyDB{}
				return mydb.New()
			},
		},
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
				return scraping.Start(c.Int("max-page"))
			},
		},
		{
			Name:    "wakati",
			Aliases: []string{"w"},
			Usage:   "Create wakati data from shoplist file",
			Action: func(c *cli.Context) error {
				return wakati.Start()
			},
		},
		{
			Name:    "generate",
			Aliases: []string{"g"},
			Usage:   "Generate shop name like huuzoku",
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "count, c",
					Usage: "word count",
				},
			},
			Action: func(c *cli.Context) error {
				shopName, err := generator.Start(c.Int("count"))
				if err != nil {
					return err
				}
				fmt.Println(shopName)
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
