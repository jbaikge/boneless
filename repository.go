package boneless

type Repository interface {
	ClassRepository
	DocumentRepository
	FileRepository
	TemplateRepository
}
