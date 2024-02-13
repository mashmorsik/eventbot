package cronkafka

import (
	"errors"
	"eventbot/Logger"
	"eventbot/data"
	"eventbot/pkg/loc"
	"fmt"
	"github.com/go-co-op/gocron"
	"time"
)

type DB struct {
	db data.DataInterface
}

type Scheduler struct {
	sched *gocron.Scheduler
}

func NewScheduler(scheduler *gocron.Scheduler) *Scheduler {
	return &Scheduler{sched: scheduler}
}

func (s *Scheduler) Sc() *gocron.Scheduler {
	return s.sched
}

func StartScheduler() *gocron.Scheduler {
	sched := gocron.NewScheduler(loc.MskLoc())
	return sched
}

func (d *DB) RunScheduler() (*gocron.Scheduler, error) {
	sc := NewScheduler(StartScheduler())
	scheduler := sc.Sc()
	defer scheduler.StartAsync()

	_, err := scheduler.Every(1).Minute().Do(func() {
		Logger.Sugar.Infoln("HandleOnceReminder fired")
		err := d.HandleOnceReminder()
		if err != nil {
			Logger.Sugar.Errorln("HandleOnceReminder failed", err)
		}
	})
	if err != nil {
		Logger.Sugar.Errorln("Couldn't create job")
	}

	jobs := scheduler.Jobs()
	if len(jobs) == 1 {
		events, err := d.db.GetCronMultipleActive()
		if err != nil {
			// if errors.Is(err, norows)
			Logger.Sugar.Errorf("Couldn't get events from DB, err: %v", err)
		}

		for _, e := range events {
			_, err = d.setCronRepeatable(e.Cron, e.ChatId, e.Name, e.EventId)
			if err != nil {
				Logger.Sugar.Errorln("Couldn't add cron jobs from DB.")
			}
		}
	}

	return scheduler, nil
}

func (d *DB) setCronRepeatable(cron string, chatID int64, eventName string, eventId int) (*gocron.Scheduler, error) {
	s, err := d.RunScheduler()
	if err != nil {
		Logger.Sugar.Errorln("Couldn't get scheduler")
		return nil, err
	}

	s.TagsUnique()

	_, err = s.Cron(cron).Tag(string(rune(eventId))).Do(func() error {
		// TODO how to send messages from the bot?
		//if err = u.SendMessage(chatID, "Don't forget about "+eventName); err != nil {
		//	return err
		//}

		err = d.db.SetLastFired(time.Now(), eventId)
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

func (d *DB) HandleOnceReminder() error {
	events, err := d.db.GetOnceNotFired()
	if err != nil {
		Logger.Sugar.Errorln("No not fired once events")
		return err
	}
	for _, e := range events {
		if time.Now().After(e.TimeDate.In(loc.CurrentLoc)) {
			// TODO how to send messages from the bot?
			//if err = u.SendMessage(e.ChatId, "Don't forget about "+e.Name); err != nil {
			//	return err
			//}
			err = d.db.SetLastFired(time.Now(), e.EventId)
			if err != nil {
				return errors.New(fmt.Sprintf("u.Data.SetLastFired failed with error: %s", err))
			}
		}
	}
	return nil
}
