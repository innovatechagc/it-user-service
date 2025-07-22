package integration

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/company/microservice-template/internal/handlers"
	"github.com/company/microservice-template/internal/services"
	testingPkg "github.com/company/microservice-template/internal/testing"
	"github.com/company/microservice-template/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type IntegrationTestSuite struct {
	suite.Suite
	containers *testingPkg.TestContainers
	router     *gin.Engine
}

func (suite *IntegrationTestSuite) SetupSuite() {
	ctx := context.Background()
	
	// Setup test containers
	containers, err := testingPkg.SetupTestContainers(ctx)
	suite.Require().NoError(err)
	suite.containers = containers

	// Setup Gin router
	gin.SetMode(gin.TestMode)
	suite.router = gin.New()
	
	// Setup services
	healthService := services.NewHealthService()
	logger := logger.NewLogger("debug")
	
	handlers.SetupRoutes(suite.router, healthService, logger)
}

func (suite *IntegrationTestSuite) TearDownSuite() {
	ctx := context.Background()
	if suite.containers != nil {
		suite.containers.Cleanup(ctx)
	}
}

func (suite *IntegrationTestSuite) TestHealthEndpoint() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/health", nil)
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "healthy")
}

func (suite *IntegrationTestSuite) TestReadinessEndpoint() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/ready", nil)
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "ready")
}

func (suite *IntegrationTestSuite) TestContainersAreRunning() {
	ctx := context.Background()
	
	// Test PostgreSQL
	pgConn, err := suite.containers.GetPostgresConnectionString(ctx)
	suite.NoError(err)
	suite.NotEmpty(pgConn)
	
	// Test Vault
	vaultAddr, err := suite.containers.GetVaultAddress(ctx)
	suite.NoError(err)
	suite.NotEmpty(vaultAddr)
	
	// Test Redis
	redisAddr, err := suite.containers.GetRedisAddress(ctx)
	suite.NoError(err)
	suite.NotEmpty(redisAddr)
}

func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}