package tests

import (
	"encoding/json"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/suite"
	"loyalty-service/internal/database/models"
	"loyalty-service/internal/handlers"
)

type WithdrawTestingSuite struct {
	suite.Suite
	server handlers.Server
	client *resty.Client
}

func (suite *WithdrawTestingSuite) SetupSuite() {
}

func (suite *WithdrawTestingSuite) TearDownSuite() {
}

func TestWithdrawSuite(t *testing.T) {
	suite.Run(t, new(WithdrawTestingSuite))
}

func (suite *WithdrawTestingSuite) TestBalance() {
	resp, err := suite.client.R().SetPathParams(map[string]string{
		"addr": suite.server.Config.ServerAddr,
	}).Get("http://{addr}/api/user/balance")
	if err != nil {
		suite.Fail("send request error", err)
	}

	suite.Equal(fiber.StatusOK, resp.StatusCode(), "User balance test")

	var userBalance struct {
		Current   float64 `json:"current"`
		Withdrawn float64 `json:"withdrawn"`
	}
	err = json.Unmarshal(resp.Body(), &userBalance)
	if err != nil {
		suite.Fail("unmarshal response error", err)
	}
}

func (suite *WithdrawTestingSuite) TestWrongRequestWithdraw() {
	resp, err := suite.client.R().SetPathParams(map[string]string{
		"addr": suite.server.Config.ServerAddr,
	}).Post("http://{addr}/api/user/balance/withdraw")
	if err != nil {
		suite.Fail("send request error", err)
	}

	suite.Equal(fiber.StatusBadRequest, resp.StatusCode(), "Wrong request withdraw test")
}

func (suite *WithdrawTestingSuite) TestNotEnoughBalanceWithdraw() {
	resp, err := suite.client.R().SetPathParams(map[string]string{
		"addr": suite.server.Config.ServerAddr,
	}).SetBody(map[string]interface{}{
		"order": "1234567890",
		"sum":   999,
	}).Post("http://{addr}/api/user/balance/withdraw")
	if err != nil {
		suite.Fail("send request error", err)
	}

	suite.Equal(fiber.StatusPaymentRequired, resp.StatusCode(), "Not enough balance withdraw test")
}

func (suite *WithdrawTestingSuite) TestWrongOrderNumWithdraw() {
	resp, err := suite.client.R().SetPathParams(map[string]string{
		"addr": suite.server.Config.ServerAddr,
	}).SetBody(map[string]interface{}{
		"order": "12341234",
		"sum":   0,
	}).Post("http://{addr}/api/user/balance/withdraw")
	if err != nil {
		suite.Fail("send request error", err)
	}

	suite.Equal(fiber.StatusUnprocessableEntity, resp.StatusCode(), "Wrong order number withdraw test")
}

func (suite *WithdrawTestingSuite) TestOrderList() {
	resp, err := suite.client.R().SetPathParams(map[string]string{
		"addr": suite.server.Config.ServerAddr,
	}).Get("http://{addr}/api/user/withdrawals")
	if err != nil {
		suite.Fail("send request error", err)
	}

	suite.Contains([]int{fiber.StatusOK, fiber.StatusNoContent}, resp.StatusCode(), "Order list test")
	if resp.StatusCode() != fiber.StatusOK {
		return
	}

	var orderList []models.Order
	err = json.Unmarshal(resp.Body(), &orderList)
	if err != nil {
		suite.Fail("unmarshal response error", err)
	}
}
