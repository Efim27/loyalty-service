package tests

import (
	"context"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/suite"
	"loyalty-service/internal/clienthttp"

	"loyalty-service/internal/handlers"
)

type MainTestingSuite struct {
	suite.Suite
	serverCtxCancel context.CancelFunc
	server          handlers.Server
	client          *resty.Client
}

func (suite *MainTestingSuite) SetupSuite() {
	suite.RunServer()

	suite.client = clienthttp.NewClientHTTP(suite.server.Config.HTTPClient)
	suite.ClientAuth()
}

func (suite *MainTestingSuite) TearDownSuite() {
	suite.serverCtxCancel()
	err := suite.server.App.Shutdown()
	if err != nil {
		suite.T().Log(err)
	}
}

func (suite *MainTestingSuite) RunServer() {
	ctx, cancel := context.WithCancel(context.Background())
	suite.serverCtxCancel = cancel
	suite.server = handlers.NewServer()
	suite.server.Prepare(ctx)
	go suite.server.Run()
}

func (suite *MainTestingSuite) ClientAuth() {
	userInfo := handlers.AuthInput{
		Login:    randomString(5),
		Password: randomString(10),
	}
	resp, err := suite.client.R().SetBody(userInfo).SetPathParams(map[string]string{
		"addr": suite.server.Config.ServerAddr,
	}).Post("http://{addr}/api/user/register")
	if err != nil || resp.StatusCode() != fiber.StatusOK {
		suite.Fail("auth failed", err)
	}
	suite.client.SetCookies(resp.Cookies())
}

func TestAllSuites(t *testing.T) {
	suite.Run(t, &MainTestingSuite{})
}

func (mainSuite *MainTestingSuite) TestAllSuites() {
	suite.Run(mainSuite.T(), &AuthTestingSuite{
		server: mainSuite.server,
	})

	suite.Run(mainSuite.T(), &OrderTestingSuite{
		server: mainSuite.server,
		client: mainSuite.client,
	})

	suite.Run(mainSuite.T(), &WithdrawTestingSuite{
		server: mainSuite.server,
		client: mainSuite.client,
	})
}
