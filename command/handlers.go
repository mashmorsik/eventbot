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

func (u UserEvent) handleNewEvent() error {
	v, ok := UserCurrentEvent[u.Message.From.ID]
	if !ok {
		UserCurrentEvent[u.Message.From.ID] = &Steps{
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
		UserCurrentEvent[u.Message.From.ID].Name = u.Message.Text
		v.CurrentStep = DateStep
		_, err := u.BotAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, sendresponse.AskForDate()))
		return err
	case DateStep:
		UserCurrentEvent[u.Message.From.ID].Date = u.Message.Text
		v.CurrentStep = TimeStep
		_, err := u.BotAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, sendresponse.AskForTime()))
		return err
	case TimeStep:
		UserCurrentEvent[u.Message.From.ID].Time = u.Message.Text
		v.CurrentStep = FrequencyStep
		_, err := u.BotAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, sendresponse.AskHowFrequently()))
		return err
	case FrequencyStep:
		UserCurrentEvent[u.Message.From.ID].Frequency = u.Message.Text

		var cron string

		cron = sendresponse.StringToCron(UserCurrentEvent[u.Message.From.ID].Date,
			UserCurrentEvent[u.Message.From.ID].Time, UserCurrentEvent[u.Message.From.ID].Frequency)

		// make it a struct
		eventId, err := u.Data.CreateEvent(u.Message.From.ID, UserCurrentEvent[u.Message.From.ID].ChatId, UserCurrentEvent[u.Message.From.ID].Name,
			sendresponse.MakeDateTimeField(UserCurrentEvent[u.Message.From.ID].Date,
				UserCurrentEvent[u.Message.From.ID].Time), cron)
		if err != nil {
			Logger.Sugar.Errorln("Couldn't get EventId from DB ", err)
		}

		if cron != "once" {
			_, err = u.setCronRepeatable(cron, u.Message.From.ID, UserCurrentEvent[u.Message.From.ID].Name, eventId)
			if err != nil {
				Logger.Sugar.Errorln("Couldn't create new cron job ", err)
			}
		}

		_, err = u.BotAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "Your event has been successfully created."))
		if err != nil {
			fmt.Println(err)
		}
		delete(UserCurrentEvent, u.Message.From.ID)
	}

	return nil
}

func (u UserEvent) setCronRepeatable(cron string, chatID int64, eventName string, eventId int) (*gocron.Scheduler, error) {
	s, err := u.RunScheduler()
	if err != nil {
		Logger.Sugar.Errorln("Couldn't get scheduler")
	}

	s.TagsUnique()

	s.Cron(cron).Tag(string(rune(eventId))).Do(func() {
		_, err = u.BotAPI.Send(tgbotapi.NewMessage(chatID, "Don't forget about "+eventName))
		if err != nil {
			fmt.Println(err)
		}
		err = u.Data.SetLastFired(time.Now(), eventId)
		if err != nil {
			Logger.Sugar.Errorln("Couldn't set last fired value for repeating event.")
		}
	})

	return s, nil
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
	_, err := u.BotAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, response))
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func (u UserEvent) handleEdit() error {
	userId := u.Message.From.ID

	var list map[int]string
	var eventsList string

	list, err := u.Data.GetEventsList(userId)
	if err != nil {
		return err
	}

	v, ok := UserCurrentEvent[userId]
	if !ok {
		UserCurrentEvent[userId] = &Steps{
			CurrentCommand: EditCommand,
			CurrentStep:    EditNameStep,
		}

		for _, event := range list {
			item := "\n" + event
			eventsList += item
		}

		_, err = u.BotAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "Which event do you want to edit? \n"+eventsList))
		return err
	}
	if u.Message.Text == "" {
		sendresponse.EmptyText()
		return nil
	}

	switch v.CurrentStep {
	case EditNameStep:

		UserCurrentEvent[u.Message.From.ID].EditName = u.Message.Text
		for id, e := range list {
			if e == UserCurrentEvent[u.Message.From.ID].EditName {
				UserCurrentEvent[u.Message.From.ID].EventId = id
				fmt.Println(UserCurrentEvent[u.Message.From.ID].EventId)
			}
		}
		fmt.Println(UserCurrentEvent[u.Message.From.ID].EventId)
		v.CurrentStep = NameStep
		_, err = u.BotAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, sendresponse.AskForName()))
		return err
	case NameStep:
		UserCurrentEvent[u.Message.From.ID].Name = u.Message.Text
		v.CurrentStep = DateStep
		_, err = u.BotAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, sendresponse.AskForDate()))
		return err
	case DateStep:
		UserCurrentEvent[u.Message.From.ID].Date = u.Message.Text
		v.CurrentStep = TimeStep
		_, err = u.BotAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, sendresponse.AskForTime()))
		return err
	case TimeStep:
		UserCurrentEvent[u.Message.From.ID].Time = u.Message.Text
		v.CurrentStep = FrequencyStep
		_, err = u.BotAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, sendresponse.AskHowFrequently()))
		return err
	case FrequencyStep:
		UserCurrentEvent[u.Message.From.ID].Frequency = u.Message.Text

		var cron string

		cron = sendresponse.StringToCron(UserCurrentEvent[u.Message.From.ID].Date,
			UserCurrentEvent[u.Message.From.ID].Time, UserCurrentEvent[u.Message.From.ID].Frequency)

		// send struct
		u.Data.UpdateEvent(UserCurrentEvent[u.Message.From.ID].EventId, UserCurrentEvent[u.Message.From.ID].Name,
			sendresponse.MakeDateTimeField(UserCurrentEvent[u.Message.From.ID].Date,
				UserCurrentEvent[u.Message.From.ID].Time), cron)

		s := u.Scheduler.Sc()

		if cron != "once" {
			err = s.RemoveByTag(string(rune(UserCurrentEvent[u.Message.From.ID].EventId)))
			if err != nil {
				Logger.Sugar.Errorln("Couldn't remove cron job after editing event.")
			}
			_, err = u.setCronRepeatable(cron, u.Message.Chat.ID, UserCurrentEvent[u.Message.From.ID].Name, UserCurrentEvent[u.Message.From.ID].EventId)
			if err != nil {
				Logger.Sugar.Errorln("Couldn't create cron job after editing.")
			}
		}

		_, err = u.BotAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "Your event has been successfully edited."))
		if err != nil {
			Logger.Sugar.Errorln("Couldn't send message after editing event.")
		}
		delete(UserCurrentEvent, u.Message.From.ID)
	}

	return nil
}

