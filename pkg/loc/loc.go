package loc

import (
	"fmt"
	"time"
)

var CurrentLoc = MskLoc()

func MskLoc() *time.Location {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		panic(fmt.Sprintf("Error loading time zone:%v", err))
	}
	return loc
}
