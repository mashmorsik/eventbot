package data_test

import (
	"database/sql"
	"eventbot/Logger"
	"eventbot/data"
	"math/rand"
	"reflect"
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
	defer func(db *data.Data, userId int64) {
		err := db.DeleteUser(userId)
		if err != nil {
			Logger.Sugar.Errorln("DeleteUser fail.")
		}
	}(db, int64(userId))

	if db.IsUser(int64(userId)) == false {
		err := db.AddUser(int64(userId))
		if err != nil {
			Logger.Sugar.Errorln("AddUser failed.")
		}
		if db.IsUser(int64(userId)) != true {
			t.Errorf("AddUser failed, got: %v, want: %v.", db.IsUser(int64(userId)), true)
		}
	} else {
		t.Error("Invalid test: User already exists.")
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

			defer func(db *data.Data, eventId int) {
				err = db.DeleteEvent(eventId)
				if err != nil {
					Logger.Sugar.Errorln("DeleteEvent failed.")
				}
			}(db, eventId)

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

func TestData_GetUserEvents(t *testing.T) {
	db := data.NewData(MustConnectTest())

	testCases := []struct {
		name         string
		userId       int64
		expectedType map[int]*data.Event
	}{
		{
			name:   "GetEvents of an existing user.",
			userId: int64(2080632788),
		},
		{
			name:         "GetEvents of a non-existing user.",
			userId:       0,
			expectedType: nil,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := db.GetUserEvents(tt.userId)
			if reflect.TypeOf(got) != reflect.TypeOf(tt.expectedType) {
				t.Error("Event type doesn't match.")
			}

			for _, e := range got {
				if e.UserId != tt.userId {
					t.Error("UserIds don't match.")
				}
			}
		})
	}
}

func TestData_UpdateEvent(t *testing.T) {
	db := data.NewData(MustConnectTest())
	id, err := db.CreateEvent(int64(2080632788), int64(2080632788), "ChristmasParty", time.Time{}, "35 11 16 01 *")
	if err != nil {
		Logger.Sugar.Errorln("Couldn't create event.")
	}
	var result bool

	defer func(db *data.Data, eventId int) {
		err = db.DeleteEvent(eventId)
		if err != nil {
			Logger.Sugar.Errorln("DeleteEvent in defer failed.")
		}
	}(db, id)

	type args struct {
		eventId   int
		eventName string
		timeDate  time.Time
		cron      string
	}
	testCases := []struct {
		name     string
		args     args
		expected bool
	}{
		{"Update existing event with valid data", args{id, "ChristmasParty", time.Now(), "once"}, true},
		{"Update existing event with invalid eventName", args{id, "", time.Now(), "once"}, false},
		{"Update non-existing event", args{0, "Basketball", time.Now(), "once"}, false},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			if db.UpdateEvent(tt.args.eventId, tt.args.eventName, tt.args.timeDate, tt.args.cron) == nil {
				result = true
			} else {
				result = false
			}
			if tt.expected != result {
				t.Error("UpdateEvent Failed.")
			}
			err, event := db.GetEvent(tt.args.eventId)
			if err != nil {
				Logger.Sugar.Errorln("GetEvent failed.")
			}
			if (event.Name == tt.args.eventName) && (event.Cron == tt.args.cron) {
				result = true
			}
			if tt.expected != result {
				t.Error("Event properties not updated.")
			}
		})
	}
}

func TestData_DeleteAllEvents(t *testing.T) {
	db := data.NewData(MustConnectTest())
	err := db.AddUser(2080632730)
	if err != nil {
		Logger.Sugar.Errorln("Couldn't add user.")
	}
	_, err = db.CreateEvent(int64(2080632730), int64(2080632730), "ChristmasParty", time.Time{}, "35 11 16 01 *")
	_, err = db.CreateEvent(int64(2080632730), int64(2080632730), "Holidays", time.Time{}, "20 12 14 02 *")
	_, err = db.CreateEvent(int64(2080632730), int64(2080632730), "Daily", time.Time{}, "10 10 11 04 *")
	if err != nil {
		Logger.Sugar.Errorln("Couldn't create event.")
	}

	var result bool

	tests := []struct {
		name     string
		userId   int64
		expected bool
	}{
		{"DeleteAllEvents for an existing user", 2080632730, true},
		{"DeleteAllEvents for a non-existing user", 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err = db.DeleteAllEvents(tt.userId)
			if err != nil {
				Logger.Sugar.Errorln("DeleteAllEvents failed.")
			}

			events, err := db.GetUserEvents(tt.userId)
			if err != nil {
				result = false
				Logger.Sugar.Errorln("GetUserEvents failed.")
			}

			if events == nil {
				result = true
			}

			if tt.expected != result {
				t.Error("DeleteAllEvents failed.")
			}
		})
	}
}

