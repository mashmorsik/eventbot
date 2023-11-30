package SendResponse

import (
	"context"
	"eventbot/Calendar"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/datepicker"
	"time"
)

var commands = []string{"/newevent", "/myevents", "/edit", "/delete", "/deleteall"}

func SendResponse(text string) string {
	switch text {
	case commands[0]:
		CreateEvent()
	case commands[1]:
		GetEvents()
	case commands[2]:
		EditEvent()
	case commands[3]:
		DeleteEvent()
	case commands[4]:
		DeleteAllEvents()
	}
	return WelcomeMessage()
}

func WelcomeMessage() string {
	return "Hello. This is reminder bot. Send /newevent to create an event, " +
		"/myevents to get the list of your events, " +
		"/edit to update your event, " +
		"/delete to delete your event, " +
		"/deleteall to delete all of your events."
}

func SendDatePicker(ctx context.Context, b *bot.Bot, update *models.Update) {
	kb := datepicker.New(b, onDatepickerSimpleSelect)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        "Select any date",
		ReplyMarkup: kb,
	})
}

func onDatepickerSimpleSelect(ctx context.Context, b *bot.Bot, mes *models.Message, date time.Time) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: mes.Chat.ID,
		Text:   "You select " + date.Format("2006-01-02"),
	})
}

func CreateEvent() tgbotapi.InlineKeyboardMarkup {
	calendar := Calendar.GenerateCalendar(2023, 11)
	return calendar
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
