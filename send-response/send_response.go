package send_response

import (
	"eventbot/Logger"
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
		"\n/deleteall to delete all of your events, " +
		"\n/disable to disable one of your events, " +
		"\n/enable to enable one of your disabled events."
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
	return "Write if I should remind you about this event once, daily, weekly, monthly or yearly. "
}

func MakeDateTimeField(Date, ETime string) time.Time {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		Logger.Sugar.Errorln(err)
		return time.Time{}
	}

	eventTime := Date + " " + ETime
	finTime, err := time.ParseInLocation("2006-01-02 15:04", eventTime, loc)
	if err != nil {
		Logger.Sugar.Errorln(err)
		return time.Time{}
	}
	return finTime
}

func StringToCron(date, timeStr, frequency string) string {
	var cron string
	switch frequency {
	case "once":
		cron = "once"
	case "daily":
		cron = timeStr[3:] + " " + timeStr[0:2] + " * * *"
	case "weekly":
		weekday, err := time.Parse("2006-01-02", date)
		if err != nil {
			panic(err)
		}
		cron = timeStr[3:] + " " + timeStr[0:2] + " * * " + weekday.Weekday().String()[0:3]
	case "monthly":
		cron = timeStr[3:] + " " + timeStr[0:2] + " " + date[8:] + " * *"
	case "yearly":
		cron = timeStr[3:] + " " + timeStr[0:2] + " " + date[8:] + " " + date[5:7] + " *"
	}
	return cron
}

func WhichEventDelete() string {
	return "Which event do you want to delete?"
}
