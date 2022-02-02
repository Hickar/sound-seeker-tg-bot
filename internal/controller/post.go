package controller

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Hickar/sound-seeker-bot/internal/entity"
	"github.com/Hickar/sound-seeker-bot/internal/usecase"
	"github.com/Hickar/sound-seeker-bot/pkg/middleware/scene"
	"github.com/Hickar/sound-seeker-bot/pkg/middleware/session"
	"gopkg.in/tucnak/telebot.v3"
)

const (
	ScenePostName                   = "post"
	PostReturnToMainCommand         = "❌ Вернуться в главное меню"
	PostEnterContentManuallyCommand = "✏️ Ввести текст поста вручную"
	PostFoundOneAcceptCommand       = "👍 Да, это он"
	PostFoundOneDeclineCommand      = "👎 Нет, не он"
	PostConfirmSendingCommand       = "📪 Точно-точно"
	PostCancelSendingCommand        = "😓 Не, всё-таки чуть переделаю"

	PostArtistReply                  = "Введи название исполнителя и альбома, либо скинь ссылку на альбом в спотифае"
	PostEnterManuallyReply           = "Введи текст поста. (можно даже приложить картинки!)"
	PostSearchingAlbumReply          = "Ищу альбомы..."
	PostFoundOneAlbumReply           = "Кажется, нашёл один вариант! Это то, что ты искал?"
	PostFoundMultipleAlbumsReply     = "Нашёл несколько вариантов. Выбери номер подходящего альбома"
	PostSelectAlbumNoSuchNumberReply = "Чтобы выбрать альбом – укажи его номер"
	PostInvalidSpotifyURLReply       = "Некорректная ссылка на альбом в спотифае!"
	PostNotFoundReply                = "К сожалению, я ничего не смог найти по этому запросу. Но ты всё ещё можешь ввести текст поста вручную!"
	PostConfirmationPromptReply      = "Точно уверен, что больше ничего не хочешь изменить в посте?"

	scenePostStateKey       = "post_state"
	scenePostFoundAlbumsKey = "post_found_albums"
	scenePostNewPostKey     = "post_user_post"

	scenePostStateWaitingForManuallyEnteredPost = "post_state_enter_manually"
	scenePostStateWaitingAlbumInfo              = "post_state_album_user_prompt"
	scenePostStateWaitingDescription            = "post_state_description_user_prompt"
	scenePostStateSelectAlbum                   = "post_state_select_from_found"
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
		menu                 = telebot.ReplyMarkup{ResizeKeyboard: true, ForceReply: true}
		cancelBtn            = menu.Text(PostReturnToMainCommand)
		enterInfoManuallyBtn = menu.Text(PostEnterContentManuallyCommand)
	)

	menu.Reply(
		menu.Row(enterInfoManuallyBtn),
		menu.Row(cancelBtn),
	)

	ssn := session.GetSession(ctx)
	ssn.Set(scenePostStateKey, scenePostStateWaitingAlbumInfo)
	ssn.Set(scenePostNewPostKey, &entity.Post{AuthorID: ctx.Message().Sender.ID})

	return ctx.Send(PostArtistReply, &menu)
}

func (pc *PostController) HandlePostUserInput(ctx telebot.Context) error {
	ssn := session.GetSession(ctx)
	ssnSceneState, ok := ssn.Get(scenePostStateKey).(string)
	if !ok {
		return errors.New("post scene state is not set")
	}

	switch ssnSceneState {
	case scenePostStateWaitingAlbumInfo:
		return pc.HandlePostUserAlbumQuery(ctx)
	case scenePostStateSelectAlbum:
		return pc.HandlePostSelectedAlbum(ctx)
	case scenePostStateWaitingDescription:
		return pc.HandleNewPostDescriptionInput(ctx)
	case scenePostStateWaitingForManuallyEnteredPost:
		return pc.HandleNewPostContentInput(ctx)
	}

	return errors.New("session is nil")
}

func (pc *PostController) HandlePostUserAlbumQuery(ctx telebot.Context) error {
	ctx.Send(PostSearchingAlbumReply)
	userQuery := ctx.Text()
	ssn := session.GetSession(ctx)

	albums, err := pc.useCase.FindAlbums(userQuery)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidSpotifyURL) {
			return ctx.Send(PostInvalidSpotifyURLReply)
		}
	}

	if len(albums) == 0 {
		return ctx.Send(PostNotFoundReply)
	}

	ssn.Set(scenePostFoundAlbumsKey, albums)
	ssn.Set(scenePostStateKey, scenePostStateSelectAlbum)
	return pc.ShowPostFoundAlbums(ctx)
}

