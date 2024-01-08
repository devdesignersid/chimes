package main

import (
	"log"
	"time"

	"github.com/devdesignersid/chimes/cmd"
	"github.com/devdesignersid/chimes/pkg/daemon"
	"github.com/devdesignersid/chimes/pkg/reminder"
	"github.com/gen2brain/beeep"
)

func main() {
	cmd.Execute()
	d := daemon.NewDaemon("chimes.pid", "chimes.log", 1*time.Second)
	_, err := d.IsAlive()
	if err != nil {
		p, err := d.Spawn()
		if p != nil {
			return
		}
		if err != nil {
			log.Fatal(err)
		}

		d.Do(job)
	}

}

func job(logger *log.Logger) {
	sqliteReminderStorage, err := reminder.GetSqliteReminderStorage()
	if err != nil {
		panic(err)
	}
	reminderService := reminder.GetReminderService(sqliteReminderStorage)

	dueReminders := reminderService.FindDueReminders()
	for _, dueReminder := range dueReminders {
		err := beeep.Notify("Chimes Reminder", dueReminder.Message, "assets/icon.png")
		if err != nil {
			panic(err)
		}
	}

}
