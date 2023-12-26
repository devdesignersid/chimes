package reminder

import "time"

type ReminderService struct {
	storage ReminderStorager
}

func NewReminderService(storage ReminderStorager) *ReminderService {
	return &ReminderService{storage}
}

func (service *ReminderService) Save(data CreateReminderData) (Reminder, error) {
	return service.storage.Save(data)
}

func (service *ReminderService) Find(filter FilterReminder) []Reminder {
	return service.storage.Find(filter)
}

func (service *ReminderService) FindOne(id string) (Reminder, error) {
	return service.storage.FindOne(id)
}

func (service *ReminderService) Update(id string, data UpdateReminderData) (Reminder, error) {
	return service.storage.Update(id, data)
}

func (service *ReminderService) Delete(id string) (bool, error) {
	return service.storage.Delete(id)
}

func (service *ReminderService) FindDueReminders(id string) []Reminder {
	currentTime := time.Now()
	return service.Find(FilterReminder{DueOn: &currentTime})
}
