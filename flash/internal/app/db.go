package app

type ResourcesRepository interface {
	LoadResources() ([]Resource, error)
}
