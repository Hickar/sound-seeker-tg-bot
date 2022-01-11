package controller

import (
	"github.com/Hickar/sound-seeker-bot/internal/usecase"
	"gopkg.in/tucnak/telebot.v3"
)

type PostController struct {
	Usecase usecase.PostUsecase
}

func NewChannelPostController(usecase usecase.PostUsecase) *PostController {
	return &PostController{usecase}
}

func (c *PostController) FindPostAlbumInfo(ctx telebot.Context) error {
	return nil
}

func (c *PostController) EnterPostDetails(ctx telebot.Context) error {
	return nil
}
