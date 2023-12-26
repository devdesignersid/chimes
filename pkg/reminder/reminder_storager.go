package reminder

type ReminderStorager interface {
	Save(CreateReminderData) (Reminder, error)
	Find(filter FilterReminder) []Reminder
	FindOne(id string) (Reminder, error)
	Update(id string, data UpdateReminderData) (Reminder, error)
	Delete(id string) (bool, error)
}
