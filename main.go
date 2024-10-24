package main

import (
	"log"
	"os"
	tgClient "read-tip-bot/clients/telegram"
	"read-tip-bot/consumer/event_consumer"
	"read-tip-bot/events/telegram"
	"read-tip-bot/storage/files"
)

const (
	batchSize = 100
	host      = "api.telegram.org"
)

func main() {
	token, storagePath := mustVariables()

	eventsProcessor := telegram.New(
		tgClient.New(host, token),
		files.New(storagePath),
	)

	log.Print("service started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)
	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}

func mustVariables() (string, string) {
	token := os.Getenv("TELEGRAM_API_KEY")
	storagePath := os.Getenv("READ_TIP_BOT_STORAGE")

	if token == "" {
		log.Fatal("TELEGRAM_API_KEY is not specified")
	}
	if storagePath == "" {
		log.Fatal("READ_TIP_BOT_STORAGE is not specified")
	}

	return token, storagePath
}
