package command

import (
	"errors"
	"eventbot/Logger"
	"eventbot/cron"
	"eventbot/data"
	sendresponse "eventbot/send-response"
	"fmt"
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
		u.Data.CreateEvent(u.Message.From.ID, UserCurrentEvent[u.Message.From.ID].ChatId, UserCurrentEvent[u.Message.From.ID].Name,
			sendresponse.MakeDateTimeField(UserCurrentEvent[u.Message.From.ID].Date,
				UserCurrentEvent[u.Message.From.ID].Time), cron)

		_, err := u.BotAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "Your event has been successfully created."))
		if err != nil {
			fmt.Println(err)
		}
		delete(UserCurrentEvent, u.Message.From.ID)
	}
	u.RerunEvents <- NewEventCommand

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

		_, err := u.BotAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "Which event do you want to edit? \n"+eventsList))
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
		_, err := u.BotAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, sendresponse.AskForName()))
		return err
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

		// send struct
		u.Data.UpdateEvent(UserCurrentEvent[u.Message.From.ID].EventId, UserCurrentEvent[u.Message.From.ID].Name,
			sendresponse.MakeDateTimeField(UserCurrentEvent[u.Message.From.ID].Date,
				UserCurrentEvent[u.Message.From.ID].Time), cron)

		_, err := u.BotAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "Your event has been successfully edited."))
		if err != nil {
			fmt.Println(err)
		}
		delete(UserCurrentEvent, u.Message.From.ID)
	}
	u.RerunEvents <- EditCommand
	return nil
}

func (u UserEvent) handleDisable() error {
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
				fmt.Println(eventId)
			}
		}

		fmt.Println(u.Message.Text, eventId)

		u.Data.DisabledTrue(eventId)
		_, err = u.BotAPI.Send(tgbotapi.NewMessage(userId, "Event disabled"))
		if err != nil {
			return err
		}

		return nil
	}
	u.RerunEvents <- DisableCommand
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

	v, ok := UserCurrentEvent[userId]
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

	if v.CurrentStep == NameStep {
		var eventId int

		for _, e := range eventsName {
			if e.Name == u.Message.Text {
				eventId = e.EventId
				fmt.Println(eventId)
			}
		}

		fmt.Println(u.Message.Text, eventId)

		u.Data.DisabledFalse(eventId)
		_, err = u.BotAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "Event enabled."))
		if err != nil {
			return err
		}

		return nil
	}
	u.RerunEvents <- EnableCommand
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

		return nil
	}

	u.RerunEvents <- DeleteCommand
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

	u.RerunEvents <- DeleteAllCommand
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

func (u UserEvent) HandleCronResponse(chatId int64, name string, eventId int) error {
	events, err := u.Data.GetAllEvents()
	if err != nil {
		Logger.Sugar.Panic(err)
	}

	for _, e := range events {
		if e.Cron != OncePeriod &&
			time.Now().After(e.TimeDate.In(loc.CurrentLoc)) &&
			e.LastFired.Equal(time.Time{}) {

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
	//_, err := u.BotAPI.Send(tgbotapi.NewMessage(chatId, "Don't forget about "+name))
	//err = u.Data.SetLastFired(time.Now(), eventId)
	//if err != nil {
	//	return err
	//}
	//
	//return nil
}
