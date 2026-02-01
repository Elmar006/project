package db

import (
	"errors"
	"fmt"
	"time"

	logger "github.com/Elmar006/project/internal/logger"
)

type Task struct {
	ID      int    `json:"id,string"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	if task.Title == "" {
		return 0, errors.New("title is required")
	}
	if task.Date == "" {
		return 0, errors.New("date is required")
	}

	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func GetTasks(limit int) ([]*Task, error) {
	if limit > 50 {
		logger.L().Info("The limit cannot exceed 50 tasks.")
		limit = 50
	}

	rows, err := DB.Query(
		`SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?`,
		limit,
	)
	if err != nil {
		logger.L().Errorf("Failed: %v", err)
		return nil, err
	}

	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		task := &Task{}

		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			logger.L().Errorf("Failed: %v", err)
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if rows.Err() != nil {
		logger.L().Error("Scan error")
		return nil, rows.Err()
	}

	return tasks, nil
}

func SearchTask(search string, limit int) ([]*Task, error) {
	if limit > 50 {
		logger.L().Info("The limit cannot be more than 50")
		limit = 50
	}
	if t, err := time.Parse("02.01.2006", search); err == nil {
		date := t.Format("20060102")
		tasks, err := searchDate(date, limit)
		if err != nil {
			logger.L().Errorf("Date search error: %v", err)
		}

		return tasks, nil
	} else {
		tasks, err := searchTitleAndComment(search, limit)
		if err != nil {
			logger.L().Errorf("Error search title/comment %v", err)
			return nil, err
		}

		return tasks, nil
	}
}

func searchDate(date string, limit int) ([]*Task, error) {
	rows, err := DB.Query(
		`SELECT * FROM scheduler WHERE date = ? LIMIT ?`,
		date, limit,
	)
	if err != nil {
		logger.L().Errorf("date search error in DB: %v", err)
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			logger.L().Errorf("Scan error: %v", err)
			return nil, err
		}

		tasks = append(tasks, &task)
	}

	if rows.Err() != nil {
		logger.L().Error(err)
		return nil, err
	}

	return tasks, nil
}

func searchTitleAndComment(search string, limit int) ([]*Task, error) {
	parseSearch := "%" + search + "%"
	rows, err := DB.Query(
		`SELECT * FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?`,
		parseSearch, parseSearch, limit,
	)
	if err != nil {
		logger.L().Errorf("Comment/Title search error in DB: %v", err)
		return nil, err
	}

	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			logger.L().Errorf("Scan error: %v", err)
			return nil, err
		}

		tasks = append(tasks, &task)
	}

	if rows.Err() != nil {
		logger.L().Error(err)
		return nil, err
	}

	return tasks, nil
}

func GetTask(id int) (*Task, error) {
	task := &Task{}
	err := DB.QueryRow(
		`SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`,
		id,
	).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		logger.L().Error(err)
		return nil, err
	}

	return task, nil
}

func UpdateTask(task *Task) error {
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		logger.L().Errorf("DB Exec error: %v", err)
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		logger.L().Error(err)
		return err
	}

	if count == 0 {
		logger.L().Errorf("Incorrect id for updating task: ID=%d not found", task.ID)
		return fmt.Errorf("Failed: count = 0")
	}

	return nil
}

func DeleteTask(id int) error {
	query := `DELETE FROM scheduler WHERE id = ?`
	_, err := DB.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}
