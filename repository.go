package gocms

type Repository interface {
	ClassRepository
	DocumentRepository
	TemplateRepository
}
