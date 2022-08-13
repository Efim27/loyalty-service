package handlers

import (
	"database/sql"
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"loyalty-service/internal/database/models"
)

// @title New order
// @Summary New order
// @Tags Order
// @Accept plain
// @Param OrderNumber body string true "Order number (Luna check)"
// @Success 200
// @Success 202
// @Failure 400
// @Failure 401
// @Failure 409
// @Failure 422
// @Failure 500
// @Router /api/user/orders [post]
func (server *Server) orderNew(c *fiber.Ctx) (err error) {
	tokenClaims, ok := c.Locals("tokenClaims").(*jwt.RegisteredClaims)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	userID, err := strconv.Atoi(tokenClaims.Issuer)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	orderNumStr := string(c.Body())
	orderNum, err := strconv.ParseInt(orderNumStr, 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	//Luna check
	orderByNum, err := models.NewOrder(uint64(orderNum))
	if errors.Is(err, models.ErrOrderNumberLunaFailed) {
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}

	err = orderByNum.GetOneByNumber(server.DB, orderByNum.Number)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return
	} else if orderByNum.UserID == uint32(userID) {
		return c.SendStatus(fiber.StatusOK)
	} else if orderByNum.UserID != 0 {
		return c.SendStatus(fiber.StatusConflict)
	}

	newOrder := models.Order{
		Number: orderNumStr,
		UserID: uint32(userID),
	}
	err = newOrder.Insert(server.DB)
	if err != nil {
		return
	}
	server.OrderAccrualHandlerChan <- orderNumStr

	return c.SendStatus(fiber.StatusAccepted)
}

// @title Order list
// @Summary Order list
// @Tags Order
// @Produce json
// @Success 200
// @Success 204
// @Failure 401
// @Failure 500
// @Router /api/user/orders [get]
func (server *Server) orderList(c *fiber.Ctx) (err error) {
	tokenClaims, ok := c.Locals("tokenClaims").(*jwt.RegisteredClaims)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	userID, err := strconv.Atoi(tokenClaims.Issuer)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	orderList, err := models.Order{}.GetAllByUserSortTime(server.DB, uint32(userID))
	if len(orderList) == 0 {
		return c.SendStatus(fiber.StatusNoContent)
	}
	if err != nil {
		return
	}

	return c.JSON(orderList)
}
