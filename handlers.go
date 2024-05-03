package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type Handler struct {
	store Store
}

func NewHandler(store Store) Handler {
	return Handler{store: store}
}

func wrappJsonError(message string) string {
	s, _ := json.Marshal(map[string]any{"error": message})
	return string(s)
}

func (h Handler) writeResponse(w http.ResponseWriter, status int, body []byte) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)

	_, err := w.Write(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h Handler) getNextDateHandler(w http.ResponseWriter, r *http.Request) {
	now := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	nextDate, err := h.store.getNextDate(now, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.writeResponse(w, http.StatusOK, []byte(nextDate))
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
		nextDate, err = h.store.getNextDate(today, task.Date, task.Repeat)
		if err != nil {
			http.Error(w, wrappJsonError(err.Error()), http.StatusBadRequest)
			return
		}
	}

	if date.Before(time.Now().Truncate(24 * time.Hour)) {
		task.Date = nextDate
	}

	id, err := h.store.postTask(task)
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

func (h Handler) getAllTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.store.getAllTasks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if tasks == nil {
		tasks = make([]Task, 0)
	}

	resp, err := json.Marshal(map[string][]Task{"tasks": tasks})
	if err != nil {
		http.Error(w, wrappJsonError(err.Error()), http.StatusInternalServerError)
		return
	}
	h.writeResponse(w, http.StatusOK, resp)
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
