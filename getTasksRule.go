package main

import (
	"encoding/json"
	"net/http"
)

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

func (h Handler) getAllTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.service.store.getAllTasks()
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
