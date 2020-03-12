package profx

import "profx/domain"

type CrawlerRepository interface {
	HasLink(hash string) (*bool, error)
	IsArticle(hash string) (*bool, error)

	GetLink(hash string) (*domain.Link, error)
	SaveLink(link *domain.Link) error
	UpdateLink(link *domain.Link) error

	SaveResource(link *domain.Resource) error
	SaveLog(log *domain.SysLog) error

	LoadSources() ([]domain.Source, error)
	LoadRules() ([]domain.Rule, error)
}
