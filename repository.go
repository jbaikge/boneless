package gocms

type RepositoryStats struct {
	Fetches int
	Inserts int
	Updates int
	Deletes int
}

type Repository interface {
	Stats() RepositoryStats
	ClassRepository
	DocumentRepository
	TemplateRepository
}
