package main

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// Получаем текущую рабочую директорию
// Создаем путь к файлу scheduler.db в текущей директории
func getDbfile() string {

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	dbFile := filepath.Join(wd, "scheduler.db")
	return dbFile
}

// Проверяем, существует ли файл
func getDb(dbFile string) (db *sql.DB) {
	_, err := os.Stat(dbFile)
	var install bool
	if err != nil {
		install = true
	}
	// Если файл не существует, создаем его
	if install {
		file, err := os.Create(dbFile)
		if err != nil {
			log.Fatal(err)
		}
		file.Close()
		log.Println("Создан файл базы данных scheduler.db")
	}
	// Открываем соединение с базой данных
	db, err = sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Создаем таблицу, если она не существует
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS scheduler 
	(id INTEGER PRIMARY KEY AUTOINCREMENT, 
	date CHAR(8) NOT NULL DEFAULT '',
	title VARCHAR(256) NOT NULL DEFAULT '',
	comment TEXT NOT NULL DEFAULT '',
	repeat VARCHAR(128) NOT NULL DEFAULT '');
	CREATE INDEX IF NOT EXISTS index_scheduler_date
	ON scheduler (date);
	`)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
