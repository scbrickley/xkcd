# XKCD Terminal viewer

A command line app that downloads all comics from [xkcd.com](https://xkcd.com) and allows you to browse through them from the terminal.

This is a rudimentary version. The following features are still in the works, and will be added as the project progresses:

- View comic title and captions
- Search by title
- Search by number
- Save favorites
- Installation via package manager

## Installation

*WARNING: THE `xkcd` executable has only been tested on Linux machines. However, the `xkcd` module should work on any OS, if you want to build your own.*

### 1. Install the Go programming language on your machine. Follow instructions [here](https://golang.org/doc/install?download=go1.12.9.linux-amd64.tar.gz).

### 1. Install `feh`

run `sudo apt-get install feh` for Debian-based distros

run `sudo pacman -S feh` for Arch-based distros

Or the equivalent for your distributions package manager.

### 1. Create your Go directory if you don't have one already.

`mkdir -p $HOME/go/src`

### 1. Navigate to that directory

`cd ~/go/src`

### 1. Fetch the repository using `go get gitlab.com/scbrickley/xkcd`

### 1. Navigate to the cmd directory, comile, and install

`cd $HOME/go/src/gitlab.com/scbrickley/xkcd/cmd/xkcd`
`go get`
`go build`
`go install`

**Note: these commands may not work if you keep your go working directory somewhere other than your home folder. Adjust these commands accordingly.**

## Usage

After installation, you should be able to type `xkcd` into your terminal to start the process. Initial download of all the comics may take a few minutes. Once it's done, a `feh` window should pop up and let you browse through the comics as you like.

The default behavior is to always pull up the newest comic. If you instead want to view the comics in a randomized order, type `xkcd -r` instead.

If you accidentally delete some of the comics, you can always run `xkcd -a` to re download them. The program will skip over any duplicate comics in the `~/.xkcd` directory, and only download the ones that are missing.

## Known issues

If `feh` is printing out error messages about incorrect sRGB profiles for .png files when viewing certain comics, follow these instructions:

### 1. Install `pngcrush` via `sudo apt-get install pngcrush` (or the equivalent for your package manager).

### 1. Navigate to the project directory: `cd ~/go/src/gitlab.com/scbrickely/xkcd`

### 1. Run the `fix` script: `./fix`

The error is actually harmless, but if the error messages are bothering you, this should fix the problem.
