package daemon

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/sevlyar/go-daemon"
	"github.com/shirou/gopsutil/process"
)

type Daemon struct {
	PidFileName string
	LogFileName string
	Wait        time.Duration
	Stop        chan struct{}
}

func NewDaemon(pidFileName string, logFileName string, wait time.Duration) *Daemon {
	return &Daemon{
		PidFileName: pidFileName,
		LogFileName: logFileName,
		Stop:        make(chan struct{}),
		Wait:        wait,
	}
}

func (d *Daemon) Spawn() (*os.Process, error) {
	context := daemon.Context{
		PidFileName: d.PidFileName,
		PidFilePerm: 0644,
		LogFileName: d.LogFileName,
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
		Args:        []string{"[]"},
	}

	child, err := context.Reborn()
	if err != nil {
		return nil, err
	}

	return child, nil
}

func (d *Daemon) Kill() error {
	p, err := d.IsAlive()
	if err != nil {
		return err
	}
	err = p.Kill()
	if err != nil {
		return err
	}
	return nil
}

func (d *Daemon) Do(job func(*log.Logger)) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	logFile, err := os.OpenFile(d.LogFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer logFile.Close()
	logger := log.New(logFile, "", log.LstdFlags)

	for {
		time.Sleep(d.Wait)
		select {
		case killSignal := <-interrupt:
			logger.Println("Got signal: ", killSignal)
			d.Stop <- struct{}{}
			logger.Println("Worker stopped")
		default:
			job(logger)
		}
	}
}

func (d *Daemon) IsAlive() (*process.Process, error) {
	pidData, err := os.ReadFile(d.PidFileName)
	if err != nil {
		return nil, err
	}

	pid, err := strconv.Atoi(string(pidData))
	if err != nil {
		return nil, err
	}

	p, err := process.NewProcess(int32(pid))
	if err != nil {
		return nil, err
	}

	return p, nil
}
