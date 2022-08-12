package workerpool

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	main_logger "loyalty-service/internal/logger"
	handler_worker "loyalty-service/internal/workerpool/worker"
)

type OrderAccrualHandler struct {
	ctx            context.Context
	workersCount   int
	retryCountHTTP uint
	handleChan     chan string
	DB             *sqlx.DB
	logger         *main_logger.Logger
	timeoutHTTP    time.Duration
}

func NewOrderAccrualHandler(ctx context.Context, workersCount uint, DB *sqlx.DB, logger *main_logger.Logger, retryCountHTTP uint, timeoutHTTP time.Duration) (inputChan chan<- string, orderAccrualHandler OrderAccrualHandler) {
	handleChan := make(chan string)
	inputChan = handleChan

	orderAccrualHandler.ctx = ctx
	orderAccrualHandler.DB = DB
	orderAccrualHandler.logger = logger
	orderAccrualHandler.handleChan = handleChan
	orderAccrualHandler.workersCount = int(workersCount)
	orderAccrualHandler.retryCountHTTP = retryCountHTTP
	orderAccrualHandler.timeoutHTTP = timeoutHTTP

	return
}

func (orderAccrualHandler OrderAccrualHandler) Start() {
	for i := 0; i < orderAccrualHandler.workersCount; i++ {
		w := &handler_worker.OutputWorker{
			Ch:            orderAccrualHandler.handleChan,
			DB:            orderAccrualHandler.DB,
			TimeoutModbus: orderAccrualHandler.timeoutHTTP}
		go w.Do()
	}

	<-orderAccrualHandler.ctx.Done()
	close(orderAccrualHandler.handleChan)
}
