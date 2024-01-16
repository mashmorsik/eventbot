package command

import (
	"errors"
	"eventbot/Logger"
	"eventbot/data"
	"eventbot/pkg/loc"
	sendresponse "eventbot/send-response"
	"fmt"
	"github.com/go-co-op/gocron"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

const (
	retryAttempts     = 3
	defaultTriesValue = 1
)

var (
	tries = defaultTriesValue
)

func (u UserEvent) HandleCommand() {
	if u.Message.Text == "" {
		sendresponse.EmptyText()
		return
	}

	err := u.Data.AddUser(u.Message.From.ID)
	if err != nil {
		panic("AddUser")
	}

	var currentCommand string
	v, ok := UserCurrentEvent[u.Message.From.ID]
	if ok {
		currentCommand = v.CurrentCommand
	} else {
		currentCommand = u.Message.Command()
	}

	switch currentCommand {
	case NewEventCommand:
		err = u.handleNewEvent()
		if err != nil {
			fmt.Println(err)
		}
	case MyEventsCommand:
		err = u.handleMyEvents()
		if err != nil {
			fmt.Println(err)
		}
	case EditCommand:
		err = u.handleEdit()
		if err != nil {
			fmt.Println(err)
		}
	case DisableCommand:
		err = u.handleDisable()
		if err != nil {
			fmt.Println(err)
		}
	case EnableCommand:
		err = u.handleEnable()
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

func (u UserEvent) handleDateStep(v *Steps) error {
	validDate := isValidDate(u.Message.Text)

	switch {
	case !validDate:
		if tries >= retryAttempts {
			if err := u.SendMessage(u.Message.Chat.ID, "I see you are not ready. Try later."); err != nil {
				return err
			}

			delete(UserCurrentEvent, u.Message.From.ID)
			break
		}

		// FIXME: use this https://github.com/go-playground/validator
		if err := u.SendMessage(u.Message.Chat.ID, "Try again. Use the YYYY-MM-DD format."); err != nil {
			return err
		}

		tries++
		return nil

	case validDate:
		UserCurrentEvent[u.Message.From.ID].Date = u.Message.Text
		v.CurrentStep = TimeStep
		_, err := u.BotAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, sendresponse.AskForTime()))
		if err != nil {
			return err
		}
		break

	default:
		panic("default case")
		return nil
	}

	tries = defaultTriesValue
	return nil
}

func (u UserEvent) handleTimeStep(v *Steps) error {
	validTime := isValidTime(u.Message.Text)

	switch {
	case !validTime:
		if tries >= retryAttempts {
			_, err := u.BotAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "Wrong format, come back later."))
			if err != nil {
				return err
			}
			delete(UserCurrentEvent, u.Message.From.ID)
			break
		}

		_, err := u.BotAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "Try again. Use the HH:MM 24h format."))
		if err != nil {
			return err
		}
		tries += 1
		return nil

	case validTime:
		UserCurrentEvent[u.Message.From.ID].Time = u.Message.Text
		v.CurrentStep = FrequencyStep
		_, err := u.BotAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, sendresponse.AskHowFrequently()))
		if err != nil {
			return err
		}
		break

	default:
		panic("default case")
		return nil
	}

	tries = defaultTriesValue
	return nil
}

