package send_response

import (
	"fmt"
	"time"
)

type Event struct {
	Id       int
	UserId   int64
	Name     string
	DateTime time.Time
	Weekly   bool
	Monthly  bool
	Yearly   bool
}

func WelcomeMessage() string {
	return "Hello! This is a reminder bot. Send \n/newevent to create an event " +
		"\n/myevents to get the list of your events " +
		"\n/edit to update your event, " +
		"\n/delete to delete one of your events, " +
		"\n/deleteall to delete all of your events."
}

func EmptyText() string {
	return "Write text blyedina"
}

func AskForName() string {
	return "Write the name of your event."
}

func AskForDate() string {
	return "Write the date of your event. Use the YYYY-MM-DD format."
}

func AskForTime() string {
	return "Write the time of your event. Use the HH:MM 24h format."
}

func AskHowFrequently() string {
	return "Write if I should remind you about this event weekly, monthly or yearly."
}

func MakeDateTimeField(Date, ETime string) time.Time {
	eventTime := Date + " " + ETime
	finTime, err := time.Parse("2006-01-02 15:04", eventTime)
	if err != nil {
		fmt.Println(err)
	}
	return finTime
}

func WhichEventDelete() string {
	return "Which event do you want to delete?"
}
