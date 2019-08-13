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
var ComicDir string
var CapDir string
var FavDir string
var FavCapDir string

func init() {
	user, _ := user.Current()
	HomeDir = fp.Join(user.HomeDir, ".xkcd")
	ComicDir = fp.Join(HomeDir, "comics")
	CapDir = fp.Join(ComicDir, "captions")
	FavDir = fp.Join(HomeDir, "favorites")
	FavCapDir = fp.Join(FavDir, "captions")
}

type Comic struct {
	Num  int
	HTML soup.Root
}

type HTTPImage struct {
    resp *http.Response
    url string
}

func (i HTTPImage) URL() string {
    return i.url
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
func (c Comic) Image() (*http.Response.Body, error) {
	comicData, err := http.Get(c.ImgSrc())
	return comicData, err
}

// Returns a consistently formatted filename for the comic
func (c Comic) FileName() string {
	return fmt.Sprintf("%04s", c.ID()) + ".png"
}

func (c Comic) FilePath() string {
	return fp.Join(ComicDir, c.FileName())
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
