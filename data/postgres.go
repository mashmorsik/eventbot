package data

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"time"
)

type Data struct {
	db *sql.DB
}

func NewData(db *sql.DB) *Data {
	if db == nil {
		panic("db is nil")
	}
	return &Data{db: db}
}

func MustConnectPostgres() *sql.DB {
	connectionStr := "postgres://postgres:admin@localhost:5432/reminder_bot?sslmode=disable&application_name=eventbot&connect_timeout=5"

	connection, err := sql.Open("postgres", connectionStr)
	if err != nil {
		panic(err)
	}

	if err = connection.Ping(); err != nil {
		panic(err)
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

func (r *Data) CreateEvent(userId int64, name string, timeDate time.Time, weekly bool, monthly bool, yearly bool) {

	sqlCreateEvent := `INSERT INTO events(user_id, name, time_date, weekly, monthly, yearly)
		VALUES($1, $2, $3, $4, $5, $6)`
	_, err := r.db.Exec(sqlCreateEvent, userId, name, timeDate, weekly, monthly, yearly)
	if err != nil {
		panic(err)
	}
}

func (r *Data) GetEventsList(userId int64) (map[int]string, error) {
	//var EventsList []string
	var EventsList = make(map[int]string)

	sqlGetEventsList := `
	SELECT id, name FROM events
	WHERE user_id = $1`

	rows, err := r.db.Query(sqlGetEventsList, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var eventName string
		var id int
		if err := rows.Scan(&id, &eventName); err != nil {
			return nil, err
		}
		EventsList[id] = eventName
	}

	//sqlGetEventsList := `
	//SELECT name FROM events
	//WHERE user_id = $1`

	//rows, err := r.db.Query(sqlGetEventsList, userId)
	//if err != nil {
	//	return nil, err
	//}
	//defer rows.Close()
	//
	//for rows.Next() {
	//	var eventName string
	//	if err := rows.Scan(&eventName); err != nil {
	//		return nil, err
	//	}
	//	EventsList = append(EventsList, eventName)
	//}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	fmt.Println(EventsList)
	return EventsList, nil
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