func (u UserEvent) handleFrequencyStep(v *Steps, finalFunc string) error {
	validFrequency := isValidFrequency(u.Message.Text)
	userId := u.Message.From.ID

	switch {
	case !validFrequency:
		if tries >= retryAttempts {
			_, err := u.BotAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "Sorry, I don't understand. Try later."))
			if err != nil {
				return err
			}
			delete(UserCurrentEvent, userId)
			break
		}

		_, err := u.BotAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "Try again. Use only: once, daily, weekly, monthly or yearly."))
		if err != nil {
			return err
		}
		tries += 1
		return nil

	case validFrequency:
		v.Frequency = u.Message.Text

		var cron string

		cron = sendresponse.StringToCron(v.Date, v.Time, v.Frequency)

		if finalFunc == "createTask" {
			eventId, err := u.Data.CreateEvent(userId, v.ChatId, v.Name, sendresponse.MakeDateTimeField(v.Date,
				v.Time), cron)
			if err != nil {
				Logger.Sugar.Errorln("Couldn't get EventId from DB ", err)
			}
			fmt.Println(eventId)

			if cron != "once" {
				_, err = u.setCronRepeatable(cron, userId, v.Name, eventId)
				if err != nil {
					Logger.Sugar.Errorln("Couldn't create new cron job ", err)
				}
			}

			_, err = u.BotAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "Your event has been successfully created."))
			if err != nil {
				fmt.Println(err)
			}
		} else {
			err := u.Data.UpdateEvent(v.EventId, v.Name, sendresponse.MakeDateTimeField(v.Date,
				v.Time), cron)
			if err != nil {
				return err
			}

			s := u.Scheduler.Sc()

			if cron != "once" {
				err = s.RemoveByTag(string(rune(v.EventId)))
				if err != nil {
					Logger.Sugar.Errorln("Couldn't remove cron job after editing event.")
				}
				_, err = u.setCronRepeatable(cron, u.Message.Chat.ID, v.Name, v.EventId)
				if err != nil {
					Logger.Sugar.Errorln("Couldn't create cron job after editing.")
				}
			}

			_, err = u.BotAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "Your event has been successfully edited."))
			if err != nil {
				Logger.Sugar.Errorln("Couldn't send message after editing event.")
			}
		}

		delete(UserCurrentEvent, userId)

		break
	default:
		panic("default case")
		return nil
	}

	tries = defaultTriesValue
	return nil
}

func (u UserEvent) handleNewEvent() error {
	userId := u.Message.From.ID

	v, ok := UserCurrentEvent[userId]
	if !ok {
		UserCurrentEvent[userId] = &Steps{
			CurrentCommand: NewEventCommand,
			CurrentStep:    NameStep,
			ChatId:         u.Message.Chat.ID,
		}

		if _, err := u.BotAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, sendresponse.AskForName())); err != nil {
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
		v.Name = u.Message.Text
		v.CurrentStep = DateStep

		return u.SendMessage(u.Message.Chat.ID, sendresponse.AskForDate())
	case DateStep:
		err := u.handleDateStep(v)
		if err != nil {
			Logger.Sugar.Errorln("HandleDateStep failed.")
		}
	case TimeStep:
		err := u.handleTimeStep(v)
		if err != nil {
			Logger.Sugar.Errorln("HandleTimeStep failed.")
		}
	case FrequencyStep:
		err := u.handleFrequencyStep(v, "createTask")
		if err != nil {
			Logger.Sugar.Errorln("HandleFrequencyStep failed.")
		}
	}

	return nil
}

