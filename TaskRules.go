package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

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
func (h Handler) getTaskHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.FormValue("id")

	if idStr == "" {
		http.Error(w, wrappJsonError("не задано поле 'id'"), http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, wrappJsonError(err.Error()), http.StatusBadRequest)
		return
	}

	task, err := h.store.getTask(id)
	if err != nil {
		http.Error(w, wrappJsonError(err.Error()), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(task)
	if err != nil {
		http.Error(w, wrappJsonError(err.Error()), http.StatusInternalServerError)
		return
	}
	h.writeResponse(w, http.StatusOK, resp)
}

func (h Handler) editTaskHandler(w http.ResponseWriter, r *http.Request) {
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

	if task.Id == "" {
		http.Error(w, wrappJsonError("не задано поле 'Id'"), http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(task.Id)
	if err != nil {
		http.Error(w, wrappJsonError(err.Error()), http.StatusBadRequest)
		return
	}

	if _, err = h.store.getTask(id); err != nil {
		http.Error(w, wrappJsonError(err.Error()), http.StatusInternalServerError)
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
		nextDate, err = h.store.getNextDate(today, task.Date, task.Repeat)
		if err != nil {
			http.Error(w, wrappJsonError(err.Error()), http.StatusBadRequest)
			return
		}
	}

	if date.Before(time.Now().Truncate(24 * time.Hour)) {
		task.Date = nextDate
	}

	if err = h.store.UpdateTask(task); err != nil {
		http.Error(w, wrappJsonError(err.Error()), http.StatusInternalServerError)
		return
	}

	h.writeResponse(w, http.StatusOK, []byte("{}"))
}

func (h Handler) postTaskDoneHandler(w http.ResponseWriter, r *http.Request) {
	today := time.Now().Format(dateFormat)
	idStr := r.FormValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, wrappJsonError(err.Error()), http.StatusBadRequest)
		return
	}

	task, err := h.store.getTask(id)
	if err != nil {
		http.Error(w, wrappJsonError(err.Error()), http.StatusInternalServerError)
		return
	}

	if task.Repeat == "" {
		err = h.store.DeleteTask(id)
		if err != nil {
			http.Error(w, wrappJsonError(err.Error()), http.StatusInternalServerError)
			return
		}
	} else {
		nextDate, err := h.store.getNextDate(today, task.Date, task.Repeat)
		if err != nil {
			http.Error(w, wrappJsonError(err.Error()), http.StatusBadRequest)
			return
		}

		task.Date = nextDate
		if err = h.store.UpdateTask(task); err != nil {
			http.Error(w, wrappJsonError(err.Error()), http.StatusInternalServerError)
			return
		}
	}
	h.writeResponse(w, http.StatusOK, []byte("{}"))
}

func (h Handler) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.FormValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, wrappJsonError(err.Error()), http.StatusBadRequest)
		return
	}
	err = h.store.DeleteTask(id)
	if err != nil {
		http.Error(w, wrappJsonError(err.Error()), http.StatusInternalServerError)
		return
	}
	h.writeResponse(w, http.StatusOK, []byte("{}"))
}
