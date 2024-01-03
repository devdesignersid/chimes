/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "chimes",
	Short: "Your elegant reminder companion. Built with Golang, this CLI app ensures you never miss a task. Set reminders and get notified on time, every time.",
	Long:  `A sophisticated Command Line Interface (CLI) application, meticulously crafted using the power and simplicity of Golang. Designed with the user in mind, Chimes allows you to set reminders for your important tasks and events. Whether it’s a meeting, a birthday, or a deadline, Chimes ensures you’re always ahead of your schedule.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
