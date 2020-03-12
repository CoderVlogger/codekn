package profx

import (
	"fmt"
	"profx/domain"
	"profx/storage/crawler"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ory/dockertest"
	"github.com/stretchr/testify/suite"
)

type CrawlerRepositoryTestSuite struct {
	suite.Suite
	repository CrawlerRepository
	db         *crawler.DB
	resource   *dockertest.Resource
	pool       *dockertest.Pool
}

func TestCrawlerRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(CrawlerRepositoryTestSuite))
}

func (s *CrawlerRepositoryTestSuite) SetupSuite() {
	pool, err := dockertest.NewPool("")
	s.NoError(err)
	s.pool = pool
	s.pool.MaxWait = time.Second * 10

	s.NoError(err)

	resource, err := s.pool.RunWithOptions(&dockertest.RunOptions{Repository: "mysql",
		Tag: "5.7.25",
		Env: []string{"MYSQL_ROOT_PASSWORD=qweasdzxcv",
			"MYSQL_USER=profx",
			"MYSQL_PASSWORD=123456789",
			"MYSQL_DATABASE=profx",
			"TIMEZONE=UTC",
		},
	}, )
	s.NoError(err)
	s.resource = resource

	port := resource.GetPort("3306/tcp")

	config := crawler.Config{
		Username:  "profx",
		Password:  "123456789",
		ReadHost:  "localhost",
		ReadPort:  port,
		WriteHost: "localhost",
		WritePort: port,
		Schema:    "profx",
	}

	err = s.pool.Retry(func() error {
		s.db, err = config.New()

		if err != nil {
			return err
		}

		return s.db.Ping()
	})
	s.NoError(err)

	err = crawler.Migrate(config)
	if err != nil {
		panic(err)
	}

	s.repository = s.db
}

func (s *CrawlerRepositoryTestSuite) TearDownSuite() {
	err := s.pool.Purge(s.resource)
	s.NoError(err)
}

func (s *CrawlerRepositoryTestSuite) generateLinks(start, end int) {
	linkSQL := "INSERT INTO links (hash, url, source, from_url) VALUES (?, ?, ?, ?)"
	resourceSQL := "INSERT INTO resources (hash, url, kind, source, from_url) VALUES (?, ?, ?, ?, ?)"

	for i := start; i <= end; i++ {
		url := fmt.Sprintf("http://example.com/link%d", i)
		hash := domain.Hash(url)

		isArticle := false
		if i%2 == 0 {
			isArticle = true
		}

		source := "example"
		fromURL := fmt.Sprintf("http://example.com/from/link%d", i)

		res, err := s.db.Write.Exec(linkSQL, hash, url, source, fromURL)
		s.NotNil(res)
		s.NoError(err)

		if isArticle {
			res, err = s.db.Write.Exec(resourceSQL, hash, url, domain.ArticleResourceKind, source, fromURL)
			s.NotNil(res)
			s.NoError(err)
		}
	}
}

func (s *CrawlerRepositoryTestSuite) TestCrawlerRepository_HasLink() {
	s.generateLinks(1, 5)

	res, err := s.repository.HasLink("qweasdzxcv")
	s.NoError(err)
	s.NotNil(res)
	if res != nil {
		s.False(*res)
	}

	hash := "bd0904098d30fb9916bfb8a7cca263a2f8fe5f4a528ee31739d79fa264b2ca54"
	res, err = s.repository.HasLink(hash)
	s.NoError(err)
	s.NotNil(res)
	if res != nil {
		s.True(*res)
	}
}

func (s *CrawlerRepositoryTestSuite) TestCrawlerRepository_IsArticle() {
	s.generateLinks(6, 10)

	// url1 := "http://example.com/link1"
	hash1 := "bd0904098d30fb9916bfb8a7cca263a2f8fe5f4a528ee31739d79fa264b2ca54"
	res, err := s.repository.IsArticle(hash1)
	s.NoError(err)
	s.NotNil(res)
	if res != nil {
		s.False(*res)
	}

	// url2 := "http://example.com/link2"
	hash2 := "a811307d3e5ceaf5cea08cd9ced1fbb93e40a4d94dc894547b89b8476ea07667"
	res, err = s.repository.IsArticle(hash2)
	s.NoError(err)
	s.NotNil(res)
	if res != nil {
		s.True(*res)
	}
}

func (s *CrawlerRepositoryTestSuite) TestCrawlerRepository_SaveLink() {
	s.generateLinks(11, 15)

	url := "http://example.com/link1001"
	hash := "aef65d99727a3c679a1e06003414da8c4f6fef4cac93894706bba902897f7abb"
	err := s.repository.SaveLink(&domain.Link{
		Hash: hash,
		URL:  url,
	})
	s.NoError(err)

	res, err := s.repository.HasLink(hash)
	s.NoError(err)
	s.NotNil(res)
	if res != nil {
		s.True(*res)
	}
}

func (s *CrawlerRepositoryTestSuite) TestCrawlerRepository_UpdateLink() {
	s.generateLinks(16, 20)

	url := "http://example.com/link16"
	hash := domain.Hash(url)

	err := s.repository.UpdateLink(&domain.Link{
		Hash:    hash,
		URL:     url,
		Source:  "example",
		FromURL: "http://example.com/from/changed-link",
	})
	s.NoError(err)

	link, err := s.repository.GetLink(hash)
	s.NoError(err)
	s.Equal("http://example.com/from/changed-link", link.FromURL)
}
