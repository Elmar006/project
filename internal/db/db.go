package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date TEXT NOT NULL,
    title TEXT NOT NULL,
    comment TEXT NOT NULL,
    repeat VARCHAR(128)
);

CREATE INDEX idx_sch_date ON scheduler(date);
`

var DB *sql.DB

func Init(dbFile string) error {
	_, err := os.Stat(dbFile)
	instal := os.IsNotExist(err)

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("Ошибка при открытии БД: %v", err)
	}
	if err := db.Ping(); err != nil {
		return fmt.Errorf("Не удалось подключится к БД: %v", err)
	}
	if instal == true {
		_, err := db.Exec(schema)
		if err != nil {
			return fmt.Errorf("Ошибка инициализации БВ: %v", err)
		}
		fmt.Println("БД созданна")
	}

	DB = db
	return nil
}
