package profx

import (
	"fmt"
	"log"
	"profx/domain"
)

// Crawler is exported.
type Crawler interface {
	Crawl() error
}

// Scraper is exported.
type Scraper interface {
	GetLinks(url string) ([]string, error)
}

// NewWebCrawler is exported.
func NewWebCrawler(sl Sourcer, s Scraper, r CrawlerRepository) WebCrawler {
	return WebCrawler{
		sourcer:    sl,
		scraper:    s,
		repository: r,
	}
}

// WebCrawler is exported.
type WebCrawler struct {
	sourcer    Sourcer
	scraper    Scraper
	repository CrawlerRepository
}

// Crawl is exported.
func (wc *WebCrawler) Crawl() error {
	counter := 0
	// Load the source links by the given source loader.
	sources, err := wc.sourcer.Load()
	if err != nil {
		// log.Printf("can't load source: %v", err)
		return fmt.Errorf("can't load source %w", err)
	}

	for _, source := range sources {
		urls, err := wc.scraper.GetLinks(source.URL)
		if err != nil {
			log.Printf("can't load get links for URL %s: %v\n", source.URL, err)
			continue
		}

		for _, url := range urls {
			hash := domain.Hash(url)
			hasLink, err := wc.repository.HasLink(hash)
			if err != nil {
				log.Printf("failed to check the link %s (%s) in db: %v\n", url, hash, err)
				continue
			}

			if !*hasLink {
				link := domain.Link{
					Hash:    hash,
					URL:     url,
					Source:  source.Name,
					FromURL: source.URL,
				}

				err = wc.repository.SaveLink(&link)
				if err != nil {
					log.Printf("failed to save link %v: %v\n", link, err)
					continue
				}

				isMatch, err := wc.sourcer.Match(source.Name, url)
				if err != nil {
					log.Printf("failed to check URL by rule regexp: %v\n", err)
					continue
				}
				if *isMatch {
					resource := domain.Resource{
						Hash:        hash,
						URL:         url,
						Kind:        source.Kind,
						Source:      source.Name,
						FromURL:     source.URL,
						Title:       "",
						Description: "",
					}
					err = wc.repository.SaveResource(&resource)
					if err != nil {
						log.Printf("failed to save resource %v: %v\n", link, err)
						continue
					}
					counter++
				}
			} else {
				// TODO: Temporary update the link with needed source info.
				dbLink, err := wc.repository.GetLink(hash)
				if err != nil {
					log.Printf("failed to get link %s (%s): %v", url, hash, err)
					continue
				}

				if dbLink != nil && (dbLink.Source == "" || dbLink.FromURL == "") {
					dbLink.Source = source.Name
					dbLink.FromURL = source.URL
					err = wc.repository.UpdateLink(dbLink)
					if err != nil {
						log.Printf("failed to update link %s (%s): %v\n", url, hash, err)
					}
				}
			}
		}
	}

	l := domain.SysLog{
		Message: fmt.Sprintf("crawler finished with counter %d", counter),
	}
	log.Println(l.Message)
	err = wc.repository.SaveLog(&l)
	if err != nil {
		// log.Printf("can't save log: %v", err)
		return fmt.Errorf("can't save log %w", err)
	}

	return nil
}
