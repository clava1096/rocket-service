package integration

import (
	"context"
	"os"
	"time"

	"github.com/clava1096/rocket-service/platform/pkg/logger"
	"github.com/clava1096/rocket-service/platform/pkg/testcontainers"
	"github.com/clava1096/rocket-service/platform/pkg/testcontainers/app"
	"github.com/clava1096/rocket-service/platform/pkg/testcontainers/mongo"
	"github.com/clava1096/rocket-service/platform/pkg/testcontainers/network"
	"github.com/clava1096/rocket-service/platform/pkg/testcontainers/path"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"
)

const (
	inventoryAppName    = "inventory-app"
	inventoryDockerFile = "deploy/docker/inventory/Dockerfile"

	grpcPortKey = "GRPC_PORT"

	loggerLevelValue = "debug"
	startupTimeOut   = 3 * time.Minute
)

type TestEnvironment struct {
	Network *network.Network
	Mongo   *mongo.Container
	App     *app.Container
}

func setupTestEnvironment(ctx context.Context) *TestEnvironment {
	logger.Info(ctx, "Setting up test environment...")
	generatedNetwork, err := network.NewNetwork(ctx, projectName)

	if err != nil {
		logger.Fatal(ctx, "Error creating network")
	}

	logger.Info(ctx, "Network created...")

	mongoUsername := getEnvWithLogginig(ctx, testcontainers.MongoUsernameKey)
	mongoPassword := getEnvWithLogginig(ctx, testcontainers.MongoPasswordKey)
	mongoImageName := getEnvWithLogginig(ctx, testcontainers.MongoImageNameKey)
	mongoDatabase := getEnvWithLogginig(ctx, testcontainers.MongoDatabaseKey)

	grpcPort := getEnvWithLogginig(ctx, grpcPortKey)

	generatedMongo, err := mongo.NewContainer(
		ctx,
		mongo.WithNetworkName(generatedNetwork.Name()),
		mongo.WithContainerName(testcontainers.MongoContainerName),
		mongo.WithImageName(mongoImageName),
		mongo.WithAuth(mongoUsername, mongoPassword),
		mongo.WithDatabase(mongoDatabase),
		mongo.WithLogger(logger.Logger()))

	if err != nil {
		cleanupTestEnvironment(ctx, &TestEnvironment{Network: generatedNetwork})
		logger.Fatal(ctx, "Error creating mongo container")
	}

	logger.Info(ctx, "Mongo container created!")

	projectRoot := path.GetProjectRoot()

	appEnv := map[string]string{
		// todo Переопределяем хост MongoDB для подключения к контейнеру из testcontainers
		testcontainers.MongoHostKey: generatedMongo.Config().ContainerName,
	}

	waitStrategy := wait.ForListeningPort(nat.Port(grpcPort + "/tcp")).WithStartupTimeout(startupTimeOut)

	appContainer, err := app.NewContainer(ctx,
		app.WithEnv(appEnv),
		app.WithName(inventoryAppName),
		app.WithPort(grpcPort),
		app.WithDockerfile(projectRoot, inventoryDockerFile),
		app.WithLogOutput(os.Stdout),
		app.WithStartupWait(waitStrategy),
		app.WithLogger(logger.Logger()),
		app.WithNetwork(generatedNetwork.Name()))

	if err != nil {
		cleanupTestEnvironment(ctx, &TestEnvironment{Network: generatedNetwork, Mongo: generatedMongo})
		logger.Fatal(ctx, "failed to run app container", zap.Error(err))
	}

	logger.Info(ctx, "App container created!")
	logger.Info(ctx, "Test environment created!")

	return &TestEnvironment{
		Network: generatedNetwork,
		Mongo:   generatedMongo,
		App:     appContainer,
	}

}

func getEnvWithLogginig(ctx context.Context, name string) string {
	value := os.Getenv(name)
	if value == "" {
		logger.Warn(ctx, "Переменная окружения не установлена", zap.String("key", name))
	}
	return value
}
