package main

import (
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"
)

func main() {
	signals := make(chan os.Signal)
	lg := configureLogger()
	lg.Info("startupping ...")

	serv := NewServer(lg)

	go func() {
		err := serv.Start()
		if err != nil {
			lg.WithError(err).Fatal("can't start the server")
		}
	}()

	signal.Notify(signals, os.Kill, os.Interrupt)
	<-signals
	lg.Info("shutting down ...")
}

// configureLogger - Настраивает логгер
func configureLogger() *logrus.Logger {
	lg := logrus.New()
	lg.SetReportCaller(false)
	lg.SetFormatter(&logrus.TextFormatter{})
	lg.SetLevel(logrus.DebugLevel)
	return lg
}
