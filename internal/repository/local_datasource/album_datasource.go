package local_datasource

import (
	"github.com/Hickar/sound-seeker-bot/internal/entity"
	"gorm.io/gorm"
)

type AlbumDatasource struct {
	db *gorm.DB
}

func New(db *gorm.DB) *AlbumDatasource {
	return &AlbumDatasource{db: db}
}

func (ds *AlbumDatasource) GetByQuery(string) ([]entity.Album, error) {
	return []entity.Album{}, nil
}

func (ds *AlbumDatasource) GetAlbumById(string) (entity.Album, error) {
	return entity.Album{}, nil
}

func (ds *AlbumDatasource) SaveAlbum(album entity.Album) error {
	return nil
}
