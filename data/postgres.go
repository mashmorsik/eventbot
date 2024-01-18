package data

import (
	"database/sql"
	"errors"
	"eventbot/Logger"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"os"
	"time"
)

type Event struct {
	EventId   int
	UserId    int64
	Name      string
	ChatId    int64
	TimeDate  time.Time
	Cron      string
	LastFired time.Time
	Disabled  bool
}

type Data struct {
	db *sql.DB
}

func NewData(db *sql.DB) *Data {
	if db == nil {
		panic("db is nil")
	}
	return &Data{db: db}
}

func (r *Data) Db() *sql.DB {
	return r.db
}

func MustConnectPostgres() *sql.DB {
	connectionStr := "postgres://postgres:mysecretpassword@localhost:5432/reminder_bot?sslmode=disable&application_name=eventbot&connect_timeout=5"

	connection, err := sql.Open("postgres", connectionStr)
	if err != nil {
		panic(err)
	}

	if err = connection.Ping(); err != nil {
		Logger.Sugar.Panic(err)
	}

	return mustMigrate(connection)
}

func mustMigrate(connection *sql.DB) *sql.DB {
	driver, err := postgres.WithInstance(connection, &postgres.Config{})
	if err != nil {
		panic(err)
	}

	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	migrationPath := fmt.Sprintf("file://%s/migration", path)
	fmt.Printf("migrationPath : %s\n", migrationPath)

	m, err := migrate.NewWithDatabaseInstance(
		migrationPath,
		"reminder_bot", driver)

	if err = m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			Logger.Sugar.Infoln("no changes in migration, skip")

		} else {
			panic(err)
		}
	}

	return connection
}

func (r *Data) IsUser(userId int64) bool {
	var result int

	sqlGetUserId := `
	SELECT * FROM users
	WHERE user_id = ($1)`
	row := r.db.QueryRow(sqlGetUserId, userId)
	err := row.Scan(result)
	if errors.Is(err, sql.ErrNoRows) {
		return false
	} else {
		return true
	}
}

func (r *Data) AddUser(userId int64) error {
	if r.IsUser(userId) == false {
		sqlAddUser := `
	INSERT INTO users(user_id)
	VALUES($1) on conflict (user_id) do nothing`
		res, err := r.db.Exec(sqlAddUser, userId)
		if err != nil {
			panic(err)
		}
		ra, _ := res.RowsAffected()
		fmt.Printf("rows affected: %v", ra)
	}
	return nil
}

func (r *Data) DeleteUser(userId int64) error {
	if r.IsUser(userId) == true {
		sqlDeleteUser := `
	DELETE FROM users
	WHERE user_id = $1`
		res, err := r.db.Exec(sqlDeleteUser, userId)
		if err != nil {
			panic(err)
		}
		ra, _ := res.RowsAffected()
		fmt.Printf("rows affected: %v", ra)
	}
	return nil
}

func (r *Data) GetEvent(eventId int) (error, *Event) {
	var e Event

	sqlGetEvent := `
	SELECT * FROM events
	WHERE id = $1`

	rows, err := r.db.Query(sqlGetEvent, eventId)
	if err != nil {
		return err, nil
	}

	if rows.Next() {
		if err = rows.Scan(&e.EventId, &e.UserId, &e.ChatId, &e.Name, &e.TimeDate, &e.Cron, &e.LastFired, &e.Disabled); err != nil {
			return err, nil
		}
		return nil, &e
	}

	if err := rows.Err(); err != nil {
		return err, nil
	}

	// No rows found, return nil for the event
	return nil, nil
}

func (r *Data) CreateEvent(userId int64, chatId int64, name string, timeDate time.Time, cron string) (int, error) {
	var e Event

	sqlCreateEvent := `
		INSERT INTO events(user_id, chat_id, name, time_date, cron, last_fired)
		VALUES($1, $2, $3, $4, $5, $6)
		RETURNING id`

	rows, err := r.db.Query(sqlCreateEvent, userId, chatId, name, timeDate, cron, time.Time{})
	if err != nil {
		return 0, err
	}

	for rows.Next() {
		if err = rows.Scan(&e.EventId); err != nil {
			return 0, err
		}
		return e.EventId, nil
	}

	if err = rows.Err(); err != nil {
		return 0, err
	}

	return e.EventId, nil
}

