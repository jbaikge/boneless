package services

type Repository interface {
	ClassRepository
	DocumentRepository
	FileRepository
	FormRepository
	TemplateRepository
}
