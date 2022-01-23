package controller

import (
	"github.com/Hickar/sound-seeker-bot/pkg/middleware/scene"
	"gopkg.in/tucnak/telebot.v3"
)

func NewMainRouter(bot *scene.Scene, controller *MainController) {
	bot.HandleStart(controller.ShowMainMenu)
	bot.Handle(MainNewPostCommand, controller.MakeNewPost)
	bot.Handle(MainShowHelpCommand, controller.ShowHelp)
	bot.Handle(MainEnterAdminCommand, controller.EnterAdminMode)
	bot.Handle(telebot.OnText, controller.ShowMainMenu)
}

func NewPostRouter(bot *scene.Scene, controller *PostController) {
	bot.HandleStart(controller.OnPostCreationStart)
	bot.Handle(telebot.OnText, controller.HandlePostAlbumInfo)
	bot.Handle(PostReturnToMainCommand, controller.ExitToMainMenu)
	bot.Handle(PostEnterContentManuallyCommand, controller.EnterPostContentManually)
}
