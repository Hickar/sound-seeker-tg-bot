package controller

import (
	"github.com/Hickar/sound-seeker-bot/pkg/middleware/scene"
	"gopkg.in/tucnak/telebot.v3"
)

const (
	SceneMainName         = "main"
	MainNewPostCommand    = "📝 Создать новый пост"
	MainEnterAdminCommand = "🔑 Режим админа"
	MainShowHelpCommand   = "📖 Помощь"

	MainStartReply = "Главное меню"
	MainEnterAdminReply = "Администраторский режим в процессе разработки..."
	MainShowHelpReply = "Супер-бот, созданный для обработки постов, предложенных пользователями, и для администрирования канала \"Sonus Emporium\""
)

func ShowMainMenu(ctx telebot.Context) error {
	var (
		menu              = telebot.ReplyMarkup{ResizeKeyboard: true, ForceReply: true}
		createPostBtn     = menu.Text(MainNewPostCommand)
		enterAdminModeBtn = menu.Text(MainEnterAdminCommand)
		showHelpBtn       = menu.Text(MainShowHelpCommand)
	)

	menu.Reply(
		menu.Row(createPostBtn),
		menu.Row(enterAdminModeBtn),
		menu.Row(showHelpBtn),
	)

	return ctx.Send(MainStartReply, &menu)
}

func MakeNewPost(ctx telebot.Context) error {
	return scene.GetScene(ctx).EnterScene(ctx, ScenePostName)
}

func EnterAdminMode(ctx telebot.Context) error {
	return ctx.Send(MainEnterAdminReply)
}

func ShowHelp(ctx telebot.Context) error {
	return ctx.Send(MainShowHelpReply)
}