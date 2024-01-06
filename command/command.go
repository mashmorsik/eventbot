package command

import (
	"database/sql"
	"eventbot/data"
	sendresponse "eventbot/send-response"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type (
	UserEvent struct {
		Data        *data.Data
		Message     *tgbotapi.Message
		BotAPI      *tgbotapi.BotAPI
		RerunEvents chan any
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
		Disabled       bool

		DeleteName string
		EditName   string
	}
)

var UserCurrentEvent = make(map[int64]*Steps)

func NewUserEvent(db *sql.DB, message *tgbotapi.Message, bot *tgbotapi.BotAPI, rerunEvents chan any) *UserEvent {
	return &UserEvent{Data: data.NewData(db), Message: message, BotAPI: bot, RerunEvents: rerunEvents}
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
	DisableCommand   = "disable"
	EnableCommand    = "enable"
	DeleteCommand    = "delete"
	DeleteAllCommand = "deleteall"

	OncePeriod = "once"
)

func (u UserEvent) HandleCommand() {
	if u.Message.Text == "" {
		sendresponse.EmptyText()
		return
	}

	u.Data.AddUser(u.Message.From.ID)

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
	case DisableCommand:
		err := u.handleDisable()
		if err != nil {
			fmt.Println(err)
		}
	case EnableCommand:
		err := u.handleEnable()
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
