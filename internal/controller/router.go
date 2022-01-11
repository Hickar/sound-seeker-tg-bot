package controller

import (
	"gopkg.in/tucnak/telebot.v3"
)

func NewMainRouter(bot *telebot.Bot) error {
	commands := []telebot.Command{
		{
			Text:        "/newpost",
			Description: "Создать новый пост в канале",
		},
		{
			Text:        "/help",
			Description: "Справка по боту",
		},
		{
			Text:        "/admin",
			Description: "Перейти в админский режим",
		},
	}

	//mainScene := telebot.NewScene("main")

	bot.Handle("/newpost", MakeNewPost)
	bot.Handle("/help", ShowHelp)
	bot.Handle("/admin", EnterAdminMode)

	if err := bot.SetCommands(commands); err != nil {
		return err
	}

	return nil
}
