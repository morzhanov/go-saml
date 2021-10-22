package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/morzhanov/go-saml/internal/idp"
	"github.com/morzhanov/go-saml/internal/idp/config"
	"github.com/morzhanov/go-saml/internal/idp/store"
	"github.com/morzhanov/go-saml/internal/logger"
	"github.com/morzhanov/go-saml/internal/psql"
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
	c, err := config.NewConfig()
	failOnError(l, "config", err)

	db, err := psql.NewDb(c.PSQLurl)
	failOnError(l, "PSQL", err)

	s := store.NewIDPStore(db)
	idpLog := logger.NewIDPLogger(l)
	i, err := idp.NewServer(idpLog, s, c.BaseURL, c.Key, c.Cert)
	failOnError(l, "service", err)

	err = i.Listen()
	failOnError(l, "listen", err)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	log.Println("App successfully started!")
	<-quit
	log.Println("received os.Interrupt, exiting...")
}
