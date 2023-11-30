package EditEvents

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
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
	//defer connection.Close()
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
		row := db.QueryRow(sqlAddUser, userId)
		if row != nil {
			panic(row)
		}
	}
	fmt.Println("User already exists")
}

func CreateEventDB(userId int64, text string, date string, time string, weekly bool, monthly bool, yearly bool) {
	db := ConnectPostgres()

	sqlCreateEvent := `
	INSERT INTO events(user_id, text, date, time, weekly, monthly, yearly)
	VALUES($1, $2, $3, $4, $5, $6, $7)`
	row := db.QueryRow(sqlCreateEvent, userId, text, date, time, weekly, monthly, yearly)
	if row != nil {
		panic(row)
	}
}

// func UpdateEvent()
// func DeleteEvent()
