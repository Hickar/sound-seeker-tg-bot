package controller

import (
	"github.com/Hickar/sound-seeker-bot/internal/usecase"
	"gopkg.in/tucnak/telebot.v3"
)

const ScenePost = "post"

type PostController struct {
	Usecase usecase.PostUsecase
}

func NewChannelPostController(usecase usecase.PostUsecase) *PostController {
	return &PostController{usecase}
}

func OnPostCreationStart(ctx telebot.Context) error {
	menu := telebot.ReplyMarkup{ResizeKeyboard: true}
	menu.Reply(
		menu.Row(menu.Text("Назад в главное меню")),
	)

	return ctx.Send("Добро пожаловать в режим администратора", menu)
}

func (c *PostController) FindPostAlbumInfo(ctx telebot.Context) error {
	return nil
}

func (c *PostController) EnterPostDetails(ctx telebot.Context) error {
	return nil
}
