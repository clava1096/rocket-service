package integration

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/clava1096/rocket-service/platform/pkg/logger"
	"github.com/joho/godotenv"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
)

const testsTimeout = 5 * time.Minute

var (
	env *TestEnvironment

	suiteCtx    context.Context
	suiteCancel context.CancelFunc
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Inventory Service Integration Test Suite")
}

var _ = BeforeSuite(func() {
	err := logger.Init(loggerLevelValue, true)
	if err != nil {
		panic(fmt.Sprintf("failed to set up a logger: %v", err))
	}

	suiteCtx, suiteCancel = context.WithTimeout(context.Background(), testsTimeout)

	envVars, err := godotenv.Read(filepath.Join("..", "..", "..", "deploy", "compose", "inventory", ".env"))
	if err != nil {
		logger.Fatal(suiteCtx, "Failed to load .env file", zap.Error(err))
	}

	for key, value := range envVars {
		_ = os.Setenv(key, value)
	}

	logger.Info(suiteCtx, "Launching the test environment...")
	env = setupTestEnvironment(suiteCtx)
})

var _ = AfterSuite(func() {
	logger.Info(context.Background(), "Completion of the test suite")
	if env != nil {
		teardownTestEnvironment(suiteCtx, env)
	}
	suiteCancel()
})
