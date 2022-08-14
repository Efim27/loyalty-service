package tests

import (
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/suite"
	"loyalty-service/internal/handlers"
)

type OrderTestingSuite struct {
	suite.Suite
	server handlers.Server
	client *resty.Client
}

func (suite *OrderTestingSuite) SetupSuite() {
}

func (suite *OrderTestingSuite) TearDownSuite() {
}

func TestOrderSuite(t *testing.T) {
	suite.Run(t, new(OrderTestingSuite))
}

func (suite *OrderTestingSuite) TestNewOrderWrongNumber() {
	const wrongOrderNum = "12341234"
	resp, err := suite.client.R().SetPathParams(map[string]string{
		"addr": suite.server.Config.ServerAddr,
	}).SetBody(wrongOrderNum).Post("http://{addr}/api/user/orders")
	if err != nil {
		suite.Fail("send request error", err)
	}

	suite.Equal(fiber.StatusUnprocessableEntity, resp.StatusCode(), "New order wrong number test")
}

func (suite *OrderTestingSuite) TestNewOrder() {
	const wrongOrderNum = "12345678903"
	resp, err := suite.client.R().SetPathParams(map[string]string{
		"addr": suite.server.Config.ServerAddr,
	}).SetBody(wrongOrderNum).Post("http://{addr}/api/user/orders")
	if err != nil {
		suite.Fail("send request error", err)
	}

	suite.Contains([]int{fiber.StatusOK, fiber.StatusAccepted, fiber.StatusConflict}, resp.StatusCode(), "New order test")
}

func (suite *OrderTestingSuite) TestOrderList() {
	resp, err := suite.client.R().SetPathParams(map[string]string{
		"addr": suite.server.Config.ServerAddr,
	}).Get("http://{addr}/api/user/orders")
	if err != nil {
		suite.Fail("send request error", err)
	}

	suite.Contains([]int{fiber.StatusOK, fiber.StatusNoContent}, resp.StatusCode(), "Order list test")
}
