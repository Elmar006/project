package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Elmar006/project/internal/db"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	ID int64 `json:"id"`
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

func writeSuccess(w http.ResponseWriter, id int64) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(SuccessResponse{ID: id})
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}
	defer r.Body.Close()

	if task.Title == "" {
		writeError(w, http.StatusBadRequest, "The 'Title' field cannot be empty")
		return
	}
	if task.Date == "" {
		task.Date = time.Now().Format("20060102")
	}
	if _, err := time.Parse("20060102", task.Date); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid date format")
		return
	}

	now := time.Now()
	t, _ := time.Parse("20060102", task.Date)
	nowNorm := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	tNorm := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())

	if tNorm.Before(nowNorm) {
		if task.Repeat != "" {
			next, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				writeError(w, http.StatusBadRequest, err.Error())
				return
			}
			task.Date = next
		} else {
			task.Date = now.Format("20060102")
		}
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Database error: "+err.Error())
		return
	}

	writeSuccess(w, id)
}
