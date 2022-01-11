package controller

import (
	"gopkg.in/tucnak/telebot.v3"
)

func MakeNewPost(ctx telebot.Context) error {
	return ctx.Send("Новый пост")
}

func EnterAdminMode(ctx telebot.Context) error {
	return ctx.Send("Режим админа")
}

func ShowHelp(ctx telebot.Context) error {
	return ctx.Send("Помощь")
}