package remote_datasource

import (
	"net/http"

	"github.com/Hickar/sound-seeker-bot/internal/entity"
)

type DiscogsAlbumDatasource struct {
	client *http.Client
}

func NewDiscogsDatasource(client *http.Client) *DiscogsAlbumDatasource {
	return &DiscogsAlbumDatasource{client: client}
}

func (ds *DiscogsAlbumDatasource) GetByQuery(string) ([]entity.Album, error) {
	return []entity.Album{}, nil
}

func (ds *DiscogsAlbumDatasource) GetAlbumById(string) (entity.Album, error) {
	return entity.Album{}, nil
}

func (ds *DiscogsAlbumDatasource) SaveAlbum(album entity.Album) error {
	return nil
}