func (u UserEvent) setCronRepeatable(cron string, chatID int64, eventName string, eventId int) (*gocron.Scheduler, error) {
	s, err := u.RunScheduler()
	if err != nil {
		Logger.Sugar.Errorln("Couldn't get scheduler")
		return nil, err
	}

	s.TagsUnique()

	_, err = s.Cron(cron).Tag(string(rune(eventId))).Do(func() error {
		if err = u.SendMessage(chatID, "Don't forget about "+eventName); err != nil {
			return err
		}

		err = u.Data.SetLastFired(time.Now(), eventId)
		if err != nil {
			Logger.Sugar.Errorln("Couldn't set last fired value for repeating event.")
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (u UserEvent) handleMyEvents() error {
	userId := u.Message.From.ID
	db := data.NewData(data.MustConnectPostgres())

	list, err := db.GetUserEvents(userId)
	if err != nil {
		return err
	}

	var eventsList string

	for _, event := range list {
		if event.Disabled == true {
			eventsList += "\n/" + event.Name + " disabled"
		} else {
			eventsList += "\n/" + event.Name
		}
	}

	response := "Your events: \n" + eventsList
	Logger.Sugar.Infoln(userId)

	return u.SendMessage(u.Message.Chat.ID, response)
}

func (u UserEvent) handleEdit() error {
	userId := u.Message.From.ID

	var list map[int]*data.Event
	var eventsList string

	list, err := u.Data.GetUserEvents(userId)
	if err != nil {
		return err
	}

	v, ok := UserCurrentEvent[userId]
	if !ok {
		UserCurrentEvent[userId] = &Steps{
			CurrentCommand: EditCommand,
			CurrentStep:    EditNameStep,
			ChatId:         u.Message.Chat.ID,
		}

		for _, event := range list {
			item := "\n" + event.Name
			eventsList += item
		}

		return u.SendMessage(u.Message.Chat.ID, "Which event do you want to edit? \n"+eventsList)
	}
	if u.Message.Text == "" {
		sendresponse.EmptyText()
		return nil
	}

	switch v.CurrentStep {
	case EditNameStep:

		v.EditName = u.Message.Text
		for id, e := range list {
			if e.Name == v.EditName {
				v.EventId = id
				v.CurrentStep = NameStep

				return u.SendMessage(u.Message.Chat.ID, sendresponse.AskForName())
			}

			if err = u.SendMessageUnknownEvent(u.Message.Chat.ID); err != nil {
				return err
			}

			delete(UserCurrentEvent, u.Message.From.ID)
			return nil
		}
		return nil
	case NameStep:
		v.Name = u.Message.Text
		v.CurrentStep = DateStep

		return u.SendMessage(u.Message.Chat.ID, sendresponse.AskForDate())
	case DateStep:
		err = u.handleDateStep(v)
		if err != nil {
			Logger.Sugar.Errorln("HandleDateStep failed.")
		}
	case TimeStep:
		err = u.handleTimeStep(v)
		if err != nil {
			Logger.Sugar.Errorln("HandleTimeStep failed.")
		}
	case FrequencyStep:
		v.Frequency = u.Message.Text
		err = u.handleFrequencyStep(v, "updateTask")
		if err != nil {
			Logger.Sugar.Errorln("HandleFrequencyStep failed.")
		}
	}

	return nil
}

func (u UserEvent) handleDisable() error {
	var (
		eventsList string
		eventsName map[int]*data.Event
	)

	userId := u.Message.From.ID

	eventsName, err := u.Data.GetUserEvents(userId)
	if err != nil {
		return err
	}

	_, ok := UserCurrentEvent[userId]
	if !ok {
		UserCurrentEvent[userId] = &Steps{
			CurrentCommand: DisableCommand,
			CurrentStep:    NameStep,
		}

		for _, event := range eventsName {
			if event.Disabled == false {
				eventsList += "\n" + event.Name
			}
		}

		return u.SendMessage(userId, "Which event do you want to disable? \n"+eventsList)
	}

	if v, _ := UserCurrentEvent[userId]; v.CurrentStep == NameStep {
		var eventId int

		for _, e := range eventsName {
			if e.Name == u.Message.Text {
				eventId = e.EventId

				err = u.Data.DisabledTrue(eventId)
				if err != nil {
					return err
				}

				if err = u.SendMessage(userId, "Event disabled"); err != nil {
					return err
				}

				s := u.Scheduler.Sc()
				if e.Cron != OncePeriod {
					err = s.RemoveByTag(string(rune(eventId)))
					if err != nil {
						Logger.Sugar.Errorln("Couldn't remove cron job after editing event.")
					}
				}
				delete(UserCurrentEvent, userId)
				return err
			}

			if err = u.SendMessageUnknownEvent(u.Message.Chat.ID); err != nil {
				return err
			}
			delete(UserCurrentEvent, userId)
			return nil
		}

		return nil
	}
	return nil
}

func (u UserEvent) handleEnable() error {
	var (
		eventsList string
		eventsName map[int]*data.Event
	)

	userId := u.Message.From.ID

	eventsName, err := u.Data.GetUserEvents(userId)
	if err != nil {
		return err
	}

	_, ok := UserCurrentEvent[userId]
	if !ok {
		UserCurrentEvent[userId] = &Steps{
			CurrentCommand: EnableCommand,
			CurrentStep:    NameStep,
		}

		for _, event := range eventsName {
			if event.Disabled == true {
				eventsList += "\n" + event.Name
			}
		}

		if len(eventsList) > 0 {
			if err = u.SendMessage(u.Message.Chat.ID, "Which event do you want to enable? \n"+eventsList); err != nil {
				return err
			}
			return nil
		}

		if err = u.SendMessage(u.Message.Chat.ID, "You don't have any disabled events."); err != nil {
			return err
		}

		delete(UserCurrentEvent, userId)
		return nil
	}

	if v, _ := UserCurrentEvent[userId]; v.CurrentStep == NameStep {
		var eventId int

		for _, e := range eventsName {
			if e.Name == u.Message.Text {
				eventId = e.EventId

				u.Data.DisabledFalse(eventId)

				if err = u.SendMessage(u.Message.Chat.ID, "Event enabled."); err != nil {
					return err
				}

				if e.Cron != "once" {
					_, err = u.setCronRepeatable(e.Cron, u.Message.Chat.ID, e.Name, eventId)
					if err != nil {
						Logger.Sugar.Errorln("Couldn't create cron job after enabling event.")
					}
				}
				delete(UserCurrentEvent, userId)
				return err
			}

			if err = u.SendMessageUnknownEvent(u.Message.Chat.ID); err != nil {
				return err
			}

			delete(UserCurrentEvent, userId)
			return nil
		}

		return nil
	}

	return nil
}

func (u UserEvent) handleDelete() error {
	var (
		eventsList string
		eventsName map[int]*data.Event
	)

	userId := u.Message.From.ID

	eventsName, err := u.Data.GetUserEvents(userId)
	if err != nil {
		return err
	}

	v, ok := UserCurrentEvent[userId]
	if !ok {
		UserCurrentEvent[userId] = &Steps{
			CurrentCommand: DeleteCommand,
			CurrentStep:    DeleteNameStep,
		}

		for _, event := range eventsName {
			eventsList += "\n" + event.Name
		}

		return u.SendMessage(u.Message.Chat.ID, "Which event do you want to delete? \n"+eventsList)
	}

	if v.CurrentStep == DeleteNameStep {
		var eventId int

		for id, e := range eventsName {
			if e.Name == u.Message.Text {
				eventId = id

				if err = u.Data.DeleteEvent(eventId); err != nil {
					return err
				}

				if err = u.SendMessage(u.Message.Chat.ID, "Event deleted"); err != nil {
					return err
				}
				delete(UserCurrentEvent, u.Message.Chat.ID)
				return nil
			}
			if err = u.SendMessageUnknownEvent(u.Message.Chat.ID); err != nil {
				return err
			}
			delete(UserCurrentEvent, u.Message.Chat.ID)
			return nil
		}

		return nil
	}

	return nil
}

func (u UserEvent) handleDeleteAll() error {
	userId := u.Message.From.ID
	db := data.NewData(data.MustConnectPostgres())
	err := db.DeleteAllEvents(userId)
	if err != nil {
		return err
	}

	return u.SendMessage(u.Message.Chat.ID, "Your events have been successfully deleted.")
}

func (u UserEvent) handleDefault() error {
	return u.SendMessage(u.Message.Chat.ID, sendresponse.WelcomeMessage())
}

func (u UserEvent) HandleOnceReminder() error {
	events, err := u.Data.GetOnceNotFired()
	if err != nil {
		Logger.Sugar.Errorln("No not fired once events")
		return err
	}
	for _, e := range events {
		if time.Now().After(e.TimeDate.In(loc.CurrentLoc)) {
			if err = u.SendMessage(e.ChatId, "Don't forget about "+e.Name); err != nil {
				return err
			}

			err = u.Data.SetLastFired(time.Now(), e.EventId)
			if err != nil {
				return errors.New(fmt.Sprintf("u.Data.SetLastFired failed with error: %s", err))
			}
		}
	}
	return nil
}

func (u UserEvent) RunScheduler() (*gocron.Scheduler, error) {
	scheduler := u.Scheduler.Sc()
	defer scheduler.StartAsync()

	_, err := scheduler.Every(1).Minute().Do(func() {
		Logger.Sugar.Infoln("HandleOnceReminder fired")
		err := u.HandleOnceReminder()
		if err != nil {
			Logger.Sugar.Errorln("HandleOnceReminder failed", err)
		}
	})
	if err != nil {
		Logger.Sugar.Errorln("Couldn't create job")
	}

	jobs := scheduler.Jobs()
	if len(jobs) == 1 {
		events, err := u.Data.GetCronMultipleActive()
		if err != nil {
			Logger.Sugar.Errorln("Couldn't get events from DB.")
		}

		for _, e := range events {
			_, err = u.setCronRepeatable(e.Cron, e.ChatId, e.Name, e.EventId)
			if err != nil {
				Logger.Sugar.Errorln("Couldn't add cron jobs from DB.")
			}
		}
	}

	return scheduler, nil
}

func (u UserEvent) SendMessage(chatID int64, text string) error {
	if _, err := u.BotAPI.Send(tgbotapi.NewMessage(chatID, text)); err != nil {
		Logger.Sugar.Errorln(err)
		return err
	}
	return nil
}

func (u UserEvent) SendMessageUnknownEvent(chatID int64) error {
	if _, err := u.BotAPI.Send(tgbotapi.NewMessage(chatID, "Unknown event")); err != nil {
		Logger.Sugar.Errorln(err)
		return err
	}
	return nil
}
