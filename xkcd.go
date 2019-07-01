package main

import (
    "io"
    "io/ioutil"
    "net/http"
    "fmt"
    "os"
    "strconv"
    "strings"
    "flag"

    "github.com/anaskhan96/soup"
)

func stringInSlice(text string, list []string) bool {
    for _, val := range list {
        if val == text {
            return true
        }
    }
    return false
}

func getFileNames(path string) []string {
    contents, _ := ioutil.ReadDir(path)

    filenames := make([]string, 1)

    for _, file := range contents {
        filenames = append(filenames, file.Name())
    }

    return filenames
}

func main() {
    all := flag.Bool("a", false, "Download all comics and skip duplicates?")
    flag.Parse()

    // Make the appropriate directories
    dir := os.Getenv("HOME") + "/.xkcd/"

    os.MkdirAll(dir + "comics/captions/", os.ModePerm)
    os.MkdirAll(dir + "favorites/captions/", os.ModePerm)

    // So we can see what comics have already been downloaded
    filenames := getFileNames(dir + "comics/")

    // Download the HTML contents of the xkcd home page
    url := "https://xkcd.com/"
    resp, _ := soup.Get(url)
    doc := soup.HTMLParse(resp)

    // Find the prev link, cast it as an int, and ++ it to get the current comic #
    prev := doc.Find("a", "rel", "prev")
    comicNum, _ := strconv.Atoi(strings.ReplaceAll(prev.Attrs()["href"], "/", ""))
    comicNum++
    comID := strconv.Itoa(comicNum)

    for prev.Attrs()["href"] != "#" {
        filename := fmt.Sprintf("%04s", comID) + ".png"

        // If filename is already in .xkcd/comics, either:
        // 1. skip it, or
        // 2. Stop the downloader
        if stringInSlice(filename, filenames) {
            if *all {
                fmt.Println("Skipping comic #" + comID, "- duplicate")
                prev = doc.Find("a", "rel", "prev")
                comID = strings.ReplaceAll(prev.Attrs()["href"], "/", "")
                resp, _ = soup.Get(url + comID)
                doc = soup.HTMLParse(resp)
                continue
            } else {
                fmt.Println("No new comics.")
                break
            }
        }

        comicPath := dir + "comics/" + filename

        resp, _ = soup.Get(url + comID)
        doc = soup.HTMLParse(resp)

        // find the comic element
        comicElem := doc.Find("div", "id", "comic").Find("img")

        // if no comic element was found on this page, move on to the next page
        if comicElem.Pointer == nil {
            fmt.Println("Skipping comic #" + comID, "- no comic element")
            prev = doc.Find("a", "rel", "prev")
            comID = strings.ReplaceAll(prev.Attrs()["href"], "/", "")
            resp, _ = soup.Get(url + comID)
            doc = soup.HTMLParse(resp)
            continue
        }

        // find the url for the image
        comicSrc := comicElem.Attrs()["src"]

        // if the url is not formatted properly, skip it
        if strings.HasPrefix(comicSrc, "//imgs.xkcd.com") == false {
            fmt.Println("Skipping comic #" + comID, "- probably a flash game")
            prev = doc.Find("a", "rel", "prev")
            comID = strings.ReplaceAll(prev.Attrs()["href"], "/", "")
            resp, _ = soup.Get(url + comID)
            doc = soup.HTMLParse(resp)
            continue
        }

        // Add the 'https:' prefix
        comicSrc = "https:" + comicSrc

        // get the image data
        comic, err := http.Get(comicSrc)
        if err != nil {
            fmt.Println("Skipping comic #" + comID, "- bad image")
            prev = doc.Find("a", "rel", "prev")
            comID = strings.ReplaceAll(prev.Attrs()["href"], "/", "")
            resp, _ = soup.Get(url + comID)
            doc = soup.HTMLParse(resp)
            continue
        }

        comicFile, err := os.Create(comicPath)
        if err != nil {
            fmt.Println("Skipping comic #" + comID, "- could not create file")
            prev = doc.Find("a", "rel", "prev")
            comID = strings.ReplaceAll(prev.Attrs()["href"], "/", "")
            resp, _ = soup.Get(url + comID)
            doc = soup.HTMLParse(resp)
            continue
        }

        _, err = io.Copy(comicFile, comic.Body)
        if err != nil {
            fmt.Println("Skipping comic #" + comID, "- could not copy data")
            prev = doc.Find("a", "rel", "prev")
            comID = strings.ReplaceAll(prev.Attrs()["href"], "/", "")
            resp, _ = soup.Get(url + comID)
            doc = soup.HTMLParse(resp)
            continue
        }

        fmt.Println("Downloading", filename)
        prev = doc.Find("a", "rel", "prev")
        comID = strings.ReplaceAll(prev.Attrs()["href"], "/", "")
        resp, _ = soup.Get(url + comID)
        doc = soup.HTMLParse(resp)
        continue
    }
}
