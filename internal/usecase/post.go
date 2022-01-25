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

	var album entity.Album
	if strings.Contains(query, "spotify") {
		if strings.Contains(query, "album") {
			id, err := uc.GetSpotifyAlbumID(query)
			if err != nil {
				return nil, err
			}

			album, err = uc.albumRepo.GetAlbumBySpotifyAlbumID(id)
			if err != nil {
				return nil, err
			}

			albums = append(albums, album)
		} else {
			return nil, ErrInvalidSpotifyURL
		}
	} else if strings.Contains(query, "discogs") {
		id, err := uc.GetDiscogsReleaseID(query)
		if err != nil {
			return nil, err
		}

		album, err = uc.albumRepo.GetAlbumByDiscogsReleaseID(id)
		if err != nil {
			return nil, err
		}

		albums = append(albums, album)
	} else {
		albums, err = uc.albumRepo.GetAlbumsByQuery(query, 3)
		if err != nil {
			return nil, err
		}
	}

	return albums, err
}

func (uc PostUsecase) GetSpotifyAlbumID(query string) (string, error) {
	uri, err := url.Parse(query)
	if err != nil {
		return "", ErrInvalidSpotifyURL
	}

	path := strings.Split(uri.Path, "/")
	return path[0], nil
}

func (uc PostUsecase) GetDiscogsReleaseID(query string) (string, error) {
	uri, err := url.Parse(query)
	if err != nil {
		return "", ErrInvalidDiscogsURL
	}

	path := strings.Split(uri.Path, "/")
	if len(path) < 3 {
		return "", ErrInvalidDiscogsURL
	}

	releaseNameParts := strings.Split(path[len(path)-1], "-")
	return releaseNameParts[0], nil
}