package data

import (
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"reflect"
	"testing"
	"time"
)

func TestData_AddUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("can't create mock: %s", err)
	}
	defer db.Close()

	data := NewData(db)
	userID := 75348957237

	// ok query
	mock.
		ExpectExec(`INSERT INTO users`).
		WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = data.AddUser(int64(userID))
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestData_GetUserEvents(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("can't create mock: %s", err)
	}
	defer db.Close()

	data := NewData(db)
	userID := 5694653497

	rows := sqlmock.NewRows([]string{"id", "user_id", "name", "time_date", "cron", "last_fired", "disabled"})
	var e = &Event{
		EventId:   0,
		UserId:    0,
		Name:      "",
		ChatId:    0,
		TimeDate:  time.Time{},
		Cron:      "",
		LastFired: time.Time{},
		Disabled:  false,
	}
	expected := map[int]*Event{}
	expected[e.EventId] = e

	for _, item := range expected {
		rows = rows.AddRow(item.EventId, item.UserId, item.Name, item.TimeDate, item.Cron, item.LastFired, item.Disabled)
	}
	mock.
		ExpectQuery(`SELECT id, user_id, name, time_date, cron, last_fired, disabled FROM events WHERE user_id = $1`).
		WithArgs(userID).
		WillReturnRows(rows)

	items, err := data.GetUserEvents(int64(userID))
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if !reflect.DeepEqual(expected, items) {
		t.Errorf("results not match")
		return
	}
}
