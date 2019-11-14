package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sync"

	"gitlab.com/scbrickley/xkcd"
)

var (
	wg        sync.WaitGroup
	numProcs  = runtime.NumCPU() * 10
	all       = flag.Bool("a", false, "Redownload all comics and skip duplicates?")
	randomize = flag.Bool("r", false, "Randomize order of comics?")
)

func init() {
	flag.Parse()
}

func main() {
	// Make the appropriate directories
	os.MkdirAll(xkcd.HomeDir, os.ModePerm)

	// Get a list of integers representing all the comics
	comicList := xkcd.ComicList()

	comicChan := make(chan int)

	go func() {
		for _, comic := range comicList {
			comicChan <- comic
		}
		close(comicChan)
	}()

	for i := 0; i < numProcs; i++ {
		wg.Add(1)
		go scraper(comicChan)
	}

	wg.Wait()
	fmt.Println("No more comics to download.")

	if *randomize {
		exec.Command("feh", "-z", "-x", xkcd.HomeDir).Run()
	} else {
		exec.Command("feh", "-n", "-x", xkcd.HomeDir).Run()
	}
}

func scraper(comics chan int) {
	defer wg.Done()
	for {
		number, ok := <-comics
		if !ok {
			return
		}

		comic, err := xkcd.NewComic(number)
		if err != nil {
			fmt.Println("Skipping comic #"+comic.ID(), "-", err)
			continue
		}

		// If comic.FileName() is already in $HOME/.xkcd, either:
		// 1. skip it, or
		// 2. Exit the program
		if comic.IsDuplicate() {
			if *all {
				fmt.Println("Skipping comic #"+comic.ID(), "- duplicate")
				continue
			}
			break
		}

		err = comic.Save()
		if err != nil {
			fmt.Println("Skipping comic #"+comic.ID(), "-", err)
			continue
		}
		fmt.Println("Downloading comic #" + comic.ID())
	}
}
