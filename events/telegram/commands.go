package telegram

import (
	"errors"
	"log"
	"net/url"
	"read-tip-bot/clients/telegram"
	"read-tip-bot/lib/e"
	"read-tip-bot/storage"
	"strings"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func NewMessageSender(chatID int, tg *telegram.Client) func(string) error {
	return func(msg string) error {
		return tg.SendMessage(chatID, msg)
	}
}

func (p *Processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%s'", text, username)

	if isAddCmd(text) {
		return p.savePage(chatID, text, username)
	}

	switch text {
	case RndCmd:
		return p.sendRandom(chatID, username)
	case HelpCmd:
		return p.sendHelp(chatID)
	case StartCmd:
		return p.sendHello(chatID)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}
}

func (p *Processor) savePage(chatID int, pageURL string, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do cmd save page", err) }()
	sendMsg := NewMessageSender(chatID, p.tg)

	page := &storage.Page{
		URL:      pageURL,
		UserName: username,
	}

	isExists, err := p.storage.IsExists(page)
	if err != nil {
		return err
	}
	if isExists {
		return sendMsg(msgAlreadyExists)
	}

	if err := p.storage.Save(page); err != nil {
		return err
	}

	return sendMsg(msgSaved)
}

func (p *Processor) sendRandom(chatID int, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do cmd send random", err) }()
	sendMsg := NewMessageSender(chatID, p.tg)

	page, err := p.storage.PickRandom(username)
	if errors.Is(err, storage.ErrNoSavedPages) {
		return sendMsg(msgNoSavedPages)
	}
	if err != nil {
		return err
	}

	if err := sendMsg(page.URL); err != nil {
		return err
	}

	return p.storage.Remove(page)
}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *Processor) sendHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	u, err := url.Parse(text)
	return err == nil && u.Host != ""
}
