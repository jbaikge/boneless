#/bin/sh

docker run --rm --name dynamodb --publish 8000:8000 amazon/dynamodb-local:latest -jar DynamoDBLocal.jar -sharedDb
