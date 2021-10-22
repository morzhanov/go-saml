package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/morzhanov/go-saml/internal/logger"
	"github.com/morzhanov/go-saml/internal/sp"
	"go.uber.org/zap"
)

func failOnError(l *zap.Logger, step string, err error) {
	if err != nil {
		l.Fatal("initialization error", zap.Error(err), zap.String("step", step))
	}
}

func main() {
	l, err := logger.NewLogger()
	if err != nil {
		log.Fatal("initialization error during logger setup")
	}
	s := sp.NewClient(l)
	err = s.Listen()
	failOnError(l, "listen", err)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	log.Println("App successfully started!")
	<-quit
	log.Println("received os.Interrupt, exiting...")
}
