package models

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jmoiron/sqlx"
)

var ErrSumMustBePositive = errors.New("sum must be positive")
var ErrSumMustBeGreaterThanBalance = errors.New("sum must be greater than the user balance")

type User struct {
	ID       uint32  `db:"id" json:"id" form:"id"`
	Login    string  `db:"login" json:"login" form:"login"`
	Password string  `db:"password" json:"password" form:"-"`
	Balance  float64 `db:"balance" json:"balance" form:"balance"`
}

func (User) GetAll(DB *sqlx.DB) ([]User, error) {
	allRows := []User{}

	return allRows, DB.Select(&allRows, `SELECT * FROM user`)
}

func (user *User) GetOne(DB *sqlx.DB, id uint32) error {
	return DB.Get(user, `SELECT * FROM "user" WHERE id=$1`, id)
}

func (user *User) GetOneByLogin(DB *sqlx.DB, login string) error {
	return DB.Get(user, `SELECT * FROM "user" WHERE login=$1`, login)
}

func (user *User) Insert(DB *sqlx.DB) error {
	return user.insertOne(DB, nil)
}

func (user *User) insertOne(DB *sqlx.DB, tx *sqlx.Tx) error {
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

	result, err := DB.NamedQuery(`INSERT INTO "user"(login, password, balance) VALUES (:login, :password, :balance) RETURNING id`, user)
	if err != nil {
		return err
	}
	defer result.Close()
	result.Next()
	err = result.Scan(&user.ID)

	if !isSetTx {
		tx.Commit()
	}

	return err
}

func (user User) InsertMany(DB *sqlx.DB, objectList []User) error {
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

func (user User) Update(DB *sqlx.DB, newObject User, tx *sqlx.Tx) error {
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

	newObject.ID = user.ID
	_, err := DB.NamedExec(`UPDATE "user"
	SET login=:login, password=:password, balance=:balance
	WHERE id=:id`, newObject)

	if !isSetTx {
		tx.Commit()
	}

	return err
}

func (user User) Delete(DB *sqlx.DB) error {
	_, err := DB.Exec(`DELETE FROM "user" WHERE id=$1;`, user.ID)

	return err
}

func (user User) TokenJWT(expiresAt time.Time, secret string) (token string, err error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    strconv.Itoa(int(user.ID)),
		ExpiresAt: jwt.NewNumericDate(expiresAt),
	})

	token, err = claims.SignedString([]byte(secret))
	return
}

func (user User) UpdateBalance(DB *sqlx.DB, delta float64, tx *sqlx.Tx) (err error) {
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

	_, err = DB.Exec(`UPDATE "user" SET balance=balance+$1 WHERE id=$2;`, delta, user.ID)

	if !isSetTx {
		tx.Commit()
	}
	return
}

func (user User) Withdraw(DB *sqlx.DB, orderNum uint64, sum float64) (err error) {
	if sum < 0 {
		return ErrSumMustBePositive
	}

	if user.Balance < sum {
		return ErrSumMustBeGreaterThanBalance
	}

	//Luna check
	_, err = NewOrder(orderNum)
	if err != nil {
		return
	}

	tx, err := DB.Beginx()
	defer tx.Rollback()
	if err != nil {
		return
	}

	withdrawal := Withdrawal{
		UserID:      user.ID,
		OrderNumber: strconv.FormatInt(int64(orderNum), 10),
		Sum:         sum,
	}

	err = withdrawal.InsertOne(DB, tx)
	if err != nil {
		return
	}

	err = user.UpdateBalance(DB, -sum, tx)
	if err != nil {
		return
	}

	tx.Commit()
	return
}
