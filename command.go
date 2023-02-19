package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func interruptibleWait(cmd *exec.Cmd, interrupt os.Signal) error {
	if cmd.Process == nil {
		return fmt.Errorf("interruptible received nil cmd.Process: has Start() been called?")
	}
	if interrupt == nil {
		return fmt.Errorf("interruptible received nil interrupt: a non-nil interrupt is needed to send to the process")
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	defer signal.Stop(c)

	errc := make(chan error, 1)
	go func() {
		signal, ok := <-c // Wait for interrupt
		if !ok {
			return // Channel closed
		}
		err := cmd.Process.Signal(interrupt)
		if err == nil {
			errc <- fmt.Errorf("signal received: %v", signal)
		} else {
			log.Println("Failed to send interrupt signal:", err)
		}
	}()

	waitErr := cmd.Wait()

	var interruptErr error
	select {
	case interruptErr = <-errc:
	default:
	}
	log.Println("interruptErr:", interruptErr)
	log.Println("waitErr:", waitErr)
	if interruptErr != nil {
		return interruptErr
	}
	return waitErr
}

func writeProcessRusageStats(u *syscall.Rusage, f flags) error {
	file, err := os.Create(f.outputVideoProcessStatsPath())
	if err != nil {
		return err
	}

	err = json.NewEncoder(file).Encode(u)
	return err
}
