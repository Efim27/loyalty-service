package handlers

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap/zapcore"
	"loyalty-service/internal/config"
	"loyalty-service/internal/database"
	main_logger "loyalty-service/internal/logger"
)

type Server struct {
	App    *fiber.App
	DB     *sqlx.DB
	Config config.Config
	Logger *main_logger.Logger
}

func fiberErrorHandler(ctx *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	return ctx.Status(code).JSON(fiber.Map{
		"status": "error",
		"error":  err.Error(),
	})
}

func GetFiberConfig() (config fiber.Config) {
	config.ErrorHandler = fiberErrorHandler
	return
}

func NewServer() Server {
	server := Server{
		App: fiber.New(GetFiberConfig()),
	}

	logger, err := main_logger.NewLogger(zapcore.DebugLevel)
	if err != nil {
		log.Fatal(err)
	}
	server.Logger = logger

	mainConfig, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	server.Config = mainConfig

	server.DB = database.NewDatabase(server.Config.DBSource, server.Logger)

	return server
}

func (server *Server) setupMiddlewares() {
	server.App.Use(recover.New())
	server.App.Use(logger.New())
	server.App.Use(compress.New())
}

func (server Server) Run() {
	server.setupMiddlewares()
	SetupRoutes(server)

	server.App.Listen(server.Config.ServerAddr)
}
