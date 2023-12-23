package data

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
)

var commands = []string{"/newevent", "/myevents", "/edit", "/delete", "/deleteall"}
var UsersSteps = make(map[int64]*Steps)

const (
	Name      = "Name"
	Date      = "Date"
	Time      = "Time"
	Frequency = "Frequency"
)

type Steps struct {
	Step      string
	Name      string
	Date      string
	Time      string
	Frequency string
}

type BotUser struct {
	UserId  int64
	Message string
	ChatId  int64
}

type Bot struct {
	bot *tgbotapi.BotAPI
}

func NewBot(bot *tgbotapi.BotAPI) *Bot {
	if bot == nil {
		panic("Bot is nil")
	}

	return &Bot{bot: bot}
}

func BotStart() *tgbotapi.BotAPI {
	token, _ := os.LookupEnv("EVENTBOT_TOKEN")
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	return bot
}

func (r *Bot) ReadMessage() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := r.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			userId := update.Message.From.ID
			msgText := update.Message.Text

			switch msgText {
			case commands[0]:
				if v, ok := UsersSteps[userId]; !ok {
					UsersSteps[userId] = &Steps{
						Step:      Name,
						Name:      "",
						Date:      "",
						Time:      "",
						Frequency: "",
					}
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, AskForName())
					_, err := r.bot.Send(msg)
					if err != nil {
						fmt.Println(err)
					}
				} else {
					switch v.Step {
					case Name:
						UsersSteps[userId].Time = update.Message.Text
						v.Step = Time
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Name")
						_, err := r.bot.Send(msg)
						if err != nil {
							fmt.Println(err)
						}
					case Time:
						UsersSteps[userId].Time = update.Message.Text
						v.Step = Frequency
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Time")
						_, err := r.bot.Send(msg)
						if err != nil {
							fmt.Println(err)
						}
					case Frequency:
						UsersSteps[userId].Frequency = update.Message.Text
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Frequency")
						_, err := r.bot.Send(msg)
						if err != nil {
							fmt.Println(err)
						}
					}
				}
			case commands[1]:
				db := NewData(MustConnectPostgres())
				list, _ := db.GetEventsList(userId)
				var eventslist string

				for _, event := range list {
					item := "\n/" + event
					eventslist += item
				}
				fmt.Println(userId)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, eventslist)
				_, err := r.bot.Send(msg)
				if err != nil {
					fmt.Println(err)
				}
			case commands[3]:
				db := NewData(MustConnectPostgres())
				list, _ := db.GetEventsList(userId)
				var eventslist = "Which event do you want to delete?"

				for _, event := range list {
					item := "\n/" + event
					eventslist += item
				}
				fmt.Println(userId)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, eventslist)
				_, err := r.bot.Send(msg)
				if err != nil {
					fmt.Println(err)
				}
			case commands[4]:
				db := NewData(MustConnectPostgres())
				db.DeleteAllEvents(userId)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Your events have been successfully deleted.")
				_, err := r.bot.Send(msg)
				if err != nil {
					fmt.Println(err)
				}
			default:
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, WelcomeMessage())
				_, err := r.bot.Send(msg)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}
}

func WelcomeMessage() string {
	return "Hello! This is a reminder bot. \nSend \n/newevent to create an event " +
		"\n/myevents to get the list of your events " +
		"\n/edit to update your event, " +
		"\n/delete to delete one of your events, " +
		"\n/deleteall to delete all of your events."
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

//func (r *Bot) SendMessage() {
//
//	msg := tgbotapi.NewMessage(User.ChatId, sendResponse.SendResponse(User.Message))
//	_, err := r.bot.Send(msg)
//	if err != nil {
//		fmt.Println(err)
//	}
//}
