package dynamodb

import (
	"testing"

	"github.com/jbaikge/boneless"
	"github.com/zeebo/assert"
)

func TestNewDynamoClass(t *testing.T) {
	class := boneless.Class{
		Id:     "from_class",
		Name:   t.Name(),
		Fields: []boneless.Field{{Name: "field_1"}},
	}

	dc := newDynamoClass(&class)

	assert.Equal(t, classPrefix+class.Id, dc.PK)
	assert.Equal(t, "class", dc.SK)
	assert.Equal(t, t.Name(), dc.Name)
	assert.DeepEqual(t, class.Fields, dc.Data)
}

func TestToClass(t *testing.T) {
	id := "to_class"
	dc := dynamoClass{
		PK:   classPrefix + id,
		SK:   "class",
		Name: t.Name(),
		Data: []boneless.Field{{Name: "field_1"}},
	}
	class := dc.ToClass()

	assert.Equal(t, id, class.Id)
	assert.Equal(t, t.Name(), class.Name)
	assert.DeepEqual(t, dc.Data, class.Fields)
}
