package command

import (
	"database/sql"
	"eventbot/data"
	sendresponse "eventbot/send-response"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

type (
	UserEvent struct {
		Db      *data.Data
		Message *tgbotapi.Message
		Bot     *tgbotapi.BotAPI
	}

	Steps struct {
		CurrentCommand string
		CurrentStep    string
		EventId        int
		ChatId         int64
		Name           string
		Date           string
		Time           string
		Frequency      string

		DeleteName string
		EditName   string
	}
)

var UserCurrentEvent = make(map[int64]*Steps)

func NewUserEvent(db *sql.DB, message *tgbotapi.Message, bot *tgbotapi.BotAPI) *UserEvent {
	return &UserEvent{Db: data.NewData(data.MustConnectPostgres()), Message: message, Bot: bot}
}

const (
	NameStep       = "Name"
	DateStep       = "Date"
	TimeStep       = "Time"
	FrequencyStep  = "Frequency"
	DeleteNameStep = "DeleteName"
	EditNameStep   = "EditName"

	NewEventCommand  = "newevent"
	MyEventsCommand  = "myevents"
	EditCommand      = "edit"
	DeleteCommand    = "delete"
	DeleteAllCommand = "deleteall"
)

func (u UserEvent) HandleCommand() {
	if u.Message.Text == "" {
		sendresponse.EmptyText()
		return
	}

	u.Db.AddUser(u.Message.From.ID)

	go u.handleReminder()

	var currentCommand string
	v, ok := UserCurrentEvent[u.Message.From.ID]
	if ok {
		currentCommand = v.CurrentCommand
	} else {
		currentCommand = u.Message.Command()
	}

	switch currentCommand {
	case NewEventCommand:
		err := u.handleNewEvent()
		if err != nil {
			fmt.Println(err)
		}
		return
	case MyEventsCommand:
		err := u.handleMyEvents()
		if err != nil {
			fmt.Println(err)
		}
	case EditCommand:
		err := u.handleEdit()
		if err != nil {
			fmt.Println(err)
		}
	case DeleteCommand:
		err := u.handleDelete()
		if err != nil {
			fmt.Println(err)
		}
	case DeleteAllCommand:
		err := u.handleDeleteAll()
		if err != nil {
			fmt.Println(err)
		}
	default:
		err := u.handleDefault()
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (u UserEvent) handleNewEvent() error {
	v, ok := UserCurrentEvent[u.Message.From.ID]
	if !ok {
		UserCurrentEvent[u.Message.From.ID] = &Steps{
			CurrentCommand: NewEventCommand,
			CurrentStep:    NameStep,
			ChatId:         u.Message.Chat.ID,
		}

		if _, err := u.Bot.Send(tgbotapi.NewMessage(u.Message.Chat.ID, sendresponse.AskForName())); err != nil {
			return err
		}
		return nil
	}
	if u.Message.Text == "" {
		sendresponse.EmptyText()
		return nil
	}

	switch v.CurrentStep {
	case NameStep:
		UserCurrentEvent[u.Message.From.ID].Name = u.Message.Text
		v.CurrentStep = DateStep
		_, err := u.Bot.Send(tgbotapi.NewMessage(u.Message.Chat.ID, sendresponse.AskForDate()))
		return err
	case DateStep:
		UserCurrentEvent[u.Message.From.ID].Date = u.Message.Text
		v.CurrentStep = TimeStep
		_, err := u.Bot.Send(tgbotapi.NewMessage(u.Message.Chat.ID, sendresponse.AskForTime()))
		return err
	case TimeStep:
		UserCurrentEvent[u.Message.From.ID].Time = u.Message.Text
		v.CurrentStep = FrequencyStep
		_, err := u.Bot.Send(tgbotapi.NewMessage(u.Message.Chat.ID, sendresponse.AskHowFrequently()))
		return err
	case FrequencyStep:
		UserCurrentEvent[u.Message.From.ID].Frequency = u.Message.Text

		var cron string

		cron = sendresponse.StringToCron(UserCurrentEvent[u.Message.From.ID].Date,
			UserCurrentEvent[u.Message.From.ID].Time, UserCurrentEvent[u.Message.From.ID].Frequency)

		u.Db.CreateEvent(u.Message.From.ID, UserCurrentEvent[u.Message.From.ID].ChatId, UserCurrentEvent[u.Message.From.ID].Name,
			sendresponse.MakeDateTimeField(UserCurrentEvent[u.Message.From.ID].Date,
				UserCurrentEvent[u.Message.From.ID].Time), cron)

		_, err := u.Bot.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "Your event has been successfully created."))
		if err != nil {
			fmt.Println(err)
		}
		delete(UserCurrentEvent, u.Message.From.ID)
	}
	return nil
}

func (u UserEvent) handleMyEvents() error {
	userId := u.Message.From.ID
	db := data.NewData(data.MustConnectPostgres())

	list, _ := db.GetEventsList(userId)
	var eventslist string

	for _, event := range list {
		item := "\n/" + event
		eventslist += item
	}

	response := "Your events: \n" + eventslist
	fmt.Println(userId)
	_, err := u.Bot.Send(tgbotapi.NewMessage(u.Message.Chat.ID, response))
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func (u UserEvent) handleEdit() error {
	userId := u.Message.From.ID

	var list map[int]string
	var eventsList string

	list, err := u.Db.GetEventsList(userId)
	if err != nil {
		return err
	}

	v, ok := UserCurrentEvent[userId]
	if !ok {
		UserCurrentEvent[userId] = &Steps{
			CurrentCommand: EditCommand,
			CurrentStep:    EditNameStep,
		}

		//list, _ = db.GetEventsList(userId)
		for _, event := range list {
			item := "\n" + event
			eventsList += item
		}

		response := "Which event do you want to edit? \n" + eventsList
		_, err := u.Bot.Send(tgbotapi.NewMessage(u.Message.Chat.ID, response))
		return err
	}
	if u.Message.Text == "" {
		sendresponse.EmptyText()
		return nil
	}

	var eventId int

	switch v.CurrentStep {
	case EditNameStep:

		UserCurrentEvent[u.Message.From.ID].EditName = u.Message.Text
		//for id, e := range list {
		//	if e == UserCurrentEvent[u.Message.From.ID].DeleteName {
		//		UserCurrentEvent[u.Message.From.ID].EventId = id
		//		fmt.Println(UserCurrentEvent[u.Message.From.ID].EventId)
		//	}
		//}
		fmt.Println(UserCurrentEvent[u.Message.From.ID].EventId)
		v.CurrentStep = NameStep
		_, err := u.Bot.Send(tgbotapi.NewMessage(u.Message.Chat.ID, sendresponse.AskForName()))
		return err
	case NameStep:
		UserCurrentEvent[u.Message.From.ID].Name = u.Message.Text
		v.CurrentStep = DateStep
		_, err := u.Bot.Send(tgbotapi.NewMessage(u.Message.Chat.ID, sendresponse.AskForDate()))
		return err
	case DateStep:
		UserCurrentEvent[u.Message.From.ID].Date = u.Message.Text
		v.CurrentStep = TimeStep
		_, err := u.Bot.Send(tgbotapi.NewMessage(u.Message.Chat.ID, sendresponse.AskForTime()))
		return err
	case TimeStep:
		UserCurrentEvent[u.Message.From.ID].Time = u.Message.Text
		v.CurrentStep = FrequencyStep
		_, err := u.Bot.Send(tgbotapi.NewMessage(u.Message.Chat.ID, sendresponse.AskHowFrequently()))
		return err
	case FrequencyStep:
		UserCurrentEvent[u.Message.From.ID].Frequency = u.Message.Text

		var weekly bool
		var monthly bool
		var yearly bool

		if UserCurrentEvent[u.Message.From.ID].Frequency == "weekly" {
			weekly = true
			monthly = false
			yearly = false
		} else if UserCurrentEvent[u.Message.From.ID].Frequency == "monthly" {
			weekly = false
			monthly = true
			yearly = false
		} else {
			weekly = false
			monthly = false
			yearly = true
		}

		for id, e := range list {
			if e == UserCurrentEvent[u.Message.From.ID].EditName {
				eventId = id
				fmt.Println(eventId)
				UserCurrentEvent[u.Message.From.ID].EventId = eventId
			}
		}

		u.Db.UpdateEvent(eventId, UserCurrentEvent[u.Message.From.ID].Name,
			sendresponse.MakeDateTimeField(UserCurrentEvent[u.Message.From.ID].Date,
				UserCurrentEvent[u.Message.From.ID].Time), weekly, monthly, yearly)

		_, err := u.Bot.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "Your event has been successfully edited."))
		if err != nil {
			fmt.Println(err)
		}
		delete(UserCurrentEvent, u.Message.From.ID)
	}

	return nil
}

