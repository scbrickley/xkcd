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

	comic, err := xkcd.LatestComic()
	if err != nil {
		fmt.Println(err)
		fmt.Println("Exiting program")
		os.Exit(1)
	}

	for comic.Num() >= 1 {
		// If comic.FileName() is already in $HOME/.xkcd, either:
		// 1. skip it, or
		// 2. Exit the program
		if comic.IsDuplicate() {
			if *all {
				fmt.Println("Skipping comic #"+comic.ID(), "- duplicate")
				err := comic.PrevComic()
				if err != nil {
					fmt.Println(err)
					fmt.Println("Exiting program")
					os.Exit(1)
				}
				continue
			} else {
				fmt.Println("No more new comics.")
				break
			}
		}

		err = comic.Save()
		if err != nil {
			fmt.Println("Skipping comic #"+comic.ID(), "-", err)
			err = comic.PrevComic()
			if err != nil {
				fmt.Println(err)
				fmt.Println("Exiting program")
				os.Exit(1)
			}
			continue
		}

		fmt.Println("Downloading comic #" + comic.ID())
		err = comic.PrevComic()
		if err != nil {
			fmt.Println(err)
			fmt.Println("Exiting program")
			os.Exit(1)
		}
	}

	if *randomize {
		exec.Command("feh", "-z", "-x", xkcd.HomeDir).Run()
	} else {
		exec.Command("feh", "-n", "-x", xkcd.HomeDir).Run()
	}
}
