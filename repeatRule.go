package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "20060102"

func (s Service) getNextDate(now string, date string, repeat string) (string, error) {
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
	if firstPart == "y" {
		nextDate = nextDate.AddDate(1, 0, 0)
		for nextDate.Before(parsedNow) || nextDate.Equal(parsedNow) {
			nextDate = nextDate.AddDate(1, 0, 0)
		}
		return nextDate.Format(dateFormat), nil
	} else if firstPart == "d" && len(repeatParts) == 2 {
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
	} else {
		return "", fmt.Errorf("ошибка: неподдерживаемый формат параметра 'repeat'")
	}
}

func (h Handler) getNextDateHandler(w http.ResponseWriter, r *http.Request) {
	now := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	nextDate, err := h.service.getNextDate(now, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.writeResponse(w, http.StatusOK, []byte(nextDate))
}

func (h Handler) writeResponse(w http.ResponseWriter, status int, body []byte) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)

	_, err := w.Write(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
