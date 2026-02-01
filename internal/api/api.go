package api

import (
	"net/http"
	"time"
)

func Init() {
	http.HandleFunc("/api/nextdate", nextDateHandler)
	http.HandleFunc("/api/task", taskHandler)
	http.HandleFunc("/api/tasks", tasksHandler)
	http.HandleFunc("/api/task/done", doneTaskHandler)
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)

	case http.MethodGet:
		getTaskByIDHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)

	case http.MethodDelete:
		deleteTaskHandler(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.URL.Query().Get("now")
	dateStr := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	_, err := time.Parse("20060102", dateStr)
	if err != nil {
		http.Error(w, "Invalid date format, expected YYYYMMDD", http.StatusBadRequest)
		return
	}
	var nowDate time.Time
	if nowStr == "" {
		nowDate = time.Now()
	} else {
		nowDate, err = time.Parse("20060102", nowStr)
		if err != nil {
			http.Error(w, "Invalid now format, expected YYYYMMDD", http.StatusBadRequest)
			return
		}
	}
	if repeat == "" {
		http.Error(w, "Repeat parametr is required", http.StatusBadRequest)
		return
	}

	result, err := NextDate(nowDate, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(result))
}
