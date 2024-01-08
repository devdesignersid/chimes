package cmd

import (
	"fmt"

	"github.com/devdesignersid/chimes/pkg/reminder"
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Deletes a reminder from the system",
	Long:  `The 'delete' command is used to delete a reminder from the system. It takes in the 'id' parameter which is the ID of the reminder to be deleted.`,
	Run: func(cmd *cobra.Command, args []string) {
		sqliteReminderStorage, err := reminder.GetSqliteReminderStorage()
		if err != nil {
			panic(err)
		}
		err = reminder.CreateTable(sqliteReminderStorage)
		if err != nil {
			panic(err)
		}
		reminderService := reminder.GetReminderService(sqliteReminderStorage)
		_, err = reminderService.FindOne(id)
		if err != nil {
			fmt.Printf("No reminder found with id: %d", id)
			return
		}

		_, err = reminderService.Delete(id)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Reminder with id %d deleted", id)
	},
}

func init() {
	deleteCmd.Flags().IntVar(&id, "id", 0, "ID of the reminder to be deleted")
	if err := deleteCmd.MarkFlagRequired("id"); err != nil {
		fmt.Println(err)
	}

	rootCmd.AddCommand(deleteCmd)
}
