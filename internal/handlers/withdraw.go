package handlers

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"loyalty-service/internal/database/models"
)

type withdrawalInput struct {
	OrderNum string  `json:"order"`
	Sum      float64 `json:"sum"`
}

// @title Balance
// @Summary Balance
// @Tags Withdraw
// @Produce json
// @Success 200
// @Failure 401
// @Failure 500
// @Router /api/user/balance [get]
func (server *Server) balanceGet(c *fiber.Ctx) (err error) {
	tokenClaims, ok := c.Locals("tokenClaims").(*jwt.RegisteredClaims)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	userID, err := strconv.Atoi(tokenClaims.Issuer)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	user := models.User{}
	err = user.GetOne(server.DB, uint32(userID))
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	var userWithdrawn float64 = 0
	userWithdrawnNull, err := models.Withdrawal{}.GetSumByUser(server.DB, user.ID)
	if err != nil {
		return
	}

	if userWithdrawnNull.Valid {
		userWithdrawn = userWithdrawnNull.Float64
	}

	return c.JSON(fiber.Map{
		"current":   user.Balance,
		"withdrawn": userWithdrawn,
	})
}

// @title Withdraw List
// @Summary Withdraw List
// @Tags Withdraw
// @Produce json
// @Success 200
// @Success 204
// @Failure 401
// @Failure 500
// @Router /api/user/withdrawals [get]
func (server *Server) withdrawalList(c *fiber.Ctx) (err error) {
	tokenClaims, ok := c.Locals("tokenClaims").(*jwt.RegisteredClaims)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	userID, err := strconv.Atoi(tokenClaims.Issuer)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	orderList, err := models.Withdrawal{}.GetAllByUserSortTime(server.DB, uint32(userID))
	if len(orderList) == 0 {
		return c.SendStatus(fiber.StatusNoContent)
	}
	if err != nil {
		return
	}

	return c.JSON(orderList)
}

// @title Withdraw
// @Summary Withdraw
// @Tags Withdraw
// @Accept json
// @Produce json
// @Param JSON body handlers.withdrawalInput true "JSON"
// @Success 200
// @Failure 401
// @Failure 422
// @Failure 500
// @Router /api/user/balance/withdraw [post]
func (server *Server) withdrawalNew(c *fiber.Ctx) (err error) {
	tokenClaims, ok := c.Locals("tokenClaims").(*jwt.RegisteredClaims)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	userID, err := strconv.Atoi(tokenClaims.Issuer)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	user := models.User{}
	err = user.GetOne(server.DB, uint32(userID))
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	withdrawalData := withdrawalInput{}
	if err = c.BodyParser(&withdrawalData); err != nil {
		err = fiber.NewError(fiber.StatusBadRequest, err.Error())
		return
	}

	orderNumInt, err := strconv.ParseInt(withdrawalData.OrderNum, 10, 64)
	if err != nil {
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}

	err = user.Withdraw(server.DB, uint64(orderNumInt), withdrawalData.Sum)
	if errors.Is(err, models.ErrOrderNumberLunaFailed) {
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}
	if errors.Is(err, models.ErrSumMustBeGreaterThanBalance) {
		return c.SendStatus(fiber.StatusPaymentRequired)
	}
	if err != nil {
		return err
	}

	return
}
