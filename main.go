package main

import (
	"flag"
	"log"
	"read-tip-bot/clients/telegram"
)

func main() {
	host, token := mustParseFlags()
	_ = telegram.New(host, token)

	// fetcher = fetcher.New()

	// processor = processor.New()

	//consumer.Start(fetcher, processor)
}

func mustParseFlags() (string, string) {
	host := flag.String("h", "", "telegram api host")
	token := flag.String("t", "", "token for access to telegram bot")

	flag.Parse()
	if *token == "" {
		log.Fatal("token is not specified")
	}
	if *host == "" {
		log.Fatal("token is not specified")
	}

	return *host, *token
}
