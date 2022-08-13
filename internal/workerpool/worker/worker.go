package worker

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"loyalty-service/internal/database/models"
	main_logger "loyalty-service/internal/logger"
)

var (
	ErrTooManyRequests            = errors.New("too many requests")
	ErrAccrualHandlingNotFinished = errors.New("accrual handling not finished")
)

type OutputWorker struct {
	Ch                chan string
	DB                *sqlx.DB
	AccrualServerAddr string
	Client            *resty.Client
	Logger            *main_logger.Logger
}

func (w *OutputWorker) handleAccrual(orderNum string) (retryAfter time.Duration, err error) {
	resp, err := w.Client.R().
		SetPathParams(map[string]string{
			"addr":     w.AccrualServerAddr,
			"orderNum": orderNum,
		}).
		Get("{addr}/api/orders/{orderNum}")
	if err != nil {
		return
	}

	respStatusCode := resp.StatusCode()
	if respStatusCode == fiber.StatusTooManyRequests {
		retryAfterStr := resp.Header().Get("Retry-After")
		var retryAfterNumSeconds int
		retryAfterNumSeconds, err = strconv.Atoi(retryAfterStr)
		if err != nil {
			return
		}

		err = ErrTooManyRequests
		retryAfter = time.Duration(retryAfterNumSeconds) * time.Second
		return
	}

	if respStatusCode == fiber.StatusOK {
		order := models.Order{}
		err = order.GetOneByNumber(w.DB, orderNum)
		if err != nil {
			return
		}

		err = json.Unmarshal(resp.Body(), &order)
		if err != nil {
			return
		}

		if order.Status == models.OrderStatusRegistered || order.Status == models.OrderStatusProcessing {
			err = ErrAccrualHandlingNotFinished
			return
		}

		if order.Status == models.OrderStatusNew {
			return
		}

		err = order.Update(w.DB, order, nil)
		if err != nil {
			return
		}

		user := models.User{}
		err = user.GetOne(w.DB, order.UserID)
		if err != nil {
			return
		}
		err = user.UpdateBalance(w.DB, order.Accrual, nil)
		if err != nil {
			return
		}

		return
	}

	errMsg := fmt.Sprintf("unhandled status code %v", respStatusCode)
	err = errors.New(errMsg)
	return
}

func (w *OutputWorker) Do() {
	for orderNum := range w.Ch {
		retryAfter, err := w.handleAccrual(orderNum)
		if errors.Is(err, ErrAccrualHandlingNotFinished) {
			w.Ch <- orderNum
		}

		if errors.Is(err, ErrTooManyRequests) {
			time.Sleep(retryAfter)
			w.Ch <- orderNum
		}

		if err != nil {
			w.Logger.Error("handler accrual error", zap.Error(err))
		}
	}
}
