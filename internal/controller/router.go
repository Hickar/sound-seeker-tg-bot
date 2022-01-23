package controller

import (
	"github.com/Hickar/sound-seeker-bot/pkg/middleware/scene"
	"gopkg.in/tucnak/telebot.v3"
)

func NewMainRouter(bot *scene.Scene) {
	bot.HandleStart(ShowMainMenu)

	bot.Handle(MainNewPostCommand, MakeNewPost)
	bot.Handle(MainShowHelpCommand, ShowHelp)
	bot.Handle(MainEnterAdminCommand, EnterAdminMode)
	bot.Handle(telebot.OnText, ShowMainMenu)
}

func NewPostRouter(bot *scene.Scene, controller *PostController) {
	bot.HandleStart(controller.OnPostCreationStart)
	bot.Handle(telebot.OnText, controller.HandlePostAlbumInfo)
	bot.Handle(PostReturnToMainCommand, controller.ExitToMainMenu)
	bot.Handle(PostEnterContentManuallyCommand, controller.EnterPostContentManually)
}
