package controller

import (
	"errors"
	"fmt"

	"github.com/Hickar/sound-seeker-bot/internal/usecase"
	"github.com/Hickar/sound-seeker-bot/pkg/middleware/scene"
	"github.com/Hickar/sound-seeker-bot/pkg/middleware/session"
	"gopkg.in/tucnak/telebot.v3"
)

const (
	ScenePostName                   = "post"
	PostReturnToMainCommand         = "❌ Вернуться в главное меню"
	PostEnterContentManuallyCommand = "✏️ Ввести текст поста вручную"

	PostArtistReply            = "Введи название исполнителя и альбома, либо скинь ссылку на альбом в спотифае"
	PostEnterManuallyReply     = "Введи текст поста. (можно даже приложить картинки!)"
	PostInvalidSpotifyURLReply = "Некорректная ссылка на альбом в спотифае!"
	PostNotFoundReply          = "По данному запросу ничего не было найдено"

	scenePostStateKey                  = "post_state"
	scenePostStateEnterContentManually = "post_state_enter_manually"
	scenePostStateWaitingAlbumInfo     = "post_state_album_user_prompt"
	scenePostStateWaitingDescription   = "post_state_description_user_prompt"
)

type PostController struct {
	useCase *usecase.PostUsecase
	stage   *scene.Stage
}

func NewChannelPostController(useCase *usecase.PostUsecase, stage *scene.Stage) *PostController {
	return &PostController{useCase: useCase, stage: stage}
}

func (pc *PostController) OnPostCreationStart(ctx telebot.Context) error {
	var (
		menu                    = telebot.ReplyMarkup{ResizeKeyboard: true, ForceReply: true}
		cancelBtn               = menu.Text(PostReturnToMainCommand)
		enterInfoManuallyButton = menu.Text(PostEnterContentManuallyCommand)
	)

	menu.Reply(
		menu.Row(enterInfoManuallyButton),
		menu.Row(cancelBtn),
	)

	ssn := session.GetSession(ctx)
	ssn.Set(scenePostStateKey, scenePostStateWaitingAlbumInfo)

	return ctx.Send(PostArtistReply, &menu)
}

func (pc *PostController) HandlePostAlbumInfo(ctx telebot.Context) error {
	ssn := session.GetSession(ctx)
	if ssn != nil {
		ssnState, ok := ssn.Get(scenePostStateKey).(string)
		if !ok {
			return errors.New("post scene state is not set")
		}

		msg := ctx.Message().Text

		switch ssnState {
		case scenePostStateWaitingAlbumInfo:
			results, err := pc.useCase.FindAlbums(msg)
			if err != nil {
				if errors.Is(err, usecase.ErrInvalidSpotifyURL) {
					return ctx.Send(PostInvalidSpotifyURLReply)
				} else {
					return ctx.Send(PostNotFoundReply)
				}
			}

			for _, result := range results {
				artistName := "None"
				if len(result.Artists) > 0 {
					artistName = result.Artists[0]
				}

				ctx.Send(fmt.Sprintf("Artist: %s\nAlbum: %s\nYear: %s\nCountry: %s\nSpotify: %s", artistName, result.Title, result.Year, result.Country, result.SpotifyLink))
			}

			return nil
		case scenePostStateWaitingDescription:
			return ctx.Send("waiting description")
		case scenePostStateEnterContentManually:
			return ctx.Send("entering post content manually")
		}
	}

	return errors.New("session is nil")
}

func (pc *PostController) EnterPostContentManually(ctx telebot.Context) error {
	var (
		menu      = telebot.ReplyMarkup{ResizeKeyboard: true, ForceReply: true}
		cancelBtn = menu.Text(PostReturnToMainCommand)
	)

	menu.Reply(menu.Row(cancelBtn))

	ssn := session.GetSession(ctx)
	ssn.Set(scenePostStateKey, scenePostStateEnterContentManually)

	return ctx.Send(PostEnterManuallyReply, &menu)
}

func (pc *PostController) ExitToMainMenu(ctx telebot.Context) error {
	session.GetSession(ctx).Delete(scenePostStateKey)
	return pc.stage.EnterScene(ctx, SceneMainName)
}