func TestData_DisabledTrue(t *testing.T) {
	db := data.NewData(MustConnectTest())
	id, err := db.CreateEvent(int64(2080632730), int64(2080632730), "Daily", time.Time{}, "10 10 11 04 *")
	if err != nil {
		Logger.Sugar.Errorln("CreateEvent failed.")
	}

	testCases := []struct {
		name     string
		eventId  int
		expected bool
	}{
		{"Disable existing event.", id, true},
		{"Disable non-existing event.", -1, false},
	}

	var result bool

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			err = db.DisabledTrue(tt.eventId)
			if err != nil {
				t.Error("DisabledTrue failed.")
			}
			// returns an event with non-existing id
			err, event := db.GetEvent(tt.eventId)
			if err != nil {
				Logger.Sugar.Errorln("GetEvent failed.")
			}
			if event == nil {
				result = false
			}
			if event.Disabled == true {
				result = true
			}
			if tt.expected != result {
				t.Error("DisabledTrue failed.")
			}
		})
	}
}

func TestData_DisabledFalse(t *testing.T) {
	db := data.NewData(MustConnectTest())
	id, err := db.CreateEvent(int64(2080632730), int64(2080632730), "Weekly", time.Time{}, "10 10 11 04 *")
	err = db.DisabledTrue(id)
	if err != nil {
		Logger.Sugar.Errorln("DisabledTrue failed.")
	}
	if err != nil {
		Logger.Sugar.Errorln("CreateEvent failed.")
	}

	testCases := []struct {
		name     string
		eventId  int
		expected bool
	}{
		{"Enable existing event.", id, true},
		{"Enable non-existing event.", -1, false},
	}

	var result bool

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			err = db.DisabledFalse(tt.eventId)
			if err != nil {
				t.Error("DisabledTrue failed.")
			}
			// returns an event with non-existing id
			err, event := db.GetEvent(tt.eventId)
			if err != nil {
				Logger.Sugar.Errorln("GetEvent failed.")
			}
			if event == nil {
				result = false
			}
			if event.Disabled == false {
				result = true
			}
			if tt.expected != result {
				t.Error("DisabledFalse failed.")
			}
		})
	}
}

func TestData_SetLastFired(t *testing.T) {
	db := data.NewData(MustConnectTest())
	id, err := db.CreateEvent(int64(2080632730), int64(2080632730), "Daily", time.Time{}, "10 10 11 04 *")
	if err != nil {
		Logger.Sugar.Errorln("CreateEvent failed.")
	}
	timeDateStr, err := time.Parse("2006-01-02 15:04", "2023-12-31 23:59")

	testCases := []struct {
		name     string
		eventId  int
		timeDate time.Time
		expected bool
	}{
		{"SetLastFired for an existing event valid data.", id, timeDateStr, true},
		{"SetLastFired for a non-existing event.", -1, timeDateStr, false},
	}

	var result bool

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			err = db.SetLastFired(tt.timeDate, tt.eventId)
			if err != nil {
				t.Error("SetLastFired failed.")
			}
			err, event := db.GetEvent(tt.eventId)
			if err != nil {
				Logger.Sugar.Errorln("GetEvent failed.")
			}
			if event == nil {
				result = false
			}
			if event.LastFired.Equal(tt.timeDate) {
				result = true
			}
			if tt.expected != result {
				t.Error("SetLastFired failed.")
			}
		})
	}
}
