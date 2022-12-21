A command line app that downloads all comics from [xkcd.com](https://xkcd.com) and allows you to browse through them from the terminal.

_Note: This project is no longer being updated or maintained. This was an early project of mine that helped me learn Go and concurrent programming. I'm leaving it up as an example of my early work for my code portfolio. The `xkcd` executable still works, but I no longer intend to add any new features._

## Installation the xkcd binary

- Install the Go programming language on your machine, if you haven't already. Follow instructions [here](https://golang.org/doc/install) to do so.
	- Note that there are other methods for installing Go besides downloading and extracting the tarball (the method outlined on the official Go website). I won't detail them here, but as long as `go version` outputs something like `go version 1.14.1 linux/amd64`, you're fine.
- Install `feh`
	- `sudo apt install feh` for Debian-based distros
	- `sudo pacman -S feh` for Arch-based distros
	- ...or the equivalent for your distribution's package manager.
- run `git clone github.com/scbrickley/xkcd` from `$HOME/go/src`, if you're still using `GOPATH`, or from whatever directory you want if have switched over to Go Modules.
- `cd xkcd/cmd/xkcd`
- `go build && go install`

## Usage

After installation, you should be able to type `xkcd` into your terminal to start the process. Initial download of all the comics may take a few minutes. Once it's done, a `feh` window should pop up and let you browse through the comics as you like.

### Flags
| Flag | Description | Default Behavior w/o Flag |
|------|-------------|---------|
| -o | Run in offline mode. | Exit program if no internet connection can be made. |
| -i | Exit program after comic scraper is done. | Load comic browser once scraper is finished. |
| -a | Scan the comic directory and download missing comics. Skip any duplicate comics. | Stop scraper once the first duplicate comic is reached. |
| -r | Randomize order of comics. | View the newest comic first. Right arrow key cycles to the next most recent comic. |

### Browser Controls

| Key | Behavior |
|-----|----------|
| Right Arrow/Left Arrow | Next/Previous Comic |
| Up Arrow/Down Arrow | Zoom In/Out |
| Ctrl+Up/Down/Left/Right | Adjust View (For comics that are zoomed in or too large too fit on the screen) |
| Q | Exit Program |
