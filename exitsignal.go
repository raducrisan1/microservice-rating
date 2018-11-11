package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"
)

func setupsignal(ticker <-chan time.Time) chan int {
	ossignal := make(chan os.Signal, 1)
	pulse := make(chan int)
	signal.Notify(ossignal,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		<-ossignal
		pulse <- 1
	}()
	go func() {
		<-ticker
		pulse <- 0
	}()
	return pulse
}
