package handlers

import (
	"database/sql"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"loyalty-service/internal/database/models"
)

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

	orderByNum := models.Order{
		Number: uint64(orderNum),
	}
	if !orderByNum.CheckLuna() {
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}

	err = orderByNum.GetOneByNumber(server.DB, orderByNum.Number)
	if err != nil && err != sql.ErrNoRows {
		return
	} else if orderByNum.UserID == uint32(userID) {
		return c.SendStatus(fiber.StatusOK)
	} else if orderByNum.UserID != 0 {
		return c.SendStatus(fiber.StatusConflict)
	}

	newOrder := models.Order{
		Number: uint64(orderNum),
		UserID: uint32(userID),
	}
	err = newOrder.Insert(server.DB)
	if err != nil {
		return
	}

	return c.SendStatus(fiber.StatusAccepted)
}

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
