/**
 * Run using:
 * 	dev help
 **/

package main

import (
	"./download"
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
			Usage: "Downloads source data from the web, storing in the local 'sources' directory",
			Action: func(c *cli.Context) error {
				download.DownloadAll()
				return nil
			},
		},
	}

	app.Run(os.Args)
}
