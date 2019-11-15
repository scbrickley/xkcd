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
	speed     = flag.Int("s", 1, "Multiply the number of spawned scrapers by this number")
	all       = flag.Bool("a", false, "Redownload all comics and skip duplicates?")
	randomize = flag.Bool("r", false, "Randomize order of comics?")
	hide      = flag.Bool("h", false, "Don't load comic browser after comic scraper finishes?")
	debug     = flag.Bool("d", false, "Print debug info?")
)

func init() {
	flag.Parse()
}

func main() {
	numProcs := runtime.NumCPU() * *speed
	runtime.GOMAXPROCS(numProcs)

	if *debug {
		fmt.Printf("Using %d routines\n", numProcs)
	}

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

	if *hide {
		return
	}

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

		comic := xkcd.NewComic(number)

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

		err := comic.Save()
		if err != nil {
			fmt.Println("Skipping comic #"+comic.ID(), "-", err)
			continue
		}
		fmt.Println("Downloading comic #" + comic.ID())
	}
}
