package repository

import "github.com/Hickar/sound-seeker-bot/internal/entity"

type AlbumDatasource interface {
	GetByQuery(string) ([]entity.Album, error)
	GetAlbumById(string) (entity.Album, error)
	SaveAlbum(album entity.Album) error
}
