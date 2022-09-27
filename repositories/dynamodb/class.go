package dynamodb

import (
	"time"

	"github.com/jbaikge/boneless"
)

const classPrefix = "class#"

func dynamoClassIds(id string) (pk string, sk string) {
	pk = classPrefix + id
	return
}

type dynamoClass struct {
	PK       string
	SK       string
	ParentId string
	Name     string
	Created  time.Time
	Updated  time.Time
	Data     []boneless.Field
}

func newDynamoClass(c *boneless.Class) (dyn *dynamoClass) {
	pk, sk := dynamoClassIds(c.Id)
	dyn = &dynamoClass{
		PK:       pk,
		SK:       sk,
		ParentId: c.ParentId,
		Name:     c.Name,
		Created:  c.Created,
		Updated:  c.Updated,
		Data:     make([]boneless.Field, len(c.Fields)),
	}
	copy(dyn.Data, c.Fields)
	return
}

func (dyn *dynamoClass) ToClass() (c boneless.Class) {
	c = boneless.Class{
		Id:       dyn.PK[len(classPrefix):],
		ParentId: dyn.ParentId,
		Name:     dyn.Name,
		Created:  dyn.Created,
		Updated:  dyn.Updated,
		Fields:   make([]boneless.Field, len(dyn.Data)),
	}
	copy(c.Fields, dyn.Data)
	return
}
