package command

import (
	"eventbot/data"
	send_response "eventbot/send-response"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type (
	UserEvent struct {
		Db      *data.Data
		Message *tgbotapi.Message
		Bot     *tgbotapi.BotAPI
	}

	Steps struct {
		CurrentStep string
		Name        string
		Date        string
		Time        string
		Frequency   string
	}
)

var UsersSteps = make(map[int64]*Steps)

func NewUserEvent(db *data.Data, message *tgbotapi.Message, bot *tgbotapi.BotAPI) *UserEvent {
	return &UserEvent{Db: db, Message: message, Bot: bot}
}

const (
	NameStep      = "Name"
	DateStep      = "Date"
	TimeStep      = "Time"
	FrequencyStep = "Frequency"

	NewEventCommand  = "newevent"
	MyEventsCommand  = "myevents"
	EditCommand      = "edit"
	DeleteCommand    = "delete"
	DeleteAllCommand = "deleteall"
)

func (u UserEvent) HandleCommand() {

	switch u.Message.Command() {
	case NewEventCommand:
		err := u.handleNewEvent()
		if err != nil {
			fmt.Println(err)
		}
		return
	case MyEventsCommand:
		//db := NewData(MustConnectPostgres())
		//list, _ := db.GetEventsList(userId)
		//var eventslist string
		//
		//for _, event := range list {
		//	item := "\n/" + event
		//	eventslist += item
		//}
		//fmt.Println(userId)
		//msg := tgbotapi.NewMessage(update.Message.Chat.ID, eventslist)
		//_, err := r.bot.Send(msg)
		//if err != nil {
		//	fmt.Println(err)
		//}
	case EditCommand:
		//db := NewData(MustConnectPostgres())
		//list, _ := db.GetEventsList(userId)
		//var eventslist = "Which event do you want to delete?"
		//
		//for _, event := range list {
		//	item := "\n/" + event
		//	eventslist += item
		//}
		//fmt.Println(userId)
		//msg := tgbotapi.NewMessage(update.Message.Chat.ID, eventslist)
		//_, err := r.bot.Send(msg)
		//if err != nil {
		//	fmt.Println(err)
		//}
	case DeleteAllCommand:
		//db := NewData(MustConnectPostgres())
		//db.DeleteAllEvents(userId)
		//msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Your events have been successfully deleted.")
		//_, err := r.bot.Send(msg)
		//if err != nil {
		//	fmt.Println(err)
		//}
	case DeleteCommand:
		// u.Db.DeleteAllEvents(u.UserMessage.UserId)
	default:
		//msg := tgbotapi.NewMessage(update.Message.Chat.ID, WelcomeMessage())
		//_, err := r.bot.Send(msg)
		//if err != nil {
		//	fmt.Println(err)
		//}
	}
}

func (u UserEvent) handleNewEvent() error {

	v, ok := UsersSteps[u.Message.From.ID]
	if !ok {
		UsersSteps[u.Message.From.ID] = &Steps{
			CurrentStep: NameStep,
		}

		_, err := u.Bot.Send(tgbotapi.NewMessage(u.Message.Chat.ID, send_response.AskForName()))
		return err
	}
	if u.Message.Text == "" {
		send_response.EmptyText()
		return nil
	}

	switch v.CurrentStep {
	case NameStep:
		UsersSteps[u.Message.From.ID].Name = u.Message.Text
		v.CurrentStep = DateStep
		_, err := u.Bot.Send(tgbotapi.NewMessage(u.Message.Chat.ID, send_response.AskForDate()))
		return err
	case DateStep:
		UsersSteps[u.Message.From.ID].Date = u.Message.Text
		v.CurrentStep = TimeStep
		_, err := u.Bot.Send(tgbotapi.NewMessage(u.Message.Chat.ID, send_response.AskForTime()))
		return err
		//case TimeStep:
		//	UsersSteps[userId].Time = update.Message.Text
		//	v.CurrentStep = Frequency
		//	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Time")
		//	_, err := r.bot.Send(msg)
		//	if err != nil {
		//		fmt.Println(err)
		//	}
		//case FrequencyStep:
		//	UsersSteps[userId].Frequency = update.Message.Text
		//	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Frequency")
		//	_, err := r.bot.Send(msg)
		//	if err != nil {
		//		fmt.Println(err)
		//	}
		//}
	}
	return nil
}
