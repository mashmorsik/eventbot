package command

import (
	"eventbot/Logger"
	cronkafka "eventbot/cron"
	"eventbot/data"
	"time"
)

type EditEvents struct {
	db data.DataInterface
}

type EventInterface interface {
	CreateNewEvent(userId int64, chatId int64, eventName string, timeDate time.Time, cron string) (int, error)
	UpdateEvent(eventId int, name string, timeDate time.Time, cron string) error
	DeleteEvent(eventId int) error
	GetEvents(userId int64) (map[int]*data.Event, error)
	DisableEvent(eventId int) error
	EnableEvent(eventId int) error
	DeleteAllEvents(userId int64) error
}

func (e *EditEvents) CreateNewEvent(userId int64, chatId int64, eventName string, timeDate time.Time, cron string) (int, error) {
	eventId, err := e.db.CreateEvent(userId, chatId, eventName, timeDate, cron)
	if err != nil {
		Logger.Sugar.Errorln("CreateEvent in DB failed.")
	}

	var event = &data.Event{
		UserId:   userId,
		Name:     eventName,
		ChatId:   chatId,
		TimeDate: timeDate,
		Cron:     cron,
		Disabled: false,
	}

	cronkafka.Producer(event)
	return eventId, nil
}

func (e *EditEvents) UpdateEvent(eventId int, name string, timeDate time.Time, cron string) error {
	err := e.db.UpdateEvent(eventId, name, timeDate, cron)
	if err != nil {
		Logger.Sugar.Errorln("UpdateEvent failed")
	}

	return nil
}

func (e *EditEvents) DeleteEvent(eventId int) error {
	err := e.db.DeleteEvent(eventId)
	if err != nil {
		Logger.Sugar.Errorln("DeleteEvent failed")
	}

	return nil
}

func (e *EditEvents) GetEvents(userId int64) (map[int]*data.Event, error) {
	events, err := e.db.GetUserEvents(userId)
	if err != nil {
		Logger.Sugar.Errorln("GetEvents failed")
	}

	return events, err
}

func (e *EditEvents) DisableEvent(eventId int) error {
	err := e.db.DisabledTrue(eventId)
	if err != nil {
		Logger.Sugar.Errorln("DisableEvent failed")
	}

	return nil
}

func (e *EditEvents) EnableEvent(eventId int) error {
	err := e.db.DisabledFalse(eventId)
	if err != nil {
		Logger.Sugar.Errorln("EnableEvent failed")
	}

	return nil
}

func (e *EditEvents) DeleteAllEvents(userId int64) error {
	err := e.db.DeleteAllEvents(userId)
	if err != nil {
		Logger.Sugar.Errorln("DeleteAllEvents failed")
	}

	return nil
}
