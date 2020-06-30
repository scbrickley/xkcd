package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sync"

	"github.com/scbrickley/xkcd"
)

var (
	wg        sync.WaitGroup
	all       = flag.Bool("a", false, "Scan comic directory for missing comics")
	randomize = flag.Bool("r", false, "Randomize order of comics")
	hide      = flag.Bool("i", false, "Hide browser after scraper finishes")
	offline   = flag.Bool("o", false, "Run in offline mode")
	numProcs  = runtime.NumCPU()
)

func init() {
	flag.Parse()
	runtime.GOMAXPROCS(numProcs)
}

func main() {
	if !*hide {
		defer func() {
			if *randomize {
				exec.Command("feh", "-z", "-K", "captions", "-F", xkcd.HomeDir).Run()
			} else {
				exec.Command("feh", "-n", "-K", "captions", "-F", xkcd.HomeDir).Run()
			}
		}()
	}

	if *offline {
		return
	}

	// Make the appropriate directories
	os.MkdirAll(xkcd.CaptionDir, os.ModePerm)

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
