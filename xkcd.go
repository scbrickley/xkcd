package xkcd

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	fp "path/filepath"
	"strconv"
	"strings"

	"github.com/anaskhan96/soup"
)

var (
	HomeDir    string
	errNoComic = errors.New("No comic element")
)

func init() {
	user, _ := user.Current()
	HomeDir = fp.Join(user.HomeDir, ".xkcd")
}

type Comic struct {
	num  int
	html soup.Root
}

func (c Comic) Num() int {
	return c.num
}

// ID returns the Comic Number in string form
func (c Comic) ID() string {
	return strconv.Itoa(c.Num())
}

// Returns the URL to the comic
func (c Comic) URL() string {
	return "https://xkcd.com/" + c.ID()
}

// Returns the text in the "href" of the "prev" link of the comic webpage
func (c Comic) PrevText() string {
	return c.html.Find("a", "rel", "prev").Attrs()["href"]
}

// Finds the previous comic, sets the Num field to the number of the new comic
// and updates the HTML content accordingly
func (c *Comic) PrevComic() {
	c.num, _ = strconv.Atoi(strings.ReplaceAll(c.PrevText(), "/", ""))
	resp, _ := soup.Get(c.URL())
	c.html = soup.HTMLParse(resp)
}

// Returns the parsed HTML content of the img element containing the actual comic
func (c Comic) ImgElem() soup.Root {
	elem := c.html.Find("div", "id", "comic")
	if elem.Error != nil {
		return soup.Root{nil, "", errNoComic}
	}

	return elem.Find("img")
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

// Checks if a comic with the same filename is already in $HOME/.xkcd/
func (c Comic) IsDuplicate() bool {
	filenames := getFileNames(HomeDir)
	return stringInSlice(c.FileName(), filenames)
}

// Returns a list of integers, starting with the newest comic's number
// and continuing in decreasing order, because we always want to get i
// the newest comics first.
func ComicList() []int {
	max := LatestComic().Num()

	list := makeRange(max)
	list = reverse(list)

	return list
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

// Gets a comic based on the number
func NewComic(comicNum int) (Comic, error) {
	url := "https://xkcd.com/" + strconv.Itoa(comicNum)
	resp, err := soup.Get(url)
	doc := soup.HTMLParse(resp)

	return Comic{comicNum, doc}, err
}

// Save the comic to $HOME/.xkcd
func (c Comic) Save() error {
	// No comic element was found
	if c.ImgElem().Error != nil {
		return errNoComic
	}

	// Improperly formatted image URL, which usually means it's not an image
	if strings.Contains(c.ImgSrc(), "imgs.xkcd.com") == false {
		return errors.New("Probably a flash game")
	}

	// Get the image data
	comicData, err := c.Image()

	// Check for possibly corrupted image data
	if err != nil {
		return errors.New("Bad image")
	}

	// Create the file where the image data will be written
	comicFile, err := os.Create(c.FilePath())
	if err != nil {
		return errors.New("Could not create file")
	}

	// Copy the image data into the filespace
	_, err = io.Copy(comicFile, comicData.Body)
	if err != nil {
		return errors.New("Could not copy data")
	}

	return nil
}

//--------------- Helper Functions --------------//
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

func makeRange(max int) []int {
	var numbers []int
	for i := 1; i <= max; i++ {
		numbers = append(numbers, i)
	}

	return numbers
}

func remove(s []int, num int) []int {
	var index int
	for i, x := range s {
		if x == num {
			index = i
			break
		}
	}

	return append(s[:index], s[index+1:]...)
}

func reverse(list []int) []int {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}

	return list
}
