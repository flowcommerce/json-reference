/**
 * Run using:
 * 	dev help
 **/

package main

import (
	"./download"
	"./cleansed"
	"github.com/codegangsta/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "reference"
	app.Usage = "Flow Reference Library"

	app.Commands = []cli.Command{
		{
			Name:  "download",
			Usage: "Downloads source data from the web, storing in the local 'data/1-sources' directory",
			Action: func(c *cli.Context) error {
				download.DownloadAll()
				return nil
			},
		},
		{
			Name:  "cleanse",
			Usage: "Cleanses downloaded files, writing all as json to 'data/1-cleansed' directory",
			Action: func(c *cli.Context) error {
				cleansed.Cleanse()
				return nil
			},
		},
	}

	app.Run(os.Args)
}
