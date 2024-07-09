# Scan Negative

CLI tool to scan film negatives with a Raspberry Pi via SSH.

## Usage

```
NAME:
   scan-negative - Scan your film from the comfort of your terminal

USAGE:
   scan-negative [global options] command [command options]

COMMANDS:
   preview, p  Preview what the scanner sees
   scan, s     Scan an image
   help, h     Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --host value  The host of the scanner either as IP or Bonjour hostname with the username (username@host). Can be assigned as FILM\_SCANNER\_HOST environment variable
   --port value  The port of the scanner used to preview the video feed. Can be assigned as FILM\_SCANNER\_PREVIEW\_PORT environment variable
   --help, -h    show help
```
