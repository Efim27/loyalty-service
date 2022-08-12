package worker

import (
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	controller_handlers "object-control/internal/server/controller-handlers"
	"object-control/internal/server/database/models"
)

type OutputWorker struct {
	Ch            chan models.Controller
	TimeoutModbus time.Duration
	DB            *sqlx.DB
}

func readPocket(controller models.Controller, DB *sqlx.DB, timeout time.Duration) {
	pocketReader, err := controller_handlers.NewControllerHandler(controller, timeout)
	defer func() {
		if err := recover(); err != nil {
			log.Println("cant connect to controller", controller.Id)
			return
		}
	}()
	defer pocketReader.Close()
	if err != nil {
		log.Println("scanner error while connecting to controller", err)
		return
	}

	_, err = pocketReader.ReadPocket(DB)
	if err != nil {
		log.Println("scanner error while reading pocket", err)
	}
}

func (w *OutputWorker) Do() {
	for controller := range w.Ch {
		readPocket(controller, w.DB, w.TimeoutModbus)
	}

	return
}
