package store

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" //
)

const (

	//postgres
	pgHost = "127.0.0.1"
	pgPort = "5432"
	pgUser = "postgres"
	pgPass = "postgres"
	pgName = "shortURL"
)

var db *sql.DB

//Close закрыть базу данных
func Close() {
	db.Close()
}

//Connect подключение к базе данных
func Connect() error {

	pgparam := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", pgHost, pgPort, pgUser, pgPass, pgName)

	var err error
	db, err = sql.Open("postgres", pgparam)
	if err != nil {
		return err
	}

	return nil
}

//Insert запись в базу данных]
func Insert(key, longlink string) error {

	var err error
	_, err = db.Exec(`INSERT INTO "URLlist" (key, longlink) VALUES ($1, $2)`, key, longlink)
	if err != nil {
		return err
	}

	return nil
}

//GetLongURL возвращает длинную ссылку по ключу
func GetLongURL(key string) (string, error) {
	var l string
	r := db.QueryRow(`SELECT longlink FROM "URLlist" WHERE key = $1`, key)
	err := r.Scan(&l)
	if err != nil {
		if err != sql.ErrNoRows {
			return "", err
		}
		return "", nil
	}
	return l, nil
}

//UnicURL проверка длинной ссылки на уникальность (true если уникальнfz)
func UnicURL(lurl string) (string, bool, error) {
	var key string
	r := db.QueryRow(`SELECT key FROM "URLlist" WHERE longlink = $1`, lurl)
	err := r.Scan(&key)
	if err != nil {
		if err != sql.ErrNoRows {
			return "", true, err
		}
		return "", true, nil
	}

	return key, false, nil
}

//UnicKey проверка ключа на уникальность (true если уникальный)
func UnicKey(key string) (bool, error) {
	if key == "" {
		return false, nil
	}

	r := db.QueryRow(`SELECT key FROM "URLlist" WHERE key = $1`, key)
	err := r.Scan(&key)
	if err != nil {
		if err != sql.ErrNoRows {
			return true, err
		}
		return true, nil
	}
	return false, nil
}
