package usecase

import (
	"net/url"
	"strings"

	"github.com/Hickar/sound-seeker-bot/internal/entity"
	"github.com/Hickar/sound-seeker-bot/internal/repository"
)

type PostUsecase struct {
	albumRepo repository.AlbumRepository
}

func NewPostUsecase(repo *repository.AlbumRepository) *PostUsecase {
	return &PostUsecase{albumRepo: *repo}
}

func (uc *PostUsecase) FindAlbums(query string) ([]entity.Album, error) {
	var (
		albums []entity.Album
		err    error
	)

	if strings.Contains(query, "spotify") {
		var album entity.Album
		if strings.Contains(query, "album") {
			uri, err := url.Parse(query)
			if err != nil {
				return nil, ErrInvalidSpotifyURL
			}

			path := strings.Split(uri.Path, "/")
			id := path[len(path)-1]

			album, err = uc.albumRepo.GetAlbumBySpotifyAlbumID(id)
			if err != nil {
				return nil, err
			}

			albums = append(albums, album)
		} else {
			return nil, ErrInvalidSpotifyURL
		}
	} else if strings.Contains(query, "discogs") {
		//album, err = uc.albumRepo.GetAlbumByDiscogsReleaseID()
	} else {
		albums, err = uc.albumRepo.GetAlbumsByQuery(query)
		if err != nil {
			return nil, err
		}
	}

	return albums, err
}
