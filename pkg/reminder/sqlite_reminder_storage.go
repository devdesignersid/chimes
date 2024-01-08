package reminder

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type SqliteReminderStorage struct {
	db *sql.DB
}

var (
	sqliteReminderStorageInstance *SqliteReminderStorage
	sqliteReminderStorageOnce     sync.Once
)

func GetSqliteReminderStorage() (*SqliteReminderStorage, error) {
	var err error

	sqliteReminderStorageOnce.Do(func() {
		sqliteReminderStorageInstance = &SqliteReminderStorage{}
		sqliteReminderStorageInstance.db, err = sql.Open("sqlite3", "./reminders.db")
		if err != nil {
			log.Fatal(err)
			return
		}
		sqliteReminderStorageInstance.db.SetMaxOpenConns(1)
		sqliteReminderStorageInstance.db.SetMaxIdleConns(1)
		sqliteReminderStorageInstance.db.SetConnMaxLifetime(time.Minute)
	})
	if err != nil {
		return nil, err
	}
	return sqliteReminderStorageInstance, nil
}

func CreateTable(storage *SqliteReminderStorage) error {
	statement, err := storage.db.Prepare(`
		CREATE TABLE IF NOT EXISTS reminders (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			message TEXT,
			due DATETIME,
			priority INTEGER,
			created_at DATETIME,
			repeat INTEGER,
			repeat_interval INTEGER
		);
	`)
	if err != nil {
		return err
	}
	_, err = statement.Exec()
	if err != nil {
		return err
	}
	return err
}

func (storage *SqliteReminderStorage) Save(data CreateReminderData) (Reminder, error) {
	reminder := Reminder{
		Message:        data.Message,
		Due:            data.Due,
		Priority:       data.Priority,
		CreatedAt:      time.Now(),
		Repeat:         data.Repeat,
		RepeatInterval: data.RepeatInterval,
	}

	statement, err := storage.db.Prepare(`INSERT INTO reminders
	(message, due, priority, created_at, repeat, repeat_interval) VALUES (?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return Reminder{}, err
	}

	_, err = statement.Exec(
		reminder.Message,
		reminder.Due.Format("2006-01-02 15:04:05"),
		reminder.Priority,
		reminder.CreatedAt.Format("2006-01-02 15:04:05"),
		reminder.Repeat,
		reminder.RepeatInterval,
	)

	if err != nil {
		return Reminder{}, err
	}

	return reminder, nil
}

func (storage *SqliteReminderStorage) Find(filter FilterReminder) []Reminder {
	values := make([]Reminder, 0)
	orderBy := ByCreatedAt
	sortOrder := Desc

	if filter.OrderBy != nil {
		orderBy = *filter.OrderBy
	}
	if filter.SortOrder != nil {
		sortOrder = *filter.SortOrder
	}

	query := `SELECT
	id, message, due, priority, created_at, repeat, repeat_interval FROM reminders WHERE 1=1`

	if filter.Priority != nil {
		query += fmt.Sprintf(" AND priority = %d", *filter.Priority)
	}

	if filter.DueBefore != nil {
		query += fmt.Sprintf(" AND due < '%s'", filter.DueBefore.Format("2006-01-02 15:04:05"))
	}
	if filter.DueAfter != nil {
		query += fmt.Sprintf(" AND due > '%s'", filter.DueAfter.Format("2006-01-02 15:04:05"))
	}
	if filter.DueOn != nil {
		query += fmt.Sprintf(" AND due = '%s'", filter.DueOn.Format("2006-01-02 15:04:05"))
	}

	switch orderBy {
	case "due":
		query += fmt.Sprintf(" ORDER BY due %s", sortOrder)
	case "createdat":
		query += fmt.Sprintf(" ORDER BY created_at %s", sortOrder)
	case "priority":
		query += fmt.Sprintf(" ORDER BY priority %s", sortOrder)
	}

	rows, err := storage.db.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var reminder Reminder
		err = rows.Scan(
			&reminder.Id,
			&reminder.Message,
			&reminder.Due,
			&reminder.Priority,
			&reminder.CreatedAt,
			&reminder.Repeat,
			&reminder.RepeatInterval,
		)
		if err != nil {
			panic(err)
		}
		values = append(values, reminder)
	}
	return values
}

func (storage *SqliteReminderStorage) FindOne(id int) (Reminder, error) {
	var reminder Reminder
	row := storage.db.QueryRow(`SELECT
	 id,
	 message,
	 due,
	 priority,
	 created_at,
	 repeat,
	 repeat_interval
	 FROM reminders WHERE id = ?`, id)

	err := row.Scan(&reminder.Id, &reminder.Message, &reminder.Due, &reminder.Priority, &reminder.CreatedAt, &reminder.Repeat, &reminder.RepeatInterval)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Reminder{}, errors.New("Reminder not found")
		}
		panic(err)
	}

	return reminder, nil
}

func (storage *SqliteReminderStorage) Update(id int, data UpdateReminderData) (Reminder, error) {
	reminder, err := storage.FindOne(id)
	if err != nil {
		return Reminder{}, err
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
	if data.Repeat != nil {
		reminder.Repeat = *data.Repeat
	}
	if data.RepeatInterval != nil {
		reminder.RepeatInterval = *data.RepeatInterval
	}

	statement, err := storage.db.Prepare(`
				UPDATE reminders
				SET 
				message = ?,
				due = ?,
				priority = ?,
				repeat = ?,
				repeat_interval = ?
				WHERE id = ?
	`)
	if err != nil {
		return Reminder{}, err
	}
	_, err = statement.Exec(reminder.Message, reminder.Due.Format("2006-01-02 15:04:05"), reminder.Priority, reminder.Repeat, reminder.RepeatInterval, reminder.Id)
	if err != nil {
		return Reminder{}, err
	}
	return reminder, nil
}

func (storage *SqliteReminderStorage) Delete(id int) (bool, error) {
	_, err := storage.FindOne(id)
	if err != nil {
		return false, err
	}
	statement, err := storage.db.Prepare(`
	DELETE FROM reminders WHERE id = ?`)
	if err != nil {
		return false, err
	}
	_, err = statement.Exec(id)
	if err != nil {
		return false, err
	}

	return true, nil
}
