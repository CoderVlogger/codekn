package app

import "time"

type Article struct {
	Title   string
	URL     string
	Created time.Time
	Source  string
}

type IndexPage struct {
	Title    string
	Date     string
	Articles []Article
}
