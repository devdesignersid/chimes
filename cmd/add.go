package cmd

import (
	"fmt"
	"time"

	"github.com/devdesignersid/chimes/pkg/reminder"
	"github.com/spf13/cobra"
)

var (
	id             int
	memo           string
	days           int
	months         int
	years          int
	hours          int
	minutes        int
	seconds        int
	date           string
	priority       int
	repeat         bool
	repeatInterval int
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a reminder to the system",
	Long:  `The 'add' command is used to add a reminder to the system. It takes in various parameters such as 'memo', 'priority', 'date', 'days', 'months', 'years', 'hours', 'minutes', 'seconds', 'repeat', and 'repeatInterval'. The 'memo' parameter is used to specify the message to be reminded. The 'priority' parameter is used to specify the priority of the reminder, with 0 being the lowest and 2 being the highest. The 'date', 'days', 'months', 'years', 'hours', 'minutes', and 'seconds' parameters are used to specify the due date of the reminder. The 'repeat' parameter is used to specify whether the reminder should repeat or not. The 'repeatInterval' parameter is used to specify the interval in seconds between each repeat.`,
	Run: func(cmd *cobra.Command, args []string) {
		var futureTime time.Time
		var err error

		if date != "" {
			futureTime, err = time.Parse("2006-01-02 15:04:05", date)
			if err != nil {
				panic(err)
			}
		} else {
			futureTime = time.Now()
			if years < 0 {
				panic("years must be a positive number")
			}
			if months < 0 || months > 12 {
				panic("months must be a positive number less than or equal to 12")
			}
			if days < 0 || days > 31 {
				panic("days must be a positive number greater than or equal to 31")
			}
			futureTime = futureTime.AddDate(years, months, days)
			futureTime = futureTime.Add(time.Duration(hours) * time.Hour)
			futureTime = futureTime.Add(time.Duration(minutes) * time.Minute)
			futureTime = futureTime.Add(time.Duration(seconds) * time.Second)
		}

		if priority < 0 || priority > 2 {
			panic("priority must be a positive number less than or equal to 2")
		}

		sqliteReminderStorage, err := reminder.GetSqliteReminderStorage()
		if err != nil {
			panic(err)
		}
		err = reminder.CreateTable(sqliteReminderStorage)
		if err != nil {
			panic(err)
		}
		reminderService := reminder.GetReminderService(sqliteReminderStorage)

		r, err := reminderService.Save(reminder.CreateReminderData{Message: memo, Due: futureTime, Priority: reminder.Priority(priority), Repeat: repeat, RepeatInterval: time.Duration(time.Duration(repeatInterval) * time.Second)})
		if err != nil {
			panic(err)
		}
		fmt.Printf(`"%s" added to reminders`, r.Message)
	},
}

func init() {
	addCmd.Flags().StringVarP(&memo, "memo", "m", "", "Memo to be reminded")
	if err := addCmd.MarkFlagRequired("memo"); err != nil {
		fmt.Println(err)
	}
	addCmd.Flags().StringVar(&date, "date", "", "Due date in format '2006-01-02 15:04:05' (optional)")
	addCmd.Flags().IntVar(&days, "days", 0, "Number of days to add to current time")
	addCmd.Flags().IntVar(&months, "months", 0, "Number of months to add to current time")
	addCmd.Flags().IntVar(&years, "years", 0, "Number of years to add to current time")
	addCmd.Flags().IntVar(&minutes, "minutes", 0, "Number of minutes to add to current time")
	addCmd.Flags().IntVar(&seconds, "seconds", 0, "Number of seconds to add to current time")
	addCmd.Flags().IntVar(&hours, "hours", 0, "Number of hours to add to current time")

	addCmd.Flags().IntVarP(&priority, "priority", "p", 0, "Priority of the reminder 0 being the lowest and 2 being the highest")
	if err := addCmd.MarkFlagRequired("priority"); err != nil {
		fmt.Println(err)
	}

	addCmd.Flags().BoolVarP(&repeat, "repeat", "r", false, "Set to true if the reminder should repeat")
	addCmd.Flags().IntVar(&repeatInterval, "repeat-interval", 0, "Interval in seconds between each repeat")

	addCmd.PreRun = func(cmd *cobra.Command, args []string) {
		if date != "" {
			days = 0
			months = 0
			years = 0
			minutes = 0
			seconds = 0
			hours = 0
		}
	}

	rootCmd.AddCommand(addCmd)
}
