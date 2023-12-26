package reminder

import "time"

type Priority int

const (
	Low Priority = iota
	Medium
	High
)

func (rp Priority) String() string {
	switch rp {
	case Low:
		return "Low"
	case Medium:
		return "Medium"
	case High:
		return "High"
	default:
		panic("Unhandled Priority!")
	}
}

type Reminder struct {
	Id        string
	Message   string
	Due       time.Time
	Priority  Priority
	CreatedAt time.Time
}

type CreateReminderData struct {
	Message  string
	Due      time.Time
	Priority Priority
}

type UpdateReminderData struct {
	Message  *string
	Due      *time.Time
	Priority *Priority
}

type SortOrderValue string

const (
	Asc  SortOrderValue = "asc"
	Desc SortOrderValue = "desc"
)

type OrderByField string

const (
	ByPriority  OrderByField = "priority"
	ByDue       OrderByField = "due"
	ByCreatedAt OrderByField = "createdat"
)

type FilterReminder struct {
	Priority  *Priority
	DueBefore *time.Time
	DueAfter  *time.Time
	DueOn     *time.Time
	OrderBy   *OrderByField
	SortOrder *SortOrderValue
}
