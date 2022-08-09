package handlers

import (
	"database/sql"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"loyalty-service/internal/database/models"
)

func (server *Server) userRegister(c *fiber.Ctx) (err error) {
	registerData := struct {
		Login    string `form:"login"`
		Password string `form:"password"`
	}{}
	if err = c.BodyParser(&registerData); err != nil {
		err = fiber.NewError(fiber.StatusBadRequest, err.Error())
		return
	}

	user := models.User{}
	err = user.GetOneByLogin(server.DB, registerData.Login)
	if !errors.Is(sql.ErrNoRows, err) {
		if err == nil {
			err = fiber.NewError(fiber.StatusConflict, "login is already taken")
			return
		}

		return err
	}

	passwordBcrypt, _ := bcrypt.GenerateFromPassword([]byte(registerData.Password), bcrypt.DefaultCost)
	newUser := models.User{
		Login:    registerData.Login,
		Password: string(passwordBcrypt),
		Balance:  0,
	}

	err = newUser.Insert(server.DB)
	if err != nil {
		return
	}

	tokenLifime := time.Now().Add(server.Config.TokenLifetime)
	tokenJWT, err := newUser.TokenJWT(tokenLifime, server.Config.Secret)
	if err != nil {
		return
	}

	cookie := fiber.Cookie{
		Name:     "Bearer",
		Value:    tokenJWT,
		Expires:  time.Now().Add(server.Config.TokenLifetime),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"status": "ok",
	})
}

func (server *Server) userLogin(c *fiber.Ctx) (err error) {
	loginData := struct {
		Login    string `form:"login"`
		Password string `form:"password"`
	}{}
	if err = c.BodyParser(&loginData); err != nil {
		err = fiber.NewError(fiber.StatusBadRequest, err.Error())
		return
	}

	user := models.User{}
	err = user.GetOneByLogin(server.DB, loginData.Login)
	if errors.Is(sql.ErrNoRows, err) {
		err = fiber.NewError(fiber.StatusUnauthorized, "user not found")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password))
	if err != nil {
		err = fiber.NewError(fiber.StatusUnauthorized, "wrong password")
		return
	}

	tokenLifime := time.Now().Add(server.Config.TokenLifetime)
	tokenJWT, err := user.TokenJWT(tokenLifime, server.Config.Secret)
	if err != nil {
		return
	}

	cookie := fiber.Cookie{
		Name:     "Bearer",
		Value:    tokenJWT,
		Expires:  time.Now().Add(server.Config.TokenLifetime),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"status": "ok",
	})
}
