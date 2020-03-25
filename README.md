# XKCD Terminal Browser
#### Now with ***Added Concurrency&#8482;***

A command line app that downloads all comics from [xkcd.com](https://xkcd.com) and allows you to browse through them from the terminal.

This is a rudimentary version. The following features are still in the works, and will be added as the project progresses:

- Custom image browser (to remove dependency on `feh`)
- ~~View captions~~
- View comic title
- Search by title
- Search by number
- Save favorites
- Installation via package manager and/or docker container
- Proper documentation of the `xkcd` module

## Installation

*WARNING: The `xkcd` executable has only been tested on Linux machines. However, the `xkcd` module should work on any OS. If you want to build your own executable, simply list `"github.com/scbrickley/xkcd"` under your imports in your own Go program.*

This project uses Go modules for dependency management, and it is recommended that you switch over to a more recent version of Go that uses modules as the default (`go 1.13` or later) before installing the program. The program should still work for older versions of Go that use the deprecated `GOPATH`, but it is not recommended.

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

## Testing the Executable

This program does not include any unit tests because: 

- Having unit tests that rely on an internet connection is bad practice

...but...

- Creating a mock server to test a program that will not run without an internet connection seems pointless, especially when the code is only designed to query one website. If the website ever changes enough to break the program, the unit tests will not reflect that fact.

Because of this, and because the program is less than 400 lines of code, and not difficult to debug, I've decided that the best way to test it is to use it.

You can help me test the program by installing it and running it for yourself, and submitting an issue if you run into any problems. You can also double check that there are no race conditions by compiling with `go build -race` when installing the program. If you get no output, then you're good to go.

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