func (r *Data) GetUserEvents(userId int64) (map[int]*Event, error) {
	var EventsList = make(map[int]*Event)

	sqlGetEventsList := `
	SELECT id, user_id, name, time_date, cron, last_fired, disabled FROM events
	WHERE user_id = $1`

	rows, err := r.db.Query(sqlGetEventsList, userId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var e Event

		if err = rows.Scan(&e.EventId, &e.UserId, &e.Name, &e.TimeDate, &e.Cron, &e.LastFired, &e.Disabled); err != nil {
			return nil, err
		}
		EventsList[e.EventId] = &e
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return EventsList, nil
}

func (r *Data) GetCronMultipleActive() (map[int]*Event, error) {
	sqlFindRemindEvent := `
	SELECT * FROM events
	WHERE cron != $1 and disabled = false`

	rows, err := r.db.Query(sqlFindRemindEvent)
	if err != nil {
		return nil, err
	}

	var currentEvents = make(map[int]*Event)

	for rows.Next() {
		var e Event

		if err = rows.Scan(&e.EventId, &e.UserId, &e.ChatId, &e.Name, &e.TimeDate, &e.Cron, &e.LastFired, &e.Disabled); err != nil {
			return nil, err
		}
		currentEvents[e.EventId] = &e
	}

	return currentEvents, nil
}

func (r *Data) GetOnceNotFired() (map[int]*Event, error) {
	sqlGetOnceNotFired := `
	SELECT * FROM events
	WHERE last_fired = $1 and cron = $2 and disabled = false`

	rows, err := r.db.Query(sqlGetOnceNotFired, time.Time{}, "once")
	if err != nil {
		return nil, err
	}

	var currentEvents = make(map[int]*Event)

	for rows.Next() {
		var e Event

		if err = rows.Scan(&e.EventId, &e.UserId, &e.ChatId, &e.Name, &e.TimeDate, &e.Cron, &e.LastFired, &e.Disabled); err != nil {
			return nil, err
		}
		currentEvents[e.EventId] = &e
	}

	return currentEvents, nil
}

func (r *Data) UpdateEvent(eventId int, name string, timeDate time.Time, cron string) error {
	sqlUpdateEvent := `
	UPDATE events
	SET name = $1, time_date = $2, cron = $3, last_fired = $4
	WHERE id = $5`

	_, err := r.db.Exec(sqlUpdateEvent, name, timeDate, cron, time.Time{}, eventId)
	if err != nil {
		Logger.Sugar.Errorln("UpdateEvent failed.")
		return err
	}
	return nil
}

func (r *Data) DeleteEvent(eventId int) error {
	sqlDeleteEvent := `
	DELETE FROM events
	WHERE id = $1`

	result, err := r.db.Exec(sqlDeleteEvent, eventId)
	if err != nil {
		Logger.Sugar.Errorln("DeleteEvent failed.")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		Logger.Sugar.Errorln("rowsAffected error.")
		return err
	}

	if rowsAffected == 0 {
		Logger.Sugar.Errorln("No events with this id found.")
		return errors.New("no events with this id found")
	} else {
		fmt.Printf("The SQL query affected %d rows.\n", rowsAffected)
	}

	return nil
}

func (r *Data) DeleteAllEvents(userId int64) error {
	sqlDeleteAllEvents := `
	DELETE FROM events
	WHERE user_id = $1`

	_, err := r.db.Exec(sqlDeleteAllEvents, userId)
	if err != nil {
		Logger.Sugar.Errorln("DeleteAllEvents failed.")
	}
	return nil
}

func (r *Data) DisabledTrue(eventId int) error {
	sqlDisabledTrue := `
	UPDATE events
	SET disabled = true
	WHERE id = $1`

	_, err := r.db.Exec(sqlDisabledTrue, eventId)
	if err != nil {
		Logger.Sugar.Panic("DisabledTrue failed")
	}
	return nil
}

func (r *Data) DisabledFalse(eventId int) error {
	sqlDisabledFalse := `
	UPDATE events
	SET disabled = false
	WHERE id = $1`

	_, err := r.db.Exec(sqlDisabledFalse, eventId)
	if err != nil {
		Logger.Sugar.Panic("DisabledFalse failed.")
	}
	return nil
}

func (r *Data) SetLastFired(lastFired time.Time, eventId int) error {
	sqlDisabledFalse := `
	UPDATE events
	SET last_fired = $1
	WHERE id = $2`

	_, err := r.db.Exec(sqlDisabledFalse, lastFired, eventId)
	if err != nil {
		Logger.Sugar.Panic("SetLastFired failed.")
	}

	return nil
}
