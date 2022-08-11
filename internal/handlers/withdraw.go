package handlers

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"loyalty-service/internal/database/models"
)

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

	var userWithdrawn int32 = 0
	userWithdrawnNull, err := models.Withdrawal{}.GetSumByUser(server.DB, user.Id)
	if err != nil {
		return
	}

	if userWithdrawnNull.Valid {
		userWithdrawn = userWithdrawnNull.Int32
	}

	return c.JSON(fiber.Map{
		"current":   user.Balance,
		"withdrawn": userWithdrawn,
	})
}

func (server *Server) withdrawalList(c *fiber.Ctx) (err error) {
	tokenClaims, ok := c.Locals("tokenClaims").(*jwt.RegisteredClaims)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	log.Println()
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
