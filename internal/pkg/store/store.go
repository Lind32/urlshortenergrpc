package store

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Lind-32/urlshortenergrpc/internal/config"
	_ "github.com/lib/pq" //...
	"github.com/pressly/goose"
)

var db *sql.DB

//Close закрыть базу данных
func Close() {
	db.Close()
}

//Connect подключение к базе данных
func Connect(cfg config.Config) error {

	p := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.Data.PgHost, cfg.Data.PgPort, cfg.Data.PgUser, cfg.Data.PgPass, cfg.Data.PgName)
	log.Printf("DB connection...\n")
	var err error
	db, err = sql.Open("postgres", p)
	if err != nil {
		return err
	}
	log.Printf("Connected: %s:%s DB: %s\n", cfg.Data.PgHost, cfg.Data.PgPort, cfg.Data.PgName)

	//очистка БД
	if cfg.Data.PgClean {
		log.Printf("Cleaning DB...\n")
		err := goose.DownTo(db, ".", 0)
		if err != nil {
			return err
		}
	}
	//миграция БД
	log.Printf("Migration DB...\n")
	err = goose.Up(db, ".")
	if err != nil {
		return err
	}

	return nil
}

//Insert запись в базу данных]
func Insert(key, llink string) error {

	var err error
	_, err = db.Exec(`INSERT INTO "shorts_url" (key, long_link) VALUES ($1, $2)`, key, llink)
	if err != nil {
		return err
	}

	return nil
}

//GetLongURL возвращает длинную ссылку по ключу
func GetLongURL(key string) (string, error) {
	var l string
	r := db.QueryRow(`SELECT long_link FROM "shorts_url" WHERE key = $1`, key)
	err := r.Scan(&l)
	if err != nil {
		if err != sql.ErrNoRows {
			return "", err
		}
		return "", nil
	}
	return l, nil
}

//UnicURL проверка длинной ссылки на уникальность (true если уникальна)
func UnicURL(lurl string) (string, bool, error) {
	var key string
	r := db.QueryRow(`SELECT key FROM "shorts_url" WHERE long_link = $1`, lurl)
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

	r := db.QueryRow(`SELECT key FROM "shorts_url" WHERE key = $1`, key)
	err := r.Scan(&key)
	if err != nil {
		if err != sql.ErrNoRows {
			return true, err
		}
		return true, nil
	}
	return false, nil
}
