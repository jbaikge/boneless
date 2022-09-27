package dynamodb

const formPrefix = "form#"

func dynamoFormIds(id string) (pk string, sk string) {
	pk = formPrefix + id
	return
}
