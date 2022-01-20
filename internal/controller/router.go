package controller

import (
	"fmt"

	"github.com/Hickar/sound-seeker-bot/internal/middleware/scene"
	"gopkg.in/tucnak/telebot.v3"
)

func NewMainRouter(bot *scene.Scene) error {
	bot.HandleStart(ShowMainMenu)

	bot.Handle("Создать пост", MakeNewPost)
	bot.Handle("Помощь", ShowHelp)
	bot.Handle("Режим админа", EnterAdminMode)

	bot.Handle(telebot.OnText, func(ctx telebot.Context) error {
		return ctx.Send("main")
	})

	return nil
}

func NewPostRouter(bot *scene.Scene) error {
	bot.Handle(telebot.OnText, func(ctx telebot.Context) error {
		msg := fmt.Sprintf("[post]: %s", ctx.Message().Text)
		return ctx.Send(msg)
	})

	bot.HandleStart(func (ctx telebot.Context) error {
		return ctx.Send("Post scene start")
	})

	return nil
}