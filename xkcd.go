package xkcd

import (
	"fmt"
	"net/http"
	"os/user"
	fp "path/filepath"
	"strconv"
	"strings"

	"github.com/anaskhan96/soup"
)

var HomeDir string

func init() {
	user, _ := user.Current()
	HomeDir = fp.Join(user.HomeDir, ".xkcd")
}

type Comic struct {
	Num  int
	HTML soup.Root
}

// ID returns the Comic Number in string form
func (c Comic) ID() string {
	return strconv.Itoa(c.Num)
}

// Returns the URL to the comic
func (c Comic) URL() string {
	return "https://xkcd.com/" + c.ID()
}

// Returns the text in the "href" of the "prev" link of the comic webpage
func (c Comic) PrevText() string {
	return c.HTML.Find("a", "rel", "prev").Attrs()["href"]
}

// Finds the previous comic, sets the Num field to the number of the new comic
// and updates the HTML content accordingly
func (c *Comic) PrevComic() {
	c.Num, _ = strconv.Atoi(strings.ReplaceAll(c.PrevText(), "/", ""))
	resp, _ := soup.Get(c.URL())
	c.HTML = soup.HTMLParse(resp)
}

// Returns the parsed HTML content of the img element containing the actual comic
func (c Comic) ImgElem() soup.Root {
	return c.HTML.Find("div", "id", "comic").Find("img")
}

// Returns the URL for the img element that contains the actual comic
func (c Comic) ImgSrc() string {
	return "https:" + c.ImgElem().Attrs()["src"]
}

// Returns the image data for the comic
func (c Comic) Image() (*http.Response, error) {
	comicData, err := http.Get(c.ImgSrc())
	return comicData, err
}

// Returns a consistently formatted filename for the comic
func (c Comic) FileName() string {
	return fmt.Sprintf("%04s", c.ID()) + ".png"
}

// Returns a string represnting the appropriate filepath for a comic
func (c Comic) FilePath() string {
	return fp.Join(HomeDir, c.FileName())
}

// Returns a Comic representing the most recent post on xkcd.com
func LatestComic() Comic {
	url := "https://xkcd.com"
	resp, _ := soup.Get(url)
	doc := soup.HTMLParse(resp)

	prev := doc.Find("a", "rel", "prev")
	comicNum, _ := strconv.Atoi(strings.ReplaceAll(prev.Attrs()["href"], "/", ""))
	comicNum++

	return Comic{comicNum, doc}
}
