package main

import (
	"flag"
	"log"
	tgClient "read-tip-bot/clients/telegram"
	"read-tip-bot/consumer/event_consumer"
	"read-tip-bot/events/telegram"
	"read-tip-bot/storage/files"
)

const (
	batchSize = 100
)

func main() {
	host, token, storagePath := mustParseFlags()

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

func mustParseFlags() (string, string, string) {
	host := flag.String("h", "", "telegram api host")
	token := flag.String("t", "", "token for access to telegram bot")
	storagePath := flag.String("p", "", "storage path")

	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified")
	}
	if *host == "" {
		log.Fatal("host is not specified")
	}
	if *storagePath == "" {
		log.Fatal("storage path is not specified")
	}

	return *host, *token, *storagePath
}
