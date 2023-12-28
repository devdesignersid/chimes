package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/devdesignersid/chimes/pkg/reminder"
	"github.com/sevlyar/go-daemon"
	"github.com/shirou/gopsutil/process"
)

func main() {
	cntxt := &daemon.Context{
		PidFileName: "chimes.pid",
		PidFilePerm: 0644,
		LogFileName: "chimes.log",
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
		Args:        []string{"[chimes-daemon]"},
	}

	d, err := cntxt.Reborn()
	if err != nil {
		fmt.Println("Unable to run: ", err)
	}

	if d != nil {
		return
	}

	defer cntxt.Release()

	fmt.Println("- - - - - - - - - - - - - - -")
	fmt.Println("daemon started")

	go worker()

	err = daemon.ServeSignals()
	if err != nil {
		fmt.Println("Error: ", err)
	}

	// Read the PID from the file
	pidData, err := os.ReadFile("chimes.pid")
	if err != nil {
		fmt.Println("Unable to read PID file:", err)
		return
	}

	// Parse the PID
	pid, err := strconv.Atoi(string(pidData))
	if err != nil {
		fmt.Println("Unable to parse PID:", err)
		return
	}

	// Get the process
	p, err := process.NewProcess(int32(pid))
	if err != nil {
		fmt.Println("Unable to get process:", err)
		return
	}

	// Kill the process
	err = p.Kill()
	if err != nil {
		fmt.Println("Unable to kill process:", err)
		return
	}

	fmt.Println("Process killed successfully")
	fmt.Println("daemon terminated")

}

var (
	stop = make(chan struct{})
)

func worker() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	logFile, err := os.OpenFile("chimes.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer logFile.Close()

	logger := log.New(logFile, "", log.LstdFlags)

	for {
		time.Sleep(1 * time.Second)

		select {
		case killSignal := <-interrupt:
			fmt.Println("Got signal:", killSignal)
			stop <- struct{}{}
			fmt.Println("Worker stopped")
		default:
			logger.Println("Checking for due reminders...")
		}
	}

}

func getSampleData() {
	fmt.Println("Chimes")

	inMemoryReminderStorage := reminder.NewInMemoryReminderStorage()
	reminderService := reminder.NewReminderService(inMemoryReminderStorage)

	currentTime := time.Now()
	futureTime := time.Now().Add(1 * time.Minute)

	reminderService.Save(reminder.CreateReminderData{Message: "Drink water", Due: futureTime, Priority: reminder.Priority(2)})
	reminderService.Save(reminder.CreateReminderData{Message: "Walk away from keyboard", Due: futureTime, Priority: reminder.Priority(1)})
	reminderService.Save(reminder.CreateReminderData{Message: "Attend standup", Due: futureTime, Priority: reminder.Priority(0)})

	sortOrder := reminder.Desc
	orderBy := reminder.ByPriority

	reminders := reminderService.Find(reminder.FilterReminder{DueAfter: &currentTime, SortOrder: &sortOrder, OrderBy: &orderBy})
	for _, reminder := range reminders {
		fmt.Printf("%s, %s, %s\n", reminder.Message, reminder.Due.Format("January 2, 2006, 3:04 PM"), reminder.Priority)
	}
}
