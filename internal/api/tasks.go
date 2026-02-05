package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Elmar006/project/internal/db"
	"github.com/Elmar006/project/internal/logger"
)

type TaskResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	if search == "" {

		tasks, err := db.GetTasks(db.LimitCount)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if tasks == nil {
			tasks = []*db.Task{}
		}

		resp := TaskResp{Tasks: tasks}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)

	} else {

		tasks, err := db.SearchTask(search, db.LimitCount)
		if err != nil {
			logger.L().Errorf("Error when searching through Query data: %v", err)
			http.Error(w, "Error when searching through Query data", http.StatusInternalServerError)
			return
		}
		if tasks == nil {
			tasks = []*db.Task{}
		}

		resp := TaskResp{Tasks: tasks}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}

}

func getTaskByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.L().Errorf("Error: couldn't convert string to integer: %v", err)
		writeError(w, http.StatusBadRequest, "Error: couldn't convert string to integer")
		return
	}

	taskID, err := db.GetTask(id)
	if err != nil {
		logger.L().Errorf("Error: couldn't convert string to integer: %v", err)
		writeError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(taskID)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}
	defer r.Body.Close()

	if task.ID == 0 {
		writeError(w, http.StatusBadRequest, "Task ID required")
		return
	}
	if task.Title == "" {
		writeError(w, http.StatusBadRequest, "The 'Title' field cannot be empty")
		return
	}
	if task.Date == "" {
		task.Date = time.Now().Format(db.TimeDateFormat)
	}
	if _, err := time.Parse(db.TimeDateFormat, task.Date); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid date format")
		return
	}

	now := time.Now()
	t, _ := time.Parse(db.TimeDateFormat, task.Date)
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
			task.Date = now.Format(db.TimeDateFormat)
		}
	}

	err := db.UpdateTask(&task)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Database error: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{})
}

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		writeError(w, http.StatusBadRequest, "Task ID required")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.L().Errorf("Error: couldn't convert string to integer: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		logger.L().Errorf("Error: couldn't convert string to integer: %v", err)
		writeError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	if task.Repeat == "" {
		err = db.DeleteTask(id)
		if err != nil {
			logger.L().Errorf("Failed to delete task: %v", err)
			writeError(w, http.StatusInternalServerError, "Failed to delete task::"+err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{})
		return
	}
	now := time.Now()
	next, err := NextDate(now, task.Date, task.Repeat)
	if err != nil {
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
	}

	task.Date = next
	err = db.UpdateTask(task)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Database error: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{})
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		writeError(w, http.StatusBadRequest, "Task ID required")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.L().Errorf("Error: couldn't convert string to integer: %v", err)
		writeError(w, http.StatusBadRequest, "Invalid task ID format")
		return
	}

	err = db.DeleteTask(id)
	if err != nil {
		logger.L().Errorf("Failed to delete task: %v", err)
		writeError(w, http.StatusInternalServerError, "Failed to delete task::"+err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{})
}
