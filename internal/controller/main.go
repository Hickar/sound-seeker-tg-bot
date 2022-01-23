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

	MainStartReply      = "Главное меню"
	MainEnterAdminReply = "Администраторский режим в процессе разработки..."
	MainShowHelpReply   = "Супер-бот, созданный для обработки постов, предложенных пользователями, и для администрирования канала \"Sonus Emporium\""
)

type MainController struct {
	stage *scene.Stage
}

func NewMainController(stage *scene.Stage) *MainController {
	return &MainController{stage: stage}
}

func (mc *MainController) ShowMainMenu(ctx telebot.Context) error {
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

func (mc *MainController) MakeNewPost(ctx telebot.Context) error {
	return mc.stage.EnterScene(ctx, ScenePostName)
}

func (mc *MainController) EnterAdminMode(ctx telebot.Context) error {
	return ctx.Send(MainEnterAdminReply)
}

func (mc *MainController) ShowHelp(ctx telebot.Context) error {
	return ctx.Send(MainShowHelpReply)
}
