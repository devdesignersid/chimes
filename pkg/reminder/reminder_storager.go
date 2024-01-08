package reminder

type ReminderStorager interface {
	Save(CreateReminderData) (Reminder, error)
	Find(filter FilterReminder) []Reminder
	FindOne(id int) (Reminder, error)
	Update(id int, data UpdateReminderData) (Reminder, error)
	Delete(id int) (bool, error)
}
