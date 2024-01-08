package reminder

import (
	"sync"
	"time"
)

type ReminderService struct {
	storage ReminderStorager
}

var (
	reminderServiceOnce     sync.Once
	reminderServiceInstance *ReminderService
)

func GetReminderService(storage ReminderStorager) *ReminderService {
	reminderServiceOnce.Do(func() {
		reminderServiceInstance = &ReminderService{storage}
	})
	return reminderServiceInstance
}

func (service *ReminderService) Save(data CreateReminderData) (Reminder, error) {
	return service.storage.Save(data)
}

func (service *ReminderService) Find(filter FilterReminder) []Reminder {
	return service.storage.Find(filter)
}

func (service *ReminderService) FindOne(id int) (Reminder, error) {
	return service.storage.FindOne(id)
}

func (service *ReminderService) Update(id int, data UpdateReminderData) (Reminder, error) {
	return service.storage.Update(id, data)
}

func (service *ReminderService) Delete(id int) (bool, error) {
	return service.storage.Delete(id)
}

func (service *ReminderService) FindDueReminders() []Reminder {
	currentTime := time.Now()
	dueReminders := service.Find(FilterReminder{DueOn: &currentTime})

	for _, reminder := range dueReminders {
		if reminder.Repeat {
			updatedDue := reminder.Due.Add(reminder.RepeatInterval)
			service.Update(reminder.Id, UpdateReminderData{Due: &updatedDue})
		}
	}

	return dueReminders
}
