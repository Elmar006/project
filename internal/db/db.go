package db

import (
	"database/sql"
	"os"

	"github.com/Elmar006/project/internal/logger"
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
		logger.L().Errorf("Ошибка при открытии БД: %v", err)
		return err
	}
	if err := db.Ping(); err != nil {
		logger.L().Errorf("Не удалось подключится к БД: %v", err)
		return err
	}
	if instal == true {
		_, err := db.Exec(schema)
		if err != nil {
			logger.L().Errorf("Ошибка инициализации БД: %v", err)
			return err
		}
		logger.L().Info("БД созданна")
	}

	DB = db
	return nil
}
