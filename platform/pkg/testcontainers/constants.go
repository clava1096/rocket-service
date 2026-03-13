package testcontainers

import "time"

const (
	// postgres container constants
	postgresPort           = "5435"
	postgresStartupTimeout = 1 * time.Minute

	// postgres environment variables
	postgresEnvUsernameKey = "order_user"
	postgresEnvPasswordKey = "order_password" //nolint:gosec

	// mongo container constants
	mongoPort          = "50051"
	MongoContainerName = "mongo"

	// mongo environment
	mongoImageName    = "MONGO_IMAGE"
	MongoImageNameKey = "MONGO_IMAGE_NAME"
	MongoHostKey      = "MONGO_HOST"
	MongoPortKey      = "MONGO_PORT"
	MongoDatabaseKey  = "MONGO_DATABASE"
	MongoUsernameKey  = "MONGO_INITDB_ROOT_USERNAME"
	MongoPasswordKey  = "MONGO_INITDB_ROOT_PASSWORD" //nolint:gosec
	MongoAuthDBKey    = "MONGO_AUTH_DB"
)
