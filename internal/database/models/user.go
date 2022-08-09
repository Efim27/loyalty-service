package models

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jmoiron/sqlx"
)

type User struct {
	Id       uint32  `db:"id" json:"id" form:"id"`
	Login    string  `db:"login" json:"login" form:"login"`
	Password string  `db:"password" json:"password" form:"-"`
	Balance  float64 `db:"balance" json:"balance" form:"balance"`
}

func (User) GetAll(DB *sqlx.DB) ([]User, error) {
	allRows := []User{}

	return allRows, DB.Select(&allRows, "SELECT * FROM user")
}

func (user *User) GetOne(DB *sqlx.DB, id uint32) error {
	return DB.Get(user, "SELECT * FROM \"user\" WHERE id=$1", id)
}

func (user *User) GetOneByLogin(DB *sqlx.DB, login string) error {
	return DB.Get(user, "SELECT * FROM \"user\" WHERE login=$1", login)
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
	err = result.Scan(&user.Id)

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

	newObject.Id = user.Id
	_, err := DB.NamedExec(`UPDATE "user"
	SET login=:login, password=:password, balance=:balance
	WHERE id=:id`, newObject)

	if !isSetTx {
		tx.Commit()
	}

	return err
}

func (user User) Delete(DB *sqlx.DB) error {
	_, err := DB.Exec(`DELETE FROM "user" WHERE id=$1;`, user.Id)

	return err
}

func (user User) TokenJWT(expiresAt time.Time, secret string) (token string, err error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    strconv.Itoa(int(user.Id)),
		ExpiresAt: jwt.NewNumericDate(expiresAt),
	})

	token, err = claims.SignedString([]byte(secret))
	return
}
