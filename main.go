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
					&cli.StringFlag{Name: "output", Aliases: []string{"o"}, Value: ".", Usage: "The destination folder. Defaults to '.'. Note, omit trailing slashes"},
					&cli.StringFlag{Name: "name", Aliases: []string{"n"}, Value: "image", Usage: "The base name of the file downloaded, excluding its file extension. Defaults to 'image', resulting in 'image.dng'"},
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

					log.Println("Scan saved at", scanFileLoc)
					log.Println("Downloading the scan...")

					downloadFileLoc := fmt.Sprintf("%s/%s.dng", ctx.String("output"), ctx.String("name"))
					downloadCmd := exec.Command("rsync", "-av", fmt.Sprintf("%s:%s", host, scanFileLoc), downloadFileLoc)
					downloadStdErr := strings.Builder{}
					downloadCmd.Stderr = &downloadStdErr
					if err := downloadCmd.Run(); err != nil {
						log.Println("Something went wrong trying to download the scan")
						log.Println(downloadStdErr.String())
						log.Fatal(err)
					}

					log.Println("Scan downloaded to", downloadFileLoc)

					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