func (pc *PostController) ShowPostFoundAlbums(ctx telebot.Context) error {
	ssn := session.GetSession(ctx)
	albums, ok := ssn.Get(scenePostFoundAlbumsKey).([]entity.Album)
	if !ok {
		return pc.ExitToMainMenu(ctx)
	}

	for i, album := range albums {
		albumMsg := "<b>" + strconv.Itoa(i+1) + ".</b> " + pc.composeAlbumInfoMsg(album)

		if album.CoverURL != "" {
			albumMsgWithCover := &telebot.Photo{
				Caption: albumMsg,
				File:    telebot.FromURL(album.CoverURL),
			}

			ctx.Send(albumMsgWithCover, telebot.ModeHTML, telebot.NoPreview)
		} else {
			ctx.Send(albumMsg, telebot.ModeHTML, telebot.NoPreview)
		}
	}

	menu := telebot.ReplyMarkup{ResizeKeyboard: true, ForceReply: true}

	if len(albums) == 1 {
		var yesBtn = menu.Text(PostFoundOneAcceptCommand)
		var noBtn = menu.Text(PostFoundOneDeclineCommand)

		menu.Reply(menu.Row(yesBtn, noBtn))

		return ctx.Send(PostFoundOneAlbumReply, &menu)
	} else {
		var menuBtns []telebot.Btn

		for i := 0; i < len(albums); i++ {
			btnText := fmt.Sprintf("№%d", i+1)
			menuBtns = append(menuBtns, menu.Text(btnText))
		}

		menu.Reply(menu.Row(menuBtns...))
		return ctx.Send(PostFoundMultipleAlbumsReply, &menu)
	}
}

func (pc *PostController) composeAlbumInfoMsg(album entity.Album) string {
	albumMsg := ""
	albumMsg += "<b>Название</b>: " + album.Title + "\n"
	albumMsg += "<b>Исполнитель</b>: " + album.Artists[0] + "\n"
	albumMsg += "<b>Страна</b>: " + album.Country + "\n"
	albumMsg += "<b>Год</b>: " + album.Year + "\n"

	if len(album.Genres) > 0 {
		albumMsg += "<b>Жанры</b>: "

		for i, genre := range album.Genres {
			genre = strings.ToLower(strings.Replace(genre, " ", "_", -1))

			if i == len(album.Genres)-1 {
				albumMsg += "#" + genre + "\n"
			} else {
				albumMsg += "#" + genre + ", "
			}
		}
	}

	if len(album.Styles) > 0 {
		albumMsg += "<b>Стили</b>: "

		for i, style := range album.Styles {
			style = strings.ToLower(strings.Replace(style, " ", "_", -1))

			if i == len(album.Styles)-1 {
				albumMsg += "#" + style + "\n\n"
			} else {
				albumMsg += "#" + style + ", "
			}
		}
	}

	if album.SpotifyLink != "" {
		albumMsg += "<a href=\"" + album.SpotifyLink + "\">Слушать в Spotify</a>\n"
	}

	return albumMsg
}

func (pc *PostController) HandlePostSelectedAlbum(ctx telebot.Context) error {
	ssn := session.GetSession(ctx)
	albums, ok := ssn.Get(scenePostFoundAlbumsKey).([]entity.Album)
	if !ok {
		return pc.ExitToMainMenu(ctx)
	}

	selectedNo, err := strconv.Atoi(ctx.Text())
	if err != nil {
		return ctx.Send(PostSelectAlbumNoSuchNumberReply)
	}

	if selectedNo > len(albums)-1 {
		return ctx.Send(PostSelectAlbumNoSuchNumberReply)
	}

	return nil
}

func (pc *PostController) HandleNewPostContentInput(ctx telebot.Context) error {
	var post entity.Post

	ssn := session.GetSession(ctx)
	err := ssn.ShouldGet(scenePostNewPostKey, &post)
	if err != nil {
		return err
	}

	post.Text = ctx.Text()
	ssn.Set(scenePostNewPostKey, post)

	var (
		menu = telebot.ReplyMarkup{ResizeKeyboard: true, ForceReply: true}
		confirmBtn = menu.Text(PostConfirmSendingCommand)
		cancelBtn  = menu.Text(PostCancelSendingCommand)
	)

	menu.Reply(menu.Row(confirmBtn), menu.Row(cancelBtn))
	return ctx.Send(PostConfirmationPromptReply, &menu)
}

func (pc *PostController) SendReport(ctx telebot.Context) error {
	var post entity.Post

	ssn := session.GetSession(ctx)
	err := ssn.ShouldGet(scenePostNewPostKey, &post)
	if err != nil {
		return err
	}

	return ctx.Send(post.Text)
}

func (pc *PostController) HandleNewPostDescriptionInput(ctx telebot.Context) error {
	return nil
}

func (pc *PostController) EnterPostContentManually(ctx telebot.Context) error {
	var (
		menu      = telebot.ReplyMarkup{ResizeKeyboard: true, ForceReply: true}
		cancelBtn = menu.Text(PostReturnToMainCommand)
	)

	menu.Reply(menu.Row(cancelBtn))

	ssn := session.GetSession(ctx)
	ssn.Set(scenePostStateKey, scenePostStateWaitingForManuallyEnteredPost)

	return ctx.Send(PostEnterManuallyReply, &menu)
}

func (pc *PostController) ExitToMainMenu(ctx telebot.Context) error {
	session.GetSession(ctx).Clear()
	return pc.stage.EnterScene(ctx, SceneMainName)
}
