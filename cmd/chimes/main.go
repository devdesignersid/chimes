package main

import (
	"fmt"
	"time"

	"github.com/devdesignersid/chimes/pkg/reminder"
)

func main() {
	fmt.Println("Chimes")

	inMemoryReminderStorage := reminder.NewInMemoryReminderStorage()
	reminderService := reminder.NewReminderService(inMemoryReminderStorage)

	currentTime := time.Now()
	futureTime := time.Now().Add(1 * time.Minute)

	reminderService.Save(reminder.CreateReminderData{Message: "Drink water", Due: futureTime, Priority: reminder.Priority(2)})
	reminderService.Save(reminder.CreateReminderData{Message: "Walk away from keyboard", Due: futureTime, Priority: reminder.Priority(1)})
	reminderService.Save(reminder.CreateReminderData{Message: "Attend standup", Due: futureTime, Priority: reminder.Priority(0)})

	sortOrder := reminder.Desc
	orderBy := reminder.ByPriority

	reminders := reminderService.Find(reminder.FilterReminder{DueAfter: &currentTime, SortOrder: &sortOrder, OrderBy: &orderBy})
	for _, reminder := range reminders {
		fmt.Printf("%s, %s, %s\n", reminder.Message, reminder.Due.Format("January 2, 2006, 3:04 PM"), reminder.Priority)
	}

}
