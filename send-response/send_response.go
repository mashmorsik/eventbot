package send_response

import (
	"eventbot/data"
	"time"
)

var commands = []string{"/newevent", "/myevents", "/edit", "/delete", "/deleteall"}

type Response struct {
	data.BotUser
}

type Event struct {
	Id       int
	UserId   int64
	Name     string
	DateTime time.Time
	Weekly   bool
	Monthly  bool
	Yearly   bool
}

func SendResponse(text string) string {
	switch text {
	case commands[0]:
		return GetEvents()
	case commands[1]:
		return GetEvents()
	case commands[2]:
		return EditEvent()
	case commands[3]:
		return DeleteEvent()
	case commands[4]:
		return DeleteAllEvents()
	}
	return WelcomeMessage()
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
	return "Write if I should remind you about this event weekly, monthly or yearly. If you want to be reminded just once, " +
		"write: Once."
}

//func MakeDateTimeField(date, time string) time.Time {
//	eventTime := date + " " + time
//	finTime := time.Parse("2006-01-02 15:04", eventTime)
//	return finTime
//}

//func MakeFrequencyField(msg string) {
//	switch msg {
//	case "Weekly":
//		Event{Weekly: true}
//		Event{Monthly: false}
//		Event{Yearly: false}
//	case "Monthly":
//		Event{Weekly: false}
//		Event{Monthly: true}
//		Event{Yearly: false}
//	case "Yearly":
//		Event{Weekly: false}
//		Event{Monthly: false}
//		Event{Yearly: true}
//	case "Once":
//		Event{Weekly: false}
//		Event{Monthly: false}
//		Event{Yearly: false}
//	}
//}

func (r *Response) CreateEvent() string {
	//e := Event{
	//	Id:       0,
	//	UserId:   0,
	//	Name:     r.BotUser.Message,
	//	DateTime: time.Time{},
	//	Weekly:   false,
	//	Monthly:  false,
	//	Yearly:   false,
	//}
	//
	//e.Name = r.BotUser.Message
	//AskForName()
	//r.ReadMessage()
	//
	//e{Name: Response{BotUser.Message}}
	//AskForDate()
	//r.ReadMessage()
	//date := Response.BotUser.Message
	//AskForTime()
	//r.ReadMessage()
	//time := Response.BotUser.Message
	//Event{DateTime: MakeDateTimeField(date, time)}
	//AskHowFrequently()
	//frequency := Response.BotUser.Message
	//MakeFrequencyField(frequency)
	//Event{UserId: Response{BotUser.UserId}}
	//db := data.NewData(data.MustConnectPostgres())
	//db.CreateEvent(Event.UserId, Event.Name, Event.DateTime, Event.Weekly, Event.Monthly, Event.Yearly)
	//message := "Your event: " + Event.Name + " has been successfully created"
	return "message"
}

func GetEvents() string {
	//r.ReadMessage()
	//db := data.NewData(data.MustConnectPostgres())
	//list := db.GetEventsList(Event.UserId)
	//var msg string
	//
	//for _, event := range list {
	//	item := "\n/" + event
	//	msg += item
	//}

	return "msg"
}

func EditEvent() string {
	//r.ReadMessage()
	//db := data.NewData(data.MustConnectPostgres())
	//
	//AskForName()
	//r.ReadMessage()
	//Event{Name: Response{BotUser.Message}}
	//AskForDate()
	//r.ReadMessage()
	//date := Response.BotUser.Message
	//AskForTime()
	//r.ReadMessage()
	//time := Response.BotUser.Message
	//Event{DateTime: MakeDateTimeField(date, time)}
	//AskHowFrequently()
	//frequency := Response.BotUser.Message
	//MakeFrequencyField(frequency)
	//Event{UserId: Response{BotUser.UserId}}
	//db := data.NewData(data.MustConnectPostgres())
	//db.UpdateEvent(Event.Id, Event.Name, Event.DateTime, Event.Weekly, Event.Monthly, Event.Yearly)
	//message := "Your event: " + Event.Name + " has been successfully edited"
	return "message"
}

func DeleteEvent() string {
	return "Deleting event"
}

func DeleteAllEvents() string {
	return "Deleting all events"
}
