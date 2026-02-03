package api

import (
	"net/http"
	"time"
)

func Init() {
	http.HandleFunc("/api/signin", signinHandler)     // Аутентификация - проверка пароля и генерация JWT токена
	http.HandleFunc("/api/nextdate", nextDateHandler) // Расчет следующей даты по правилу повторения

	http.HandleFunc("/api/task", auth(taskHandler))          // CRUD операции с задачей (создание, чтение, обновление, удаление)
	http.HandleFunc("/api/tasks", auth(tasksHandler))        // Получение списка задач и поиск
	http.HandleFunc("/api/task/done", auth(doneTaskHandler)) // Отметка задачи как выполненной с учетом правила повторения
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r) // Создание задачи
	case http.MethodGet:
		getTaskByIDHandler(w, r) // Получение задачи по ID
	case http.MethodPut:
		updateTaskHandler(w, r) // Обновление задачи
	case http.MethodDelete:
		deleteTaskHandler(w, r) // Удаление задачи
	default:
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.URL.Query().Get("now")
	dateStr := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	_, err := time.Parse("20060102", dateStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid date format, expected YYYYMMDD")
		return
	}
	var nowDate time.Time
	if nowStr == "" {
		nowDate = time.Now()
	} else {
		nowDate, err = time.Parse("20060102", nowStr)
		if err != nil {
			writeError(w, http.StatusBadRequest, "Invalid now format, expected YYYYMMDD")
			return
		}
	}
	if repeat == "" {
		writeError(w, http.StatusBadRequest, "Repeat parameter is required")
		return
	}

	result, err := NextDate(nowDate, dateStr, repeat)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(result))
}
