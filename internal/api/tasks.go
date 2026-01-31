package api

import (
	"encoding/json"
	"net/http"

	"github.com/Elmar006/project/internal/db"
	"github.com/Elmar006/project/internal/logger"
)

type TaskResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	if search == "" {

		tasks, err := db.GetTasks(50)
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

		tasks, err := db.SearchTask(search, 50)
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
