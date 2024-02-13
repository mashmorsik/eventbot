package telegram

import (
	"database/sql"
	"eventbot/cron"
	"eventbot/data"
	"eventbot/internal/command"
	"github.com/go-co-op/gocron"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type (
	UserEvent struct {
		Data        *data.Data
		Message     *tgbotapi.Message
		BotAPI      *tgbotapi.BotAPI
		RerunEvents chan any
		Scheduler   *cronkafka.Scheduler
		command     command.EventInterface
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

func NewUserEvent(db *sql.DB, message *tgbotapi.Message, bot *tgbotapi.BotAPI, rerunEvents chan any, sc *gocron.Scheduler,
) command.UserEventer {
	return UserEvent{
		Data: data.NewData(db), Message: message, BotAPI: bot, RerunEvents: rerunEvents, Scheduler: cronkafka.NewScheduler(sc),
	}
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
