# Scan Negative

CLI tool to scan film negatives with a Raspberry Pi via SSH.

## Contents

- [Usage](#usage)
- [Prerequisites](#prerequisites)

## Prerequisites

- A Raspberry Pi with a camera looking at some film negatives.
- SSH access to your Raspberry Pi via public key pair.
- Libcamera support on the raspberry pi. The legacy API won't work.

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

### Scan Command

Scan a film negative

```
NAME:
   scan-negative scan - Scan an image

USAGE:
   scan-negative scan [command options]

OPTIONS:
   --output value, -o value  The destination folder. Defaults to '.'. Note, omit trailing slashes (default: ".")
   --name value, -n value    The base name of the file downloaded, excluding its file extension. Defaults to 'image', resulting in 'image.dng' (default: "image")
   --help, -h                show help
```

### Preview Command

Preview what is viewable from the scanner's camera. Useful to frame and find focus.

```
NAME:
   scan-negative preview - Preview what the scanner sees

USAGE:
   scan-negative preview [command options]

OPTIONS:
   --help, -h  show help
```
