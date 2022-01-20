package controller

import (
	"github.com/Hickar/sound-seeker-bot/internal/middleware/scene"
	"gopkg.in/tucnak/telebot.v3"
)

const (
	SceneMain = "main"
	MainNewPostCommand = "Создать новый пост"
	MainEnterAdminCommand = "Режим админа"
	MainShowHelpCommand = "Помощь"
	MainGreeteing = "Главное меню"
)

func ShowMainMenu(ctx telebot.Context) error {
	return ctx.Send(MainGreeteing)
}

func MakeNewPost(ctx telebot.Context) error {
	return scene.GetScene(ctx).EnterScene(ctx, "post")
}

func EnterAdminMode(ctx telebot.Context) error {
	return ctx.Send("Режим админа")
}

func ShowHelp(ctx telebot.Context) error {
	return ctx.Send("Бот для предложения музыки в личный канал")
}