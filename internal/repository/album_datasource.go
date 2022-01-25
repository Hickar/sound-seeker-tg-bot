package repository

import "github.com/Hickar/sound-seeker-bot/internal/entity"

type AlbumDatasource interface {
	GetByQuery(string, int) ([]entity.Album, error)
	GetAlbumById(string) (entity.Album, error)
	SaveAlbum(entity.Album) error
}
