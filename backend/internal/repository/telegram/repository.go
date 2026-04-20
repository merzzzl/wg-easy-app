package telegram

import "github.com/go-telegram/bot"

type Repository struct {
	bot *bot.Bot
}

func New(botClient *bot.Bot) *Repository {
	return &Repository{bot: botClient}
}
