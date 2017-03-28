/**
 * Run using:
 * 	dev help
 **/

package main

import (
	"./cleanse"
	"./download"
	"./final"
	"./javascript"
	"fmt"
	"github.com/codegangsta/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "reference"
	app.Usage = "Flow Reference Library"

	app.Commands = []cli.Command{
		{
			Name:  "all",
			Usage: "Runs all scripts",
			Action: func(c *cli.Context) error {
				fmt.Println("Downloading data...")
				fmt.Println("------------------------------")
				download.DownloadAll()

				fmt.Println("\nCleansing data...")
				fmt.Println("------------------------------")
				cleanse.Cleanse()

				fmt.Println("\nGenerating final models...")
				fmt.Println("------------------------------")
				final.Generate()

				fmt.Println("\nGenerating javascript models...")
				fmt.Println("------------------------------")
				javascript.Generate()

				fmt.Println("\nDone\n")
				return nil
			},
		},

		{
			Name:  "download",
			Usage: "Downloads source data from the web, storing in the local 'data/source' directory",
			Action: func(c *cli.Context) error {
				download.DownloadAll()
				return nil
			},
		},

		{
			Name:  "cleanse",
			Usage: "Cleanses downloaded files, writing all as json to 'data/cleanse' directory",
			Action: func(c *cli.Context) error {
				cleanse.Cleanse()
				return nil
			},
		},

		{
			Name:  "final",
			Usage: "Pulls together all the cleanse data into the final final reference data. Writes to 'data/final' directory",
			Action: func(c *cli.Context) error {
				final.Generate()
				return nil
			},
		},

		{
			Name:  "javascript",
			Usage: "Generates data used by our javascript libraries. Writes to 'data/javascript' directory",
			Action: func(c *cli.Context) error {
				javascript.Generate()
				return nil
			},
		},
	}

	app.Run(os.Args)
}