func (u UserEvent) handleDelete() error {
	var (
		eventsList string
		eventsName map[int]string
	)

	userId := u.Message.From.ID

	eventsName, err := u.Db.GetEventsList(userId)
	if err != nil {
		return err
	}

	// first message
	v, ok := UserCurrentEvent[userId]
	if !ok {
		UserCurrentEvent[userId] = &Steps{
			CurrentCommand: DeleteCommand,
			CurrentStep:    DeleteNameStep,
		}

		for _, event := range eventsName {
			eventsList += "\n" + event
		}

		_, err = u.Bot.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "Which event do you want to delete? \n"+eventsList))
		if err != nil {
			return err
		}

		return nil
	}

	if v.CurrentStep == DeleteNameStep {
		var eventId int

		for id, e := range eventsName {
			if e == u.Message.Text {
				eventId = id
				fmt.Println(eventId)
			}
		}

		fmt.Println(u.Message.Text, eventId)

		u.Db.DeleteEvent(eventId)
		_, err = u.Bot.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "Event deleted"))
		if err != nil {
			return err
		}

		return nil
	}
	return nil
}

func (u UserEvent) handleDeleteAll() error {
	userId := u.Message.From.ID
	db := data.NewData(data.MustConnectPostgres())
	db.DeleteAllEvents(userId)
	_, err := u.Bot.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "Your events have been successfully deleted."))
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func (u UserEvent) handleDefault() error {
	msg := tgbotapi.NewMessage(u.Message.Chat.ID, sendresponse.WelcomeMessage())
	_, err := u.Bot.Send(msg)
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func (u UserEvent) handleReminder() {
	for {
		var remindEvents, _ = u.Db.FindRemindEvent()
		for _, v := range remindEvents {
			_, err := u.Bot.Send(tgbotapi.NewMessage(v.ChatId, v.Name))
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(v.Name)
		}
		time.Sleep(30 * time.Second)
	}
}

func (u UserEvent) HandleCronResponse(chatId int64, name string) error {
	_, err := u.Bot.Send(tgbotapi.NewMessage(chatId, "Don't forget about "+name))
	if err != nil {
		return err
	}

	return nil
}
