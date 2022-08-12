package models

import (
	"errors"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"loyalty-service/internal/utils/luna"
)

var ErrOrderNumberLunaFailed = errors.New("luna check failed")

const (
	OrderStatusNew        = "NEW"
	OrderStatusProcessing = "PROCESSING"
	OrderStatusInvalid    = "INVALID"
	OrderStatusProcessed  = "PROCESSED"
)

type Order struct {
	Id        uint32    `db:"id" json:"-" form:"id"`
	Number    string    `db:"number" json:"number" form:"number"`
	UserID    uint32    `db:"user_id" json:"-" form:"user_id"`
	Status    string    `db:"status" json:"status" form:"status"`
	Accrual   float64   `db:"accrual" json:"accrual,omitempty" form:"accrual"`
	CreatedAt time.Time `db:"created_at" json:"uploaded_at" form:"uploaded_at"`
}

func NewOrder(number uint64) (order *Order, err error) {
	order = &Order{}
	if !luna.Luna(int(number)) {
		err = ErrOrderNumberLunaFailed
		return
	}

	order.Number = strconv.FormatInt(int64(number), 10)
	return
}

func (Order) GetAll(DB *sqlx.DB) ([]Order, error) {
	allRows := []Order{}

	return allRows, DB.Select(&allRows, `SELECT * FROM "order"`)
}

func (Order) GetAllByUserSortTime(DB *sqlx.DB, userID uint32) ([]Order, error) {
	allRows := []Order{}

	return allRows, DB.Select(&allRows, `SELECT * FROM "order" WHERE user_id=$1 ORDER BY created_at DESC`, userID)
}

func (order *Order) GetOne(DB *sqlx.DB, id uint32) error {
	return DB.Get(order, `SELECT * FROM "order" WHERE id=$1`, id)
}

func (order *Order) GetOneByNumber(DB *sqlx.DB, number string) error {
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
	orderNum, err := strconv.Atoi(order.Number)
	if err != nil {
		return false
	}

	return luna.Luna(orderNum)
}
