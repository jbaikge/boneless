package gocms

import "time"

type dynamoDocument struct {
	DocumentId string
	ClassId    string
	ParentId   string
	TemplateId string
	VersionId  string
	Title      string
	Version    int
	Created    time.Time
	Updated    time.Time
}

func (dyn dynamoDocument) FromDocument(d *Document) {
	dyn.DocumentId = d.Id
	dyn.ClassId = d.ClassId
	dyn.ParentId = d.ParentId
	dyn.TemplateId = d.TemplateId
	dyn.Title = d.Title
	dyn.Version = d.Version
	dyn.Created = d.Created
	dyn.Updated = d.Updated
}

func (dyn dynamoDocument) ToDocument() (d Document) {
	d.Id = dyn.DocumentId
	d.ClassId = dyn.ClassId
	d.ParentId = dyn.ParentId
	d.TemplateId = dyn.TemplateId
	d.Title = dyn.Title
	d.Version = dyn.Version
	d.Created = dyn.Created
	d.Updated = dyn.Updated
	return
}
