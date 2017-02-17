/**
 * Run using:
 * 	dev help
 **/

package main

import (
	"./download"
	"./cleanse"
	"./flow"
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

				fmt.Println("\nGenerating flow models...")
				fmt.Println("------------------------------")
				flow.Generate()

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
			Usage: "Cleanses downloaded files, writing all as json to 'data/2-cleanse' directory",
			Action: func(c *cli.Context) error {
				cleanse.Cleanse()
				return nil
			},
		},

		{
			Name:  "flow",
			Usage: "Pulls together all the cleanse data into the final flow reference data. Writes to 'data/final' directory",
			Action: func(c *cli.Context) error {
				flow.Generate()
				return nil
			},
		},
	}

	app.Run(os.Args)
}
