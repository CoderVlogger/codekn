package profx

import (
	"log"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/pkg/errors"
)

// CollyScraper is exported.
type CollyScraper struct {
}

// NewCollyScraper is exported.
func NewCollyScraper() *CollyScraper {
	return &CollyScraper{}
}

// GetLinks is exported.
func (p *CollyScraper) GetLinks(url string) ([]string, error) {
	var urls []string
	collector := colly.NewCollector()
	collector.SetRequestTimeout(3 * time.Second)

	collector.OnHTML("a", func(e *colly.HTMLElement) {
		url := e.Attr("href")
		urls = append(urls, strings.TrimSpace(url))
	})

	collector.OnRequest(func(r *colly.Request) {
		log.Println("visiting", r.URL)
	})

	err := collector.Visit(url)
	if err != nil {
		return nil, errors.Wrap(err, "error on page visit")
	}

	return urls, nil
}
