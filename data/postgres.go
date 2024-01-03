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

type Data struct {
	db *sql.DB
}

type Event struct {
	EventId  int
	UserId   int64
	Name     string
	ChatId   int64
	TimeDate string
	Cron     string
}

func NewData(db *sql.DB) *Data {
	if db == nil {
		panic("db is nil")
	}
	return &Data{db: db}
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

// FIXME add result erorr
func (r *Data) AddUser(userId int64) {
	if r.IsUser(userId) == false {
		sqlAddUser := `
	INSERT INTO users(user_id)
	VALUES($1)`
		res, err := r.db.Exec(sqlAddUser, userId)
		if err != nil {
			panic(err)
		}
		ra, _ := res.RowsAffected()
		//liId, err := res.LastInsertId()
		//if err != nil {
		//	fmt.Println(err)
		//}
		fmt.Printf("rows affected: %v", ra)
	} else {
		fmt.Println("BotUser already exists")
	}
}

func (r *Data) GetUsersList() ([]int64, error) {
	var UsersList []int64

	sqlGetUsersList := `
	SELECT * FROM users`

	rows, err := r.db.Query(sqlGetUsersList)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var UserId int64
		if err := rows.Scan(&UserId); err != nil {
			return nil, err
		}
		UsersList = append(UsersList, UserId)
	}
	fmt.Println(UsersList)
	return UsersList, nil
}

func (r *Data) DeleteUser(userId int64) {
	sqlDeleteEvent := `
	DELETE FROM users
	WHERE user_id = $1`

	_, err := r.db.Exec(sqlDeleteEvent, userId)
	if err != nil {
		panic(err)
	}
}

func (r *Data) CreateEvent(userId int64, chatId int64, name string, timeDate time.Time, cron string) {

	sqlCreateEvent := `INSERT INTO events(user_id, chat_id, name, time_date, cron)
		VALUES($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(sqlCreateEvent, userId, chatId, name, timeDate, cron)
	if err != nil {
		panic(err)
	}
}

func (r *Data) GetEventsList(userId int64) (map[int]string, error) {
	var EventsList = make(map[int]string)

	sqlGetEventsList := `
	SELECT id, name FROM events
	WHERE user_id = $1`

	rows, err := r.db.Query(sqlGetEventsList, userId)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		if err := rows.Close(); err != nil {
			panic(err)
			return
		}
	}(rows)

	for rows.Next() {
		// FIXME: marshall into structure - events
		var eventName string
		var id int
		if err := rows.Scan(&id, &eventName); err != nil {
			return nil, err
		}
		EventsList[id] = eventName
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	fmt.Println(EventsList)
	return EventsList, nil
}

func (r *Data) GetAllEvents() (map[int]*Event, error) {
	sqlFindRemindEvent := `
	SELECT * FROM events`

	rows, err := r.db.Query(sqlFindRemindEvent)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		if err := rows.Close(); err != nil {
			panic(err)
			return
		}
	}(rows)

	var currentEvents = make(map[int]*Event)
	var (
		id       int
		userId   int64
		chatId   int64
		name     string
		timeDate string
		cron     string
	)

	for rows.Next() {
		if err := rows.Scan(&id, &userId, &chatId, &name, &timeDate, &cron); err != nil {
			return nil, err
		}
		currentEvents[id] = &Event{
			Name:     name,
			ChatId:   chatId,
			UserId:   userId,
			TimeDate: timeDate,
			Cron:     cron,
		}
	}

	return currentEvents, nil
}

func (r *Data) UpdateEvent(eventId int, name string, timeDate time.Time, weekly bool, monthly bool, yearly bool) {
	sqlUpdateEvent := `
	UPDATE events
	SET name = $1, time_date = $2, weekly = $3, monthly = $4, yearly = $5
	WHERE id = $6`

	_, err := r.db.Exec(sqlUpdateEvent, name, timeDate, weekly, monthly, yearly, eventId)
	if err != nil {
		panic(err)
	}
}

// FIXME: return error
func (r *Data) DeleteEvent(eventId int) {
	sqlDeleteEvent := `
	DELETE FROM events
	WHERE id = $1`

	_, err := r.db.Exec(sqlDeleteEvent, eventId)
	if err != nil {
		panic(err)
	}
}

func (r *Data) DeleteAllEvents(userId int64) {
	sqlDeleteAllEvents := `
	DELETE FROM events
	WHERE user_id = $1`

	_, err := r.db.Exec(sqlDeleteAllEvents, userId)
	if err != nil {
		panic(err)
	}
}

func (r *Data) FindRemindEvent() (map[int]*Event, error) {
	sqlFindRemindEvent := `
	SELECT id, chat_id, name FROM events
	WHERE time_date = $1`

	loc, _ := time.LoadLocation("Asia/Shanghai")
	currentTime := time.Now().In(loc)
	rows, err := r.db.Query(sqlFindRemindEvent, currentTime)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		if err := rows.Close(); err != nil {
			panic(err)
			return
		}
	}(rows)

	var currentEvents = make(map[int]*Event)
	var id int
	var chatId int64
	var name string

	for rows.Next() {
		if err := rows.Scan(&id, &chatId, &name); err != nil {
			return nil, err
		}
		currentEvents[id] = &Event{
			Name:   name,
			ChatId: chatId,
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	fmt.Println("This is the request")
	fmt.Println(currentEvents)
	return currentEvents, nil
}
