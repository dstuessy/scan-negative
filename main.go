package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/urfave/cli/v2"
)

const scanFileLoc = "~/.open-scanner/image.dng"
const scanCommand = "libcamera-still --raw on --nopreview -o %s"
const previewCommand = "libcamera-vid -t 0 --width 1920 --height 1080 --codec h264 --inline --listen -o tcp://0.0.0.0:8888"

const hostEnvVar = "FILM_SCANNER_HOST"

func main() {

	app := &cli.App{
		Name:  "Terminal Scanner",
		Usage: "Scan your film from the comfort of your terminal",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "host", Value: os.Getenv(hostEnvVar), Usage: fmt.Sprintf("The host of the scanner either as IP or Bonjour hostname with the username (username@host). Can be assigned as %s environment variable", hostEnvVar)},
		},
		Commands: []*cli.Command{
			{
				Name:    "preview",
				Aliases: []string{"p"},
				Usage:   "Preview what the scanner sees",
				Action: func(ctx *cli.Context) error {
					log.Println("Asking scanner to preview its view...")

					host := ctx.String("host")

					if host == "" {
						log.Fatal("No host provided. Please see --help for usage")
					}

					previewCmd := exec.Command("ssh", "-v", host, previewCommand)

					previewStdOut, err := previewCmd.StdoutPipe()
					if err != nil {
						log.Fatal(err)
					}
					go io.Copy(os.Stdout, previewStdOut)

					previewStdErr, err := previewCmd.StderrPipe()
					if err != nil {
						log.Fatal(err)
					}
					go io.Copy(os.Stderr, previewStdErr)

					previewStdIn, err := previewCmd.StdinPipe()
					if err != nil {
						log.Fatal(err)
					}
					go io.Copy(previewStdIn, os.Stdin)

					if err := previewCmd.Run(); err != nil {
						log.Fatal(err)
					}

					log.Println("Preview available at tcp://", host, ":8888")
					fmt.Println("Ctrl+C to stop the preview...")

					return nil
				},
			},
			{
				Name:    "scan",
				Aliases: []string{"s"},
				Usage:   "Scan an image",
				Flags: []cli.Flag{
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

					log.Println("Scan saved at", fmt.Sprintf("%s:%s", host, scanFileLoc))
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
