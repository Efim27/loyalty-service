package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/suite"

	"loyalty-service/internal/handlers"
)

type AuthTestingSuite struct {
	suite.Suite
	serverCtxCancel context.CancelFunc
	server          handlers.Server
}

func (suite *AuthTestingSuite) SetupSuite() {
	ctx, cancel := context.WithCancel(context.Background())
	suite.serverCtxCancel = cancel
	suite.server = handlers.NewServer()
	suite.server.Prepare(ctx)
}

func (suite *AuthTestingSuite) TearDownSuite() {
	suite.serverCtxCancel()
	err := suite.server.App.Shutdown()
	if err != nil {
		suite.T().Log(err)
	}
}

func TestAuthSuite(t *testing.T) {
	suite.Run(t, new(AuthTestingSuite))
}

func (suite *AuthTestingSuite) TestLoginEmpty() {
	req := httptest.NewRequest("POST", "/api/user/login", nil)
	resp, _ := suite.server.App.Test(req, -1)

	suite.Equal(fiber.StatusBadRequest, resp.StatusCode, "Login empty test")
}

func (suite *AuthTestingSuite) TestRegisterEmpty() {
	req := httptest.NewRequest("POST", "/api/user/register", nil)
	resp, _ := suite.server.App.Test(req, -1)

	suite.Equal(fiber.StatusBadRequest, resp.StatusCode, "Register empty test")
}

func (suite *AuthTestingSuite) TestRegisterAndLogin() {
	userInfo := handlers.AuthInput{
		Login:    randomString(5),
		Password: randomString(10),
	}

	userInfoBytes, err := json.Marshal(&userInfo)
	if err != nil {
		suite.Fail("marshalling userInfo error", err)
	}

	req := httptest.NewRequest("POST", "/api/user/register", bytes.NewBuffer(userInfoBytes))
	req.Header.Set("Content-Type", fiber.MIMEApplicationJSON)
	resp, _ := suite.server.App.Test(req, -1)
	suite.Equal(fiber.StatusOK, resp.StatusCode, "Register new user")

	req = httptest.NewRequest("POST", "/api/user/register", bytes.NewBuffer(userInfoBytes))
	req.Header.Set("Content-Type", fiber.MIMEApplicationJSON)
	resp, _ = suite.server.App.Test(req, -1)
	suite.Equal(fiber.StatusConflict, resp.StatusCode, "Register current user again")

	req = httptest.NewRequest("POST", "/api/user/login", bytes.NewBuffer(userInfoBytes))
	req.Header.Set("Content-Type", fiber.MIMEApplicationJSON)
	resp, _ = suite.server.App.Test(req, -1)
	suite.Equal(fiber.StatusOK, resp.StatusCode, "Login current user")
}
