package workerpool

import (
	"context"

	"github.com/go-resty/resty/v2"
	"github.com/jmoiron/sqlx"
	main_logger "loyalty-service/internal/logger"
	handler_worker "loyalty-service/internal/workerpool/worker"
)

type OrderAccrualHandler struct {
	ctx               context.Context
	AccrualServerAddr string
	workersCount      int
	handleChan        chan string
	client            *resty.Client
	DB                *sqlx.DB
	logger            *main_logger.Logger
}

func NewOrderAccrualHandler(ctx context.Context, AccrualServerAddr string, workersCount uint, DB *sqlx.DB, logger *main_logger.Logger, client *resty.Client) (orderAccrualHandler OrderAccrualHandler, inputChan chan<- string) {
	handleChan := make(chan string)
	inputChan = handleChan

	orderAccrualHandler.ctx = ctx
	orderAccrualHandler.AccrualServerAddr = AccrualServerAddr
	orderAccrualHandler.DB = DB
	orderAccrualHandler.logger = logger
	orderAccrualHandler.handleChan = handleChan
	orderAccrualHandler.workersCount = int(workersCount)
	orderAccrualHandler.client = client

	return
}

func (orderAccrualHandler OrderAccrualHandler) Start() {
	for i := 0; i < orderAccrualHandler.workersCount; i++ {
		w := &handler_worker.OutputWorker{
			Ch:                orderAccrualHandler.handleChan,
			DB:                orderAccrualHandler.DB,
			AccrualServerAddr: orderAccrualHandler.AccrualServerAddr,
			Client:            orderAccrualHandler.client,
			Logger:            orderAccrualHandler.logger,
		}
		go w.Do()
	}

	<-orderAccrualHandler.ctx.Done()
	close(orderAccrualHandler.handleChan)
}
