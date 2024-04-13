package storage

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Stern-Ritter/go_task_manager/internal/errors"
	"github.com/Stern-Ritter/go_task_manager/internal/model"
	"github.com/Stern-Ritter/go_task_manager/internal/utils"
)

type TaskStore struct {
	db *sql.DB
}

func NewTaskStore(db *sql.DB) TaskStore {
	return TaskStore{db: db}
}

func (s TaskStore) Create(t model.Task) (int, error) {
	res, err := s.db.Exec(`
		INSERT INTO scheduler (date, title, comment, repeat) 
		VALUES (:date, :title, :comment, :repeat)
	`,
		sql.Named("date", t.Date.Format("20060102")),
		sql.Named("title", t.Title),
		sql.Named("comment", t.Comment),
		sql.Named("repeat", t.Repeat))

	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (s TaskStore) Update(t model.Task) error {
	res, err := s.db.Exec(`
		UPDATE scheduler 
		SET date = :date, title = :title, comment = :comment, repeat = :repeat 
		WHERE id = :id
	`,
		sql.Named("id", t.ID),
		sql.Named("date", t.Date.Format("20060102")),
		sql.Named("title", t.Title),
		sql.Named("comment", t.Comment),
		sql.Named("repeat", t.Repeat))

	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil || rows == 0 {
		return errors.NewTaskNotExists(fmt.Sprintf("Task with id: %d doesn`t exist", t.ID), err)
	}
	return nil
}

func (s TaskStore) Complete(t model.Task) error {
	now := time.Now()
	nextDate, err := utils.NextDate(now, t.Date.Format("20060102"), t.Repeat)
	if err != nil {
		return err
	}

	res, err := s.db.Exec(`
		UPDATE scheduler 
		SET date = :date 
		WHERE id = :id
	`,
		sql.Named("id", t.ID),
		sql.Named("date", nextDate))

	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil || rows == 0 {
		return errors.NewTaskNotExists(fmt.Sprintf("Task with id: %d doesn`t exist", t.ID), err)
	}
	return nil
}

func (s TaskStore) Delete(id int) error {
	res, err := s.db.Exec(`
		DELETE FROM scheduler 
		WHERE id = :id
	`,
		sql.Named("id", id))

	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil || rows == 0 {
		return errors.NewTaskNotExists(fmt.Sprintf("Task with id: %d doesn`t exist", id), err)
	}
	return nil
}

func (s TaskStore) GetByID(id int) (model.Task, error) {
	row := s.db.QueryRow(`
		SELECT id, date, title, comment, repeat 
		FROM scheduler 
		WHERE id = :id
	`,
		sql.Named("id", id))

	t := model.Task{}
	var date string
	err := row.Scan(&t.ID, &date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		return t, err
	}
	t.Date, err = time.Parse("20060102", date)
	if err != nil {
		return t, err
	}

	return t, nil
}

func (s TaskStore) GetAll() ([]model.Task, error) {
	rows, err := s.db.Query(`
		SELECT id, date, title, comment, repeat 
		FROM scheduler 
		ORDER BY date
	`)

	var res []model.Task

	if err != nil {
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		t := model.Task{}
		var date string
		err := rows.Scan(&t.ID, &date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return res, err
		}
		t.Date, err = time.Parse("20060102", date)
		if err != nil {
			return res, err
		}
		res = append(res, t)
	}

	err = rows.Err()
	return res, err
}

func (s TaskStore) GetAllByTitleOrComment(search string) ([]model.Task, error) {
	rows, err := s.db.Query(`
		SELECT id, date, title, comment, repeat 
		FROM scheduler 
		WHERE title LIKE :search OR 
		comment LIKE :search ORDER BY date
	`,
		sql.Named("search", search))

	var res []model.Task

	if err != nil {
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		t := model.Task{}
		var date string
		err := rows.Scan(&t.ID, &date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return res, err
		}
		t.Date, err = time.Parse("20060102", date)
		if err != nil {
			return res, err
		}
		res = append(res, t)
	}

	err = rows.Err()
	return res, err
}

func (s TaskStore) GetAllByDate(date string) ([]model.Task, error) {
	rows, err := s.db.Query(`
		SELECT id, date, title, comment, repeat 
		FROM scheduler 
		WHERE date = :date
	`,
		sql.Named("date", date))

	var res []model.Task

	if err != nil {
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		t := model.Task{}
		var date string
		err := rows.Scan(&t.ID, &date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return res, err
		}
		t.Date, err = time.Parse("20060102", date)
		if err != nil {
			return res, err
		}
		res = append(res, t)
	}

	err = rows.Err()
	return res, err
}
