package main

import (
	"database/sql"
)

func (s Store) addTask(task Task) (int, error) {
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

// func (s Store) postTaskHandler(w http.ResponseWriter, r *http.Request) {
// 	var task Task
// 	var buf bytes.Buffer
// 	today := time.Now().Truncate(24 * time.Hour).Format(dateFormat)

// 	_, err := buf.ReadFrom(r.Body)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
// 	defer r.Body.Close()

// 	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	if task.Title == "" {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	dateTime, err := time.Parse(dateFormat, task.Date)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	if task.Date == "" {
// 		task.Date = today
// 	}

// 	if task.Repeat == "" {
// 		task.Date = today
// 	}
// 	nextDate, err := getNextDate(today, task.Date, task.Repeat)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
// 	if dateTime.Before(time.Now().Truncate(24 * time.Hour)) {
// 		task.Date = nextDate
// 	}

// 	// Добавление новой задачи в таблицу
// 	db, err := sql.Open("sqlite3", "scheduler.db")
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	defer db.Close()

// 	id, err := s.addTask(task)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	resp, err := json.Marshal(map[string]any{"id": id})
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

// 	_, err = w.Write([]byte(resp))
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}
// }
