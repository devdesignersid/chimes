package cmd

import (
	"fmt"
	"time"

	"github.com/devdesignersid/chimes/pkg/reminder"
	"github.com/spf13/cobra"
)

var (
	wasDue    bool
	willBeDue bool
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		reminderFilter := reminder.FilterReminder{}

		sqliteReminderStorage, err := reminder.GetSqliteReminderStorage()
		if err != nil {
			panic(err)
		}
		err = reminder.CreateTable(sqliteReminderStorage)
		if err != nil {
			panic(err)
		}
		reminderService := reminder.GetReminderService(sqliteReminderStorage)
		currentTime := time.Now()

		if wasDue {
			reminderFilter.DueBefore = &currentTime
		}
		if willBeDue {
			reminderFilter.DueAfter = &currentTime
		}

		dueReminders := reminderService.Find(reminderFilter)

		for i, dueReminder := range dueReminders {
			fmt.Printf("%d. %s - %s\n", i+1, dueReminder.Message, dueReminder.Due.Format("02 January 2006 03:04:05 PM"))
		}

	},
}

func init() {
	listCmd.Flags().BoolVar(&wasDue, "was-due", false, "filter reminders that were due in the past")
	listCmd.Flags().BoolVar(&willBeDue, "will-be-due", false, "filter reminders that will be due in the future")

	rootCmd.AddCommand(listCmd)
}
