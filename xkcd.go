package xkcd

import (
    //"fmt"
    "os/user"
    fp "path/filepath"

    //"github.com/anaskhan96/soup"
)

var HomeDir string

func init() {
    user, _ := user.Current()
    HomeDir = fp.Join(user.HomeDir, ".xkcd", "comics")
}

type Comic struct {
    Num int
    URL string
}

func (c Comic) Id() string {
    
}

func (c Comic) FetchLatest {
    
}