func (u UserEvent) handleDisable() error {
	var (
		eventsList string
		eventsName map[int]*data.Event
	)

	userId := u.Message.From.ID

	// get only current user events
	eventsName, err := u.Data.GetAllEvents()
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

		_, err = u.BotAPI.Send(tgbotapi.NewMessage(userId, "Which event do you want to disable? \n"+eventsList))
		if err != nil {
			return err
		}
	}

	if v, _ := UserCurrentEvent[userId]; v.CurrentStep == NameStep {
		var eventId int

		for _, e := range eventsName {
			if e.Name == u.Message.Text {
				eventId = e.EventId

				u.Data.DisabledTrue(eventId)
				_, err = u.BotAPI.Send(tgbotapi.NewMessage(userId, "Event disabled"))
				if err != nil {
					return err
				}

				s := u.Scheduler.Sc()
				if e.Cron != "once" {
					err = s.RemoveByTag(string(rune(eventId)))
					if err != nil {
						Logger.Sugar.Errorln("Couldn't remove cron job after editing event.")
					}
				}
			}
		}
		delete(UserCurrentEvent, userId)
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

	eventsName, err := u.Data.GetAllEvents()
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
			_, err = u.BotAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "Which event do you want to enable? \n"+eventsList))
			if err != nil {
				return err
			}
		} else {
			_, err = u.BotAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "You don't have any disabled events."))
			if err != nil {
				return err
			}
		}
	}

	if v, _ := UserCurrentEvent[userId]; v.CurrentStep == NameStep {
		var eventId int

		for _, e := range eventsName {
			if e.Name == u.Message.Text {
				eventId = e.EventId

				u.Data.DisabledFalse(eventId)
				_, err = u.BotAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "Event enabled."))
				if err != nil {
					return err
				}

				if e.Cron != "once" {
					_, err = u.setCronRepeatable(e.Cron, u.Message.Chat.ID, e.Name, eventId)
					if err != nil {
						Logger.Sugar.Errorln("Couldn't create cron job after enabling event.")
					}
				}
			}
		}
		delete(UserCurrentEvent, u.Message.From.ID)
		return nil
	}

	return nil
}

func (u UserEvent) handleDelete() error {
	var (
		eventsList string
		eventsName map[int]string
	)

	userId := u.Message.From.ID

	eventsName, err := u.Data.GetEventsList(userId)
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
			eventsList += "\n" + event
		}

		_, err = u.BotAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "Which event do you want to delete? \n"+eventsList))
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

		u.Data.DeleteEvent(eventId)
		_, err = u.BotAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "Event deleted"))
		if err != nil {
			return err
		}
		delete(UserCurrentEvent, u.Message.Chat.ID)
		return nil
	}

	return nil
}

func (u UserEvent) handleDeleteAll() error {
	userId := u.Message.From.ID
	db := data.NewData(data.MustConnectPostgres())
	db.DeleteAllEvents(userId)
	_, err := u.BotAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "Your events have been successfully deleted."))
	if err != nil {
		fmt.Println(err)
	}

	return nil
}

func (u UserEvent) handleDefault() error {
	msg := tgbotapi.NewMessage(u.Message.Chat.ID, sendresponse.WelcomeMessage())
	_, err := u.BotAPI.Send(msg)
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func (u UserEvent) HandleOnceReminder() error {
	events, err := u.Data.GetOnceNotFired()
	if err != nil {
		Logger.Sugar.Errorln("No not fired once events")
		return err
	}
	for _, e := range events {
		if time.Now().After(e.TimeDate.In(loc.CurrentLoc)) {
			_, err = u.BotAPI.Send(tgbotapi.NewMessage(e.ChatId, "Don't forget about "+e.Name))
			if err != nil {
				return errors.New(fmt.Sprintf("u.BotAPI.Send failed with error: %s", err))
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
	sched := u.Scheduler.Sc()

	_, err := sched.Every(1).Minute().Do(func() {
		Logger.Sugar.Infoln("HandleOnceReminder fired")
		err := u.HandleOnceReminder()
		if err != nil {
			Logger.Sugar.Errorln("HandleOnceReminder failed", err)
		}
	})
	if err != nil {
		Logger.Sugar.Errorln("Couldn't create job")
	}

	jobs := sched.Jobs()
	if len(jobs) == 1 {
		events, err := u.Data.GetAllEvents()
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

	sched.StartAsync()

	return sched, nil
}
