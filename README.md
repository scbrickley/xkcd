# XKCD Terminal Browser

A command line app that downloads all comics from [xkcd.com](https://xkcd.com) and allows you to browse through them from the terminal.

This is a rudimentary version. The following features are still in the works, and will be added as the project progresses:

- View comic title and captions
- Search by title
- Search by number
- Save favorites
- Installation via package manager
- Proper Documentation of the `xkcd` module

## Installation

*WARNING: THE `xkcd` executable has only been tested on Linux machines. However, the `xkcd` module should work on any OS. If you want to build your own executable simply list `"gitlab.com/scbrickley/xkcd"` under your imports in your `main.go` file.*

- Install the Go programming language on your machine, if you haven't already. Follow instructions [here](https://golang.org/dl/) to do so.
    - Don't forget to add `export PATH="$HOME/go/bin:$PATH"` to the end of `$HOME/.profile`
- Install `feh`
    - `sudo apt-get install feh` for Debian-based distros
    - `sudo pacman -S feh` for Arch-based distros
    - ...or the equivalent for your distribution's package manager.
- Create your Go directory if you don't have one already: `mkdir -p $HOME/go/src`
- Navigate to that directory `cd $HOME/go/src`
- Fetch the repository using `go get gitlab.com/scbrickley/xkcd`
- Navigate to the cmd directory, compile, and install
    - `cd $HOME/go/src/gitlab.com/scbrickley/xkcd/cmd/xkcd`
    - `go get`
    - `go build`
    - `go install`

**Note: these commands may not work if you keep your Go working directory somewhere other than your home folder. Adjust these commands accordingly.**

## Usage

After installation, you should be able to type `xkcd` into your terminal to start the process. Initial download of all the comics may take a few minutes. Once it's done, a `feh` window should pop up and let you browse through the comics as you like.

The default behavior is to always pull up the newest comic. If you instead want to view the comics in a randomized order, type `xkcd -r` instead.

If you accidentally delete some of the comics, you can always run `xkcd -a` to re-download them. The program will skip over any duplicate comics in the `$HOME/.xkcd` directory, and only download the ones that are missing.

## Known Issues

If `feh` is printing out error messages about incorrect sRGB profiles for .png files when viewing certain comics, follow these instructions:

- Install `pngcrush` via `sudo apt-get install pngcrush` (or the equivalent for your package manager).

- Navigate to the project directory: `cd ~/go/src/gitlab.com/scbrickely/xkcd`

- Run the `fix` script: `./fix`

The error is actually harmless, but if the error messages are bothering you, this should fix the problem.
