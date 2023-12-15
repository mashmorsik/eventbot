package EditEvents

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"time"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "admin"
	dbname   = "reminder_bot"
)

func ConnectPostgres() *sql.DB {
	connectionStr := "user=postgres password=admin dbname=reminder_bot port=5432 sslmode=disable"

	connection, err := sql.Open("postgres", connectionStr)
	if err != nil {
		panic(err)
	}

	return connection
}

func IsUser(userId int64) bool {
	db := ConnectPostgres()

	sqlGetUserId := `
	SELECT * FROM users
	WHERE user_id = ($1)`
	row := db.QueryRow(sqlGetUserId, userId)
	//if row != nil {
	//	panic(row)
	//}
	if errors.Is(row.Err(), sql.ErrNoRows) {
		return false
	}
	return true
}

func AddUser(userId int64) {
	db := ConnectPostgres()

	if IsUser(userId) == false {
		sqlAddUser := `
	INSERT INTO users(user_id)
	VALUES($1)`
		res, err := db.Exec(sqlAddUser, userId)
		if err != nil {
			panic(err)
		}
		ra, _ := res.RowsAffected()
		liId, err := res.LastInsertId()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("rows affected: %v, last inserted id: %v", ra, liId)
	}
	fmt.Println("User already exists")
}

func CreateEvent(userId int64, text string, timeDate time.Time, weekly bool, monthly bool, yearly bool) {
	db := ConnectPostgres()

	sqlCreateEvent := `INSERT INTO events(user_id, name, time_date, weekly, monthly, yearly) 
		VALUES($1, $2, $3, $4, $5, $6)`
	_, err := db.Exec(sqlCreateEvent, 6537489202, "text", time.Now(), false, false, true)

	if err != nil {
		panic(err)
	}
}

func GetEventsList(userId int64) []string {
	db := ConnectPostgres()
	var EventsList []string

	sqlGetEventsList := `
	SELECT * FROM events 
	WHERE user_id = $1`
	row := db.QueryRow(sqlGetEventsList, userId)
	if row.Err() != nil {
		panic(row.Err())
	}

	err := row.Scan(EventsList)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return EventsList
		}

		panic(err)
	}
	return EventsList
}

// func UpdateEvent()
// func DeleteEvent()
