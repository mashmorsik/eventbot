package cron

import (
	"eventbot/pkg/loc"
	"github.com/go-co-op/gocron"
)

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

//func RunScheduler(userEvent *command.UserEvent) (*gocron.Scheduler, error) {
//	s := userEvent.Scheduler.Sc()
//
//	_, err := s.Every(1).Minute().Do(func() {
//		Logger.Sugar.Infoln("HandleOnceReminder fired")
//		err := userEvent.HandleOnceReminder()
//		if err != nil {
//			Logger.Sugar.Errorln("HandleOnceReminder failed", err)
//		}
//	})
//	if err != nil {
//		Logger.Sugar.Errorln("Couldn't create job")
//	}
//
//	s.StartAsync()
//
//	return s, nil
//}
