package handlers

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	client_http "loyalty-service/internal/client-http"
	"loyalty-service/internal/config"
	"loyalty-service/internal/database"
	main_logger "loyalty-service/internal/logger"
	"loyalty-service/internal/workerpool"
)

type Server struct {
	App                     *fiber.App
	DB                      *sqlx.DB
	Config                  config.Config
	Logger                  *main_logger.Logger
	OrderAccrualHandlerChan chan<- string
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
	logger.Info("config loaded", zap.Any("config", server.Config))

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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	clientHTTP := client_http.NewClientHTTP(server.Config.HTTPClient)
	orderAccrualHandler, OrderAccrualHandlerChan := workerpool.NewOrderAccrualHandler(ctx, server.Config.HTTPClient.AccrualAddr, 4, server.DB, server.Logger, clientHTTP)
	server.OrderAccrualHandlerChan = OrderAccrualHandlerChan
	go orderAccrualHandler.Start()

	server.App.Listen(server.Config.ServerAddr)
}
