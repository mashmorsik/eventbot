package SendResponse

var commands = []string{"/newevent", "/myevents", "/edit", "/delete", "/deleteall"}

func SendResponse(text string) string {
	switch text {
	case commands[0]:
		CreateEvent()
	case commands[1]:
		GetEvents()
	case commands[2]:
		EditEvent()
	case commands[4]:
		DeleteEvent()
	case commands[5]:
		DeleteAllEvents()
	}
	return WelcomeMessage()
}

func WelcomeMessage() string {
	return "Hello. This is reminder bot. Send /newevent to create an event."
}

func CreateEvent() string {
	//date := datepicker.New(bot, onDatepickerSimpleSelect)
	return "Creating event"
}

func GetEvents() string {
	return "Returning list of events"
}

func EditEvent() string {
	return "Editing event"
}

func DeleteEvent() string {
	return "Deleting event"
}

func DeleteAllEvents() string {
	return "Deleting all events"
}
