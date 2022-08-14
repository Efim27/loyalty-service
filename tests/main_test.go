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

func (mainSuite *MainTestingSuite) SetupSuite() {
	mainSuite.RunServer()

	mainSuite.client = clienthttp.NewClientHTTP(mainSuite.server.Config.HTTPClient)
	mainSuite.ClientAuth()
}

func (mainSuite *MainTestingSuite) TearDownSuite() {
	mainSuite.serverCtxCancel()
	err := mainSuite.server.App.Shutdown()
	if err != nil {
		mainSuite.T().Log(err)
	}
}

func (mainSuite *MainTestingSuite) RunServer() {
	ctx, cancel := context.WithCancel(context.Background())
	mainSuite.serverCtxCancel = cancel
	mainSuite.server = handlers.NewServer()
	mainSuite.server.Prepare(ctx)
	go mainSuite.server.Run()
}

func (mainSuite *MainTestingSuite) ClientAuth() {
	userInfo := handlers.AuthInput{
		Login:    randomString(5),
		Password: randomString(10),
	}
	resp, err := mainSuite.client.R().SetBody(userInfo).SetPathParams(map[string]string{
		"addr": mainSuite.server.Config.ServerAddr,
	}).Post("http://{addr}/api/user/register")
	if err != nil || resp.StatusCode() != fiber.StatusOK {
		mainSuite.Fail("auth failed", err)
	}
	mainSuite.client.SetCookies(resp.Cookies())
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
