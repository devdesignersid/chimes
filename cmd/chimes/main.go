package main

import (
	"fmt"
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
	} else {
		d.Kill()
	}

}

func job(logger *log.Logger) {
	logger.Println("Checking for due reminder...")
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

func seedData() {
	inMemoryReminderStorage := reminder.GetInMemoryReminderStorage()
	reminderService := reminder.GetReminderService(inMemoryReminderStorage)
	futureTime := time.Now().Add(3 * time.Second)

	reminderService.Save(reminder.CreateReminderData{Message: "Drink water", Due: futureTime, Priority: reminder.Priority(2), Repeat: true, RepeatInterval: 3 * time.Second})

	reminders := reminderService.Find(reminder.FilterReminder{})
	fmt.Printf("%#v", reminders)

}
