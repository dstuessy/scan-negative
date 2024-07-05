package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/urfave/cli/v2"
)

const scanFileLoc = "/home/danielstuessy/.open-scanner/image.dng"
const scanCommand = "libcamera-still --raw on --nopreview -o %s"

// const host = "danielstuessy@filmscanner.local"

const hostEnvVar = "FILM_SCANNER_HOST"

func main() {

	app := &cli.App{
		Name:  "Terminal Scanner",
		Usage: "Scan your film from the comfort of your terminal",
		Commands: []*cli.Command{
			{
				Name:    "preview",
				Aliases: []string{"p"},
				Usage:   "Preview what the scanner sees",
				Action: func(cCtx *cli.Context) error {
					log.Println("Start preview...")
					return nil
				},
			},
			{
				Name:    "scan",
				Aliases: []string{"s"},
				Usage:   "Scan an image",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "host", Value: os.Getenv(hostEnvVar), Usage: fmt.Sprintf("The host of the scanner either as IP or Bonjour hostname with the username (username@host). Can be assigned as %s environment variable", hostEnvVar)},
				},
				Action: func(ctx *cli.Context) error {
					log.Println("Asking scanner to make a scan...")

					host := ctx.String("host")

					if host == "" {
						log.Fatal("No host provided. Please see --help for usage")
					}

					scanCmd := exec.Command("ssh", "-v", host, fmt.Sprintf(scanCommand, scanFileLoc))
					scanStdErr := strings.Builder{}
					scanCmd.Stderr = &scanStdErr
					if err := scanCmd.Run(); err != nil {
						log.Println("Something went wrong trying to make a scan")
						log.Println(scanStdErr.String())
						log.Fatal(err)
					}

					log.Println("Scan saved at /home/danielstuessy/.open-scanner/image.dng")
					log.Println("Downloading the scan...")

					downloadCmd := exec.Command("rsync", "-av", fmt.Sprintf("%s:%s", host, scanFileLoc), ".")
					downloadStdErr := strings.Builder{}
					downloadCmd.Stderr = &downloadStdErr
					if err := downloadCmd.Run(); err != nil {
						log.Println("Something went wrong trying to download the scan")
						log.Println(downloadStdErr.String())
						log.Fatal(err)
					}

					log.Println("Scan downloaded to ./image.dng")

					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
