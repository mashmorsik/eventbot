package data_test

import (
	"database/sql"
	"eventbot/Logger"
	"eventbot/data"
	"math/rand"
	"testing"
	"time"
)

func MustConnectTest() *sql.DB {
	connectionStr := "postgres://postgres:mysecretpassword@localhost:5435/reminder_test?sslmode=disable&application_name=eventbot&connect_timeout=5"

	connection, err := sql.Open("postgres", connectionStr)
	if err != nil {
		panic(err)
	}

	if err = connection.Ping(); err != nil {
		Logger.Sugar.Panic(err)
	}

	return connection
}

func TestData_AddUserNewUser(t *testing.T) {
	db := data.NewData(MustConnectTest())

	userId := rand.Intn(10000000000)

	if db.IsUser(int64(userId)) == false {
		err := db.AddUser(int64(userId))
		if err != nil {
			Logger.Sugar.Errorln("AddUser failed.")
		}
		result := db.IsUser(int64(userId))
		if result != true {
			t.Errorf("AddUser failed, got: %v, want: %v.\n", result, true)
		}
	} else {
		t.Error("Invalid test: User already exists.\n")
	}
}

func TestData_AddUserExistingUser(t *testing.T) {
	db := data.NewData(MustConnectTest())

	userId := 2080632788

	if db.IsUser(int64(userId)) == true {
		result := db.AddUser(int64(userId))
		if result != nil {
			t.Error("Add user failed.")
		}
	} else {
		t.Error("Invalid test: User doesn't exist.")
	}
}

func TestData_CreateEvent(t *testing.T) {
	db := data.NewData(MustConnectTest())

	testCases := []struct {
		name      string
		userId    int64
		chatId    int64
		eventName string
		timeDate  time.Time
		cron      string
		expected  bool
	}{
		{"Create event with valid fields", int64(2080632788), int64(2080632788), "Swimming",
			time.Now(), "once", true},
		{"Create event with invalid userId", int64(1111111111), int64(2080632788), "Dentist",
			time.Now(), "once", false},
		{"Create event with invalid eventName", int64(2080632788), int64(2080632788), "",
			time.Now(), "once", false},
		{"Create event with invalid cron", int64(2080632788), int64(2080632788), "Dancing",
			time.Now(), "", false},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			var result bool
			eventId, err := db.CreateEvent(tt.userId, tt.chatId, tt.eventName, tt.timeDate, tt.cron)
			if err != nil {
				result = false
			}
			events, err := db.GetUserEvents(tt.userId)
			if err != nil {
				t.Error("Couldn't get events list.")
			}
			for _, e := range events {
				if e.EventId == eventId && e.Name == tt.eventName && e.Cron == tt.cron {
					result = true
				}
			}
			if result != tt.expected {
				t.Error("CreateEvent failed.")
			}
		})
	}
}

func TestData_DeleteEvent(t *testing.T) {
	db := data.NewData(MustConnectTest())
	id, err := db.CreateEvent(int64(2080632788), int64(2080632788), "Theatre", time.Now(), "once")
	if err != nil {
		Logger.Sugar.Errorln("Couldn't create event.")
	}

	testCases := []struct {
		name     string
		eventId  int
		expected bool
	}{
		{"Delete an existing event.", id, true},
		{"Delete a non-existing event.", 0, false},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			var result bool
			if db.DeleteEvent(tt.eventId) == nil {
				result = true
			}
			if result != tt.expected {
				t.Error("DeleteEvent failed.")
			}
		})
	}
}
