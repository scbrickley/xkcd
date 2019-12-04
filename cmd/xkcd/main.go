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
	all       = flag.Bool("a", false, "Redownload all comics and skip duplicates?")
	randomize = flag.Bool("r", false, "Randomize order of comics?")
	hide      = flag.Bool("h", false, "Don't load comic browser after comic scraper finishes?")
	offline   = flag.Bool("o", false, "Run in offline mode?")
	numProcs  = runtime.NumCPU()
)

func init() {
	flag.Parse()
}

func main() {
	if !*hide {
		defer func() {
			if *randomize {
				exec.Command("feh", "-z", "-x", xkcd.HomeDir).Run()
			} else {
				exec.Command("feh", "-n", "-x", xkcd.HomeDir).Run()
			}
		}()
	}

	if *offline {
		return
	}

	runtime.GOMAXPROCS(numProcs)

	// Make the appropriate directories
	os.MkdirAll(xkcd.HomeDir, os.ModePerm)

	// Get a list of integers representing all the comics
	comicList, err := xkcd.ComicList()
	if err != nil {
		fmt.Println(err)
		fmt.Println("Exiting program")
		os.Exit(1)
	}

	comicChan := make(chan int, len(comicList))

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
			fmt.Println(err)
			fmt.Println("Exiting program")
			os.Exit(1)
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
