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
	_ "loyalty-service/docs"
	client_http "loyalty-service/internal/clienthttp"
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

	mainConfig, err := config.LoadConfig()
	if err != nil {
		server.Logger.Fatal("loading config error", zap.Error(err))
	}
	server.Config = mainConfig

	logger, err := main_logger.NewLogger(zapcore.DebugLevel, server.Config.LogFile)
	if err != nil {
		log.Fatal(err)
	}
	server.Logger = logger
	server.Logger.Info("config loaded", zap.Any("config", server.Config))

	server.DB = database.NewDatabase(server.Config.DBSource, server.Logger)

	return server
}

func (server *Server) setupMiddlewares() {
	server.App.Use(recover.New())
	server.App.Use(logger.New())
	server.App.Use(compress.New())
}

// @title Loyalty Service
// @version 1.1
// @description Cumulative loyalty system "Gofermart"

// @contact.name Efim
// @contact.url https://t.me/hima27
// @contact.email efim-02@mail.ru

// @schemes http
func (server Server) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	clientHTTP := client_http.NewClientHTTP(server.Config.HTTPClient)
	orderAccrualHandler, OrderAccrualHandlerChan := workerpool.NewOrderAccrualHandler(ctx, server.Config.HTTPClient.AccrualAddr, 4, server.DB, server.Logger, clientHTTP)

	server.OrderAccrualHandlerChan = OrderAccrualHandlerChan
	go orderAccrualHandler.Start()

	server.setupMiddlewares()
	SetupRoutes(server)

	err := server.App.Listen(server.Config.ServerAddr)
	if err != nil {
		server.Logger.Error("stopping down server error", zap.Error(err))
	}
}
