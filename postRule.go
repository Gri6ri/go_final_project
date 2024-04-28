package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

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

func (h Handler) postTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task Task
	var buf bytes.Buffer
	today := time.Now().Format(dateFormat)

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, wrappJsonError(err.Error()), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(w, wrappJsonError(err.Error()), http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		http.Error(w, wrappJsonError("не задано поле 'Title'"), http.StatusBadRequest)
		return
	}

	if task.Date == "" {
		task.Date = today
	}

	date, err := time.Parse(dateFormat, task.Date)
	if err != nil {
		http.Error(w, wrappJsonError(err.Error()), http.StatusBadRequest)
		return
	}

	var nextDate string
	if task.Repeat == "" {
		nextDate = today
	} else {
		nextDate, err = h.service.getNextDate(today, task.Date, task.Repeat)
		if err != nil {
			http.Error(w, wrappJsonError(err.Error()), http.StatusBadRequest)
			return
		}
	}

	if date.Before(time.Now().Truncate(24 * time.Hour)) {
		task.Date = nextDate
	}

	id, err := h.service.store.postTask(task)
	if err != nil {
		http.Error(w, wrappJsonError(err.Error()), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(map[string]any{"id": id})
	if err != nil {
		http.Error(w, wrappJsonError(err.Error()), http.StatusInternalServerError)
		return
	}

	h.writeResponse(w, http.StatusCreated, resp)
}
