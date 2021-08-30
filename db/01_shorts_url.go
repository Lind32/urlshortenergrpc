package db

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(Up, Down)
}

//Up генерация таблицы
func Up(tx *sql.Tx) error {

	_, err := tx.Exec(`
	CREATE TABLE IF NOT EXISTS shorts_url(
		"id" SERIAL PRIMARY KEY,
		"key" VARCHAR(10),
		"long_link" VARCHAR);
		`)
	if err != nil {
		return err
	}

	return nil
}

//Down дроп таблицы
func Down(tx *sql.Tx) error {

	_, err := tx.Exec(`
	DROP TABLE shorts_url;
		`)
	if err != nil {
		return err
	}

	return nil
}
