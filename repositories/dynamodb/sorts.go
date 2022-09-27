package dynamodb

import (
	"fmt"
	"time"
)

const sortPrefix = "sort#"

func dynamoSortIds(classId string, key string, docId string, value interface{}) (pk string, sk string) {
	pk = sortPrefix + classId + "#" + key
	if t, ok := value.(time.Time); ok {
		value = t.UTC().Format(time.RFC3339)
	}
	sk = fmt.Sprintf("%.*s#%s", 64, fmt.Sprintf("%v", value), docId)
	return
}
