package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"gitlab.com/scbrickley/xkcd"
)

func main() {
	all := flag.Bool("a", false, "Download all comics and skip duplicates?")
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

		// if no comic element was found on this page, move on to the next page
		if comic.ImgElem().Pointer == nil {
			fmt.Println("Skipping comic #"+comic.ID(), "- no comic element")
			comic.PrevComic()
			continue
		}

		// if the url is not formatted properly, skip it
		if strings.Contains(comic.ImgSrc(), "imgs.xkcd.com") == false {
			fmt.Println("Skipping comic #"+comic.ID(), "- probably a flash game")
			comic.PrevComic()
			continue
		}

		// get the image data
		comicData, err := comic.Image()
		if err != nil {
			fmt.Println("Skipping comic #"+comic.ID(), "- bad image")
			comic.PrevComic()
			continue
		}

		comicFile, err := os.Create(comic.FilePath())
		if err != nil {
			fmt.Println("Skipping comic #"+comic.ID(), "- could not create file")
			comic.PrevComic()
			continue
		}

		_, err = io.Copy(comicFile, comicData.Body)
		if err != nil {
			fmt.Println("Skipping comic #"+comic.ID(), "- could not copy data")
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
