package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
)

const scanFileLoc = "~/.open-scanner/image.dng"
const scanCommand = "libcamera-still --raw on --nopreview -o %s"
const previewCommand = "libcamera-vid -t 0 --width 1920 --height 1080 --codec h264 --inline --listen -o tcp://0.0.0.0:%s"

const hostEnvVar = "FILM_SCANNER_HOST"
const previewPortEnvVar = "FILM_SCANNER_PREVIEW_PORT"

const defaultPort = "8888"

func main() {
	app := &cli.App{
		Name:  "scan-negative",
		Usage: "Scan your film from the comfort of your terminal",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "host", Value: os.Getenv(hostEnvVar), Usage: fmt.Sprintf("The host of the scanner either as IP or hostname with the username (username@host). Can be assigned as %s environment variable", hostEnvVar)},
			&cli.StringFlag{Name: "preview-port", Value: os.Getenv(previewPortEnvVar), Usage: fmt.Sprintf("The port of the scanner used to preview the video feed. Can be assigned as %s environment variable", previewPortEnvVar)},
		},
		Commands: []*cli.Command{
			{
				Name:    "preview",
				Aliases: []string{"p"},
				Usage:   "Preview what the scanner sees in VLC Player",
				Action: func(ctx *cli.Context) error {
					log.Println("Asking scanner to preview its view...")

					host := ctx.String("host")

					if host == "" {
						log.Fatal("No host provided. Please see --help for usage")
					}

					port := ctx.String("preview-port")
					if port == "" {
						log.Println("No port provided. Using default port ", defaultPort)
						port = defaultPort
					}

					previewCmd := exec.Command("ssh", "-v", host, fmt.Sprintf(previewCommand, port))

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

					go func() {
						if err := previewCmd.Run(); err != nil {
							log.Fatal(err)
						}
					}()

					// Giving the preview some time to start
					time.Sleep(1 * time.Second)

					videoCmd := exec.Command("vlc", fmt.Sprintf("tcp/h264://%s:%s/", host, port))
					videoStdErr := strings.Builder{}
					videoCmd.Stderr = &videoStdErr
					if err := videoCmd.Run(); err != nil {
						log.Println("Something went wrong trying to play the video feed in VLC")
						log.Println(videoStdErr.String())
						log.Fatal(err)
					}

					return nil
				},
			},
			{
				Name:    "scan",
				Aliases: []string{"s"},
				Usage:   "Scan an image",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "output", Aliases: []string{"o"}, Value: ".", Usage: "The destination directory. Omit trailing slashes"},
					&cli.StringFlag{Name: "base-name", Aliases: []string{"n"}, Value: "image", Usage: "The base name of the file downloaded, excluding its file extension."},
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

					downloadFileLoc := fmt.Sprintf("%s/%s.dng", ctx.String("output"), ctx.String("base-name"))
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
