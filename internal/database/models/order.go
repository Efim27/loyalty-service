package models

import (
	"time"

	"github.com/jmoiron/sqlx"
	"loyalty-service/internal/utils/luna"
)

const (
	OrderStatusNew        = "NEW"
	OrderStatusProcessing = "PROCESSING"
	OrderStatusInvalid    = "INVALID"
	OrderStatusProcessed  = "PROCESSED"
)

type Order struct {
	Id        uint32    `db:"id" json:"-" form:"id"`
	Number    uint64    `db:"number" json:"number" form:"number"`
	UserID    uint32    `db:"user_id" json:"user_id" form:"user_id"`
	Status    string    `db:"status" json:"status" form:"status"`
	Accrual   float64   `db:"accrual" json:"accrual" form:"accrual"`
	CreatedAt time.Time `db:"created_at" json:"uploaded_at" form:"uploaded_at"`
}

func (Order) GetAll(DB *sqlx.DB) ([]Order, error) {
	allRows := []Order{}

	return allRows, DB.Select(&allRows, `SELECT * FROM "order"`)
}

func (order *Order) GetOne(DB *sqlx.DB, id uint32) error {
	return DB.Get(order, `SELECT * FROM "order" WHERE id=$1`, id)
}

func (order *Order) GetOneByNumber(DB *sqlx.DB, number uint64) error {
	return DB.Get(order, `SELECT * FROM "order" WHERE "number"=$1`, number)
}

func (order *Order) Insert(DB *sqlx.DB) error {
	return order.insertOne(DB, nil)
}

func (order *Order) insertOne(DB *sqlx.DB, tx *sqlx.Tx) error {
	isSetTx := true
	if tx == nil {
		var err error
		tx, err = DB.Beginx()
		if err != nil {
			return err
		}

		isSetTx = false
		defer tx.Rollback()
	}

	order.Status = OrderStatusNew
	result, err := DB.NamedQuery(`INSERT INTO "order"("number", user_id, status, accrual, created_at) VALUES (:number, :user_id, :status, :accrual, CURRENT_TIMESTAMP) RETURNING id`, order)
	if err != nil {
		return err
	}
	defer result.Close()
	result.Next()
	err = result.Scan(&order.Id)

	if !isSetTx {
		tx.Commit()
	}

	return err
}

func (order Order) InsertMany(DB *sqlx.DB, objectList []Order) error {
	tx, err := DB.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, object := range objectList {
		err = object.insertOne(DB, tx)
		if err != nil {
			return err
		}
	}
	tx.Commit()

	return nil
}

func (order Order) Update(DB *sqlx.DB, newObject Order, tx *sqlx.Tx) error {
	isSetTx := true
	if tx == nil {
		var err error
		tx, err = DB.Beginx()
		if err != nil {
			return err
		}

		isSetTx = false
		defer tx.Rollback()
	}

	newObject.Id = order.Id
	_, err := DB.NamedExec(`UPDATE "order"
	SET number=:number, status=:status, accrual=:accrual, created_at=:created_at, user_id=:user_id
	WHERE id=:id`, newObject)

	if !isSetTx {
		tx.Commit()
	}

	return err
}

func (order Order) Delete(DB *sqlx.DB) error {
	_, err := DB.Exec(`DELETE FROM "order" WHERE id=$1;`, order.Id)

	return err
}

func (order Order) CheckLuna() bool {
	return luna.Luna(int(order.Number))
}
