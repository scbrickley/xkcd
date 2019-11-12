package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"

	"gitlab.com/scbrickley/xkcd"
)

func main() {
	all := flag.Bool("a", false, "Redownload all comics and skip duplicates?")
	randomize := flag.Bool("r", false, "Randomize order of comics?")
	flag.Parse()

	// Make the appropriate directories
	os.MkdirAll(xkcd.HomeDir, os.ModePerm)

	comic := xkcd.LatestComic()

	for comic.Num >= 1 {
		// If comic.FileName() is already in .xkcd/comics, either:
		// 1. skip it, or
		// 2. Exit the program
		if comic.IsDuplicate() {
			if *all {
				fmt.Println("Skipping comic #"+comic.ID(), "- duplicate")
				comic.PrevComic()
				continue
			} else {
				fmt.Println("No more new comics.")
				break
			}
		}

		err := comic.Save()
		if err != nil {
			fmt.Println("Skipping comic #"+comic.ID(), "-", err)
			comic.PrevComic()
			continue
		}

		fmt.Println("Downloading comic #" + comic.ID())
		comic.PrevComic()
	}

	if *randomize {
		exec.Command("feh", "-z", "-x", xkcd.HomeDir).Run()
	} else {
		exec.Command("feh", "-n", "-x", xkcd.HomeDir).Run()
	}
}
