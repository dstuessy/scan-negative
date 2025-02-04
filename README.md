# Scan Negative

Scan your film from the comfort of your terminal.

## Contents

- [Prerequisites](#prerequisites)
- [Usage](#usage)

## Prerequisites

- Run this command on a UNIX based machine i.e. macOS, Linux, or BSD.
- A Raspberry Pi with a camera looking at some film negatives.
- SSH access to your Raspberry Pi via public key pair.
- Libcamera support on the raspberry pi. The legacy API won't work.
- A folder called `.open-scanner` in your Raspberry Pi's home directory. This will depend on which user you use.
- [SSH](https://en.wikipedia.org/wiki/Secure_Shell) installed on your client machine[^1].
- [Rsync](https://en.wikipedia.org/wiki/Rsync) installed on your client machine[^1].
- [VLC](https://en.wikipedia.org/wiki/VLC_media_player) installed on your client machine[^1].

[^1]: i.e. the machine you're using to run this command.

## Usage

```
NAME:
   scan-negative - Scan your film from the comfort of your terminal

USAGE:
   scan-negative [global options] command [command options] 

COMMANDS:
   preview, p  Preview what the scanner sees in VLC Player
   scan, s     Scan an image
   help, h     Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --host value          The host of the scanner either as IP or hostname with the username (username@host). Can be assigned as FILM_SCANNER_HOST environment variable
   --preview-port value  The port of the scanner used to preview the video feed. Can be assigned as FILM_SCANNER_PREVIEW_PORT environment variable
   --help, -h            show help
```

### Scan Command

Scan a film negative

```
NAME:
   scan-negative scan - Scan an image

USAGE:
   scan-negative scan [command options]

OPTIONS:
   --output value, -o value     The destination directory. Omit trailing slashes (default: ".")
   --base-name value, -n value  The base name of the file downloaded, excluding its file extension. (default: "image")
   --help, -h                   show help
```

### Preview Command

Preview what is viewable from the scanner's camera. Useful to frame and find focus.

```
NAME:
   scan-negative preview - Preview what the scanner sees in VLC Player

USAGE:
   scan-negative preview [command options]

OPTIONS:
   --help, -h  show help
```
