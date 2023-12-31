package reminder

import (
	"errors"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
)

type InMemoryReminderStorage struct {
	data map[string]Reminder
}

var (
	inMemoryReminderStorageOnce     sync.Once
	inMemoryReminderStorageInstance *InMemoryReminderStorage
)

func GetInMemoryReminderStorage() *InMemoryReminderStorage {
	inMemoryReminderStorageOnce.Do(func() {
		inMemoryReminderStorageInstance = &InMemoryReminderStorage{
			data: make(map[string]Reminder),
		}
	})
	return inMemoryReminderStorageInstance

}

func (storage *InMemoryReminderStorage) Save(data CreateReminderData) (Reminder, error) {
	reminder := Reminder{Id: uuid.NewString(), Message: data.Message, Due: data.Due, Priority: data.Priority, CreatedAt: time.Now(), Repeat: data.Repeat, RepeatInterval: data.RepeatInterval}
	storage.data[reminder.Id] = reminder
	return reminder, nil
}

func (storage *InMemoryReminderStorage) Find(filter FilterReminder) []Reminder {
	values := make([]Reminder, 0, len(storage.data))
	orderBy := ByCreatedAt
	sortOrder := Desc

	if filter.OrderBy != nil {
		orderBy = *filter.OrderBy
	}
	if filter.SortOrder != nil {
		sortOrder = *filter.SortOrder
	}

	for _, value := range storage.data {
		if filter.Priority != nil && value.Priority != *filter.Priority {
			continue
		}
		if filter.DueBefore != nil && value.Due.After(*filter.DueBefore) {
			continue
		}
		if filter.DueAfter != nil && value.Due.Before(*filter.DueAfter) {
			continue
		}
		if filter.DueOn != nil && value.Due.Unix() != (*filter.DueOn).Unix() {
			continue
		}
		values = append(values, value)
	}

	sort.Slice(values, func(i, j int) bool {
		switch orderBy {
		case "due":
			if sortOrder == Asc {
				return values[i].Due.Before(values[j].Due)
			}
			return values[i].Due.After(values[j].Due)
		case "createdat":
			if sortOrder == Asc {
				return values[i].CreatedAt.Before(values[j].CreatedAt)
			}
			return values[i].CreatedAt.After(values[j].CreatedAt)
		case "priority":
			if sortOrder == Asc {
				return values[i].Priority < values[j].Priority
			}
			return values[i].Priority > values[j].Priority
		default:
			return false
		}
	})

	return values
}

func (storage *InMemoryReminderStorage) FindOne(id string) (Reminder, error) {
	value, ok := storage.data[id]
	if !ok {
		return Reminder{}, errors.New("Reminder not found")
	}
	return value, nil
}

func (storage *InMemoryReminderStorage) Update(id string, data UpdateReminderData) (Reminder, error) {
	reminder, error := storage.FindOne(id)
	if error != nil {
		return Reminder{}, error
	}

	if data.Message != nil {
		reminder.Message = *data.Message
	}

	if data.Due != nil {
		reminder.Due = *data.Due
	}

	if data.Priority != nil {
		reminder.Priority = *data.Priority
	}

	storage.data[id] = reminder
	return reminder, nil
}

func (storage *InMemoryReminderStorage) Delete(id string) (bool, error) {
	_, error := storage.FindOne(id)
	if error != nil {
		return false, error
	}
	delete(storage.data, id)
	return true, nil
}
