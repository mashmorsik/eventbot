package cron

import (
	"eventbot/Logger"
	"eventbot/command"
	"eventbot/data"
	"eventbot/pkg/loc"
	"fmt"
	"github.com/go-co-op/gocron"
)

func RerunScheduler(userEvent *command.UserEvent, data *data.Data, sched *gocron.Scheduler) error {
	sched.Stop()

	jobs := sched.Jobs()

	//if err != nil {
	//	return fmt.Errorf("RemoveJob: %s", err)
	//}

	for _, j := range jobs {
		Logger.Sugar.Infof("remove event, name:%s, isRunning:%v", j.GetName(), j.IsRunning())

		err := sched.RemoveByID(j)
		if err != nil {
			return fmt.Errorf("RemoveJob: %s", err)
		}
	}

	_, err := RunScheduler(userEvent, data)
	if err != nil {
		return fmt.Errorf("RunScheduler: %s", err)
	}

	return nil
}

func RunScheduler(userEvent *command.UserEvent, data *data.Data) (*gocron.Scheduler, error) {
	s := gocron.NewScheduler(loc.MskLoc())
	//if err != nil {
	//	return nil, fmt.Errorf("NewScheduler: %s", err)
	//}

	// добавлять тег к джобе с айди ивентом, создавать джоб в createtask функции
	//один крон на все once
	// не вычитывать ивенты, если там есть last-fired
	_, err := s.Every(1).Minute().Do(func() {
		Logger.Sugar.Infoln("HandleOnceReminder fired")
		err := userEvent.HandleOnceReminder()
		if err != nil {
			Logger.Sugar.Errorln("HandleOnceReminder failed", err)
		}
	})
	if err != nil {
		Logger.Sugar.Errorln("Couldn't create job")
	}

	//for _, event := range events {
	//	Logger.Sugar.Infof("set up new cron job: %#v", event)
	//
	//	var (
	//		cronExp            = ""
	//		cronHandler func() = nil
	//	)
	//
	//	switch event.Cron {
	//	case "once":
	//		cronExp = "* * * * *"
	//		cronHandler = func() {
	//			if userEvent.HandleOnceReminder(events) != nil {
	//				Logger.Sugar.Errorln("HandleOnceReminder", err)
	//			}
	//		}
	//	default:
	//		//cronExp = event.Cron
	//		//cronHandler = func() {
	//		//	if userEvent.HandleOnceReminder(events) != nil {
	//		//		Logger.Sugar.Errorln("HandleOnceReminder", err)
	//		//	}
	//		//	//if ue.HandleCronResponse(event.ChatId, event.Name, event.EventId) != nil {
	//		//	//	Logger.Sugar.Errorln(err)
	//		//	//	return
	//		//	//}
	//		//}
	//	}
	//
	//	e, err := s.Cron(cronExp).Do(cronHandler)
	//
	//	if err != nil {
	//		Logger.Sugar.Panic(err)
	//	}

	//Logger.Sugar.Infof("run cron job name: %s", job.GetName())

	s.StartAsync()

	return s, nil
}
