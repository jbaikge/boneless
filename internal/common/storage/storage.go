package storage

import (
	"os"
	"strings"
)

type Option string

const (
	DynamoDB = Option("dynamodb")
	SQLite   = Option("sqlite")
	Memory   = Option("memory")
)

func FromEnv() Option {
	switch strings.ToLower(os.Getenv("BONELESS_STORE")) {
	case string(DynamoDB):
		return DynamoDB
	case string(SQLite):
		return SQLite
	case string(Memory):
		return Memory
	default:
		return Memory
	}
}
