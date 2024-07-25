package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "20060102"

type Task struct {
	Id      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return Store{db: db}
}

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
func isDbHere(dbFile string) bool {
	_, err := os.Stat(dbFile)
	return err != nil
}

// Если файл не существует, создаем его
func createDb(dbFile string, db *sql.DB) {
	file, err := os.Create(dbFile)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()
	log.Println("Создан файл базы данных scheduler.db")

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
}

func (s Store) getNextDate(now string, date string, repeat string) (string, error) {
	parsedNow, err := time.Parse(dateFormat, now)
	if err != nil {
		return "", fmt.Errorf("ошибка при парсинге 'now': %w", err)
	}
	nextDate, err := time.Parse(dateFormat, date)
	if err != nil {
		return "", fmt.Errorf("ошибка при парсинге 'date': %w", err)
	}
	if repeat == "" {
		return "", fmt.Errorf("ошибка: параметр 'repeat' пустой")
	}
	repeatParts := strings.Split(repeat, " ")
	firstPart := repeatParts[0]
	switch {
	case firstPart == "y":
		nextDate = nextDate.AddDate(1, 0, 0)
		for nextDate.Before(parsedNow) || nextDate.Equal(parsedNow) {
			nextDate = nextDate.AddDate(1, 0, 0)
		}
		return nextDate.Format(dateFormat), nil
	case firstPart == "d" && len(repeatParts) == 2:
		numOfDays, err := strconv.Atoi(repeatParts[1])
		if err != nil {
			return "", fmt.Errorf("ошибка при преобразовании параметра 'repeat' в int: %w", err)
		}
		if numOfDays > 400 {
			return "", fmt.Errorf("ошибка: число дней превышает 400")
		}
		nextDate = nextDate.AddDate(0, 0, numOfDays)
		for nextDate.Before(parsedNow) || nextDate.Equal(parsedNow) {
			nextDate = nextDate.AddDate(0, 0, numOfDays)
		}
		return nextDate.Format(dateFormat), nil
	default:
		return "", fmt.Errorf("ошибка: неподдерживаемый формат параметра 'repeat'")
	}
}

func (s Store) getAllTasks() ([]Task, error) {
	rows, err := s.db.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT 30")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var allTasks []Task

	for rows.Next() {
		task := Task{}
		if err := rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, err
		}
		allTasks = append(allTasks, task)
	}

	return allTasks, nil
}

func (s Store) postTask(task Task) (int, error) {
	res, err := s.db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (s Store) getTask(id int) (Task, error) {
	task := Task{}

	row := s.db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id", sql.Named("id", id))

	err := row.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return task, err
	}

	return task, nil
}

func (s Store) UpdateTask(task Task) error {
	_, err := s.db.Exec("UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id",
		sql.Named("id", task.Id),
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))

	return err
}

func (s Store) DeleteTask(id int) error {
	_, err := s.db.Exec("DELETE FROM scheduler WHERE id = :id", sql.Named("id", id))

	return err
}
