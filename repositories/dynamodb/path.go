package dynamodb

const pathPrefix = "path#"

func dynamoPathIds(path string) (pk string, sk string) {
	pk = pathPrefix + path
	sk = path
	return
}
