package dynamodb

import "fmt"

const templatePrefix = "template#"

func dynamoTemplateIds(id string, version int) (pk string, sk string) {
	pk = templatePrefix + id
	sk = templatePrefix + fmt.Sprintf("v%06d", version)
	return
}
