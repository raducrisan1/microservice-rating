package main

import (
	"os"
	"os/signal"
	"syscall"
)

func setupsignal() chan os.Signal {
	ossignal := make(chan os.Signal, 1)
	signal.Notify(ossignal,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	return ossignal
}
