package models

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

type Withdrawal struct {
	ID          uint32    `db:"id" json:"-" form:"id"`
	OrderNumber string    `db:"order_number" json:"order" form:"order"`
	UserID      uint32    `db:"user_id" json:"-" form:"user_id"`
	Sum         float64   `db:"sum" json:"sum" form:"sum"`
	ProcessedAt time.Time `db:"processed_at" json:"processed_at" form:"processed_at"`
}

func (Withdrawal) GetAll(DB *sqlx.DB) ([]Withdrawal, error) {
	allRows := []Withdrawal{}

	return allRows, DB.Select(&allRows, `SELECT * FROM "withdrawal"`)
}

func (Withdrawal) GetAllByUserSortTime(DB *sqlx.DB, userID uint32) ([]Withdrawal, error) {
	allRows := []Withdrawal{}

	return allRows, DB.Select(&allRows, `SELECT * FROM "withdrawal" WHERE user_id=$1 ORDER BY processed_at DESC`, userID)
}

func (withdrawal *Withdrawal) GetOne(DB *sqlx.DB, id uint32) error {
	return DB.Get(withdrawal, `SELECT * FROM "withdrawal" WHERE id=$1`, id)
}

func (withdrawal *Withdrawal) GetOneByOrderNumber(DB *sqlx.DB, orderNumber uint32) error {
	return DB.Get(withdrawal, `SELECT * FROM "withdrawal" WHERE order_number=$1`, orderNumber)
}

func (withdrawal *Withdrawal) Insert(DB *sqlx.DB) error {
	return withdrawal.InsertOne(DB, nil)
}

func (withdrawal *Withdrawal) InsertOne(DB *sqlx.DB, tx *sqlx.Tx) error {
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

	result, err := DB.NamedQuery(`INSERT INTO "withdrawal"(order_number, user_id, sum, processed_at) VALUES (:order_number, :user_id, :sum, CURRENT_TIMESTAMP) RETURNING id`, withdrawal)
	if err != nil {
		return err
	}
	err = result.Err()
	if err != nil {
		return err
	}
	defer result.Close()
	result.Next()
	err = result.Scan(&withdrawal.ID)

	if !isSetTx {
		tx.Commit()
	}

	return err
}

func (withdrawal Withdrawal) InsertMany(DB *sqlx.DB, objectList []Withdrawal) error {
	tx, err := DB.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, object := range objectList {
		err = object.InsertOne(DB, tx)
		if err != nil {
			return err
		}
	}
	tx.Commit()

	return nil
}

func (withdrawal Withdrawal) Update(DB *sqlx.DB, newObject Withdrawal, tx *sqlx.Tx) error {
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

	newObject.ID = withdrawal.ID
	_, err := DB.NamedExec(`UPDATE "withdrawal"
	SET order_number=:order_number, user_id=:user_id, sum=:sum, processed_at=:processed_at
	WHERE id=:id`, newObject)

	if !isSetTx {
		tx.Commit()
	}

	return err
}

func (withdrawal Withdrawal) Delete(DB *sqlx.DB) error {
	_, err := DB.Exec(`DELETE FROM "withdrawal" WHERE id=$1;`, withdrawal.ID)

	return err
}

func (withdrawal Withdrawal) GetSumByUser(DB *sqlx.DB, userID uint32) (userSum sql.NullFloat64, err error) {
	result, err := DB.Query(`SELECT SUM(sum) FROM withdrawal WHERE user_id=$1;`, userID)
	if err != nil {
		return
	}
	err = result.Err()
	if err != nil {
		return
	}
	defer result.Close()

	result.Next()
	err = result.Scan(&userSum)

	return
}
