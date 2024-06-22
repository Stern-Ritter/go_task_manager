package model

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/Stern-Ritter/go_task_manager/internal/errors"
	"github.com/Stern-Ritter/go_task_manager/internal/utils"
)

type Task struct {
	ID      int
	Date    time.Time
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func (t *Task) UnmarshalJSON(data []byte) error {
	type TaskAlias Task

	aliasTask := &struct {
		*TaskAlias
		Date string `json:"date"`
		ID   string `json:"id"`
	}{
		TaskAlias: (*TaskAlias)(t),
	}

	if err := json.Unmarshal(data, aliasTask); err != nil {
		return err
	}

	if len(strings.TrimSpace(aliasTask.ID)) != 0 {
		id, err := strconv.Atoi(aliasTask.ID)
		if err != nil {
			return err
		}
		t.ID = id
	}

	if len(strings.TrimSpace(aliasTask.Title)) == 0 {
		return errors.NewInvalidTitleFormat("task title is empty", nil)
	}

	currDateTime := time.Now()
	now := time.Date(currDateTime.Year(), currDateTime.Month(), currDateTime.Day(), 0, 0, 0, 0, time.UTC)
	var date time.Time
	nextDate := ""

	if len(strings.TrimSpace(aliasTask.Date)) != 0 {
		value, err := time.Parse("20060102", aliasTask.Date)
		if err != nil {
			return errors.NewInvalidDateFormat("invalid task date format", err)
		}
		date = value
	} else {
		date = now
	}

	if len(strings.TrimSpace(aliasTask.Repeat)) != 0 {
		value, err := utils.NextDate(now, date.Format("20060102"), aliasTask.Repeat)
		if err != nil {
			return err
		}
		nextDate = value
	}

	if date.Before(now) {
		if nextDate == "" {
			date = now
		} else {
			value, err := time.Parse("20060102", nextDate)
			if err != nil {
				return errors.NewInvalidDateFormat("invalid task next date format", err)
			}
			date = value
		}
	}

	t.Date = date
	return nil
}
