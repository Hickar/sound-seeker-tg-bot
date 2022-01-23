package remote_datasource

import (
	"net/http"

	"github.com/Hickar/sound-seeker-bot/internal/config"
	"github.com/Hickar/sound-seeker-bot/internal/entity"
)

type SpotifyAlbumDatasource struct {
	client      *http.Client
	credentials config.SpotifyConfig
}

func NewSpotifyDatasource(client *http.Client, credentials config.SpotifyConfig) *SpotifyAlbumDatasource {
	return &SpotifyAlbumDatasource{client: client, credentials: credentials}
}

func (ds *SpotifyAlbumDatasource) GetByQuery(string) ([]entity.Album, error) {
	return []entity.Album{}, nil
}

func (ds *SpotifyAlbumDatasource) GetById(string) (entity.Album, error) {
	return entity.Album{}, nil
}

func (ds *SpotifyAlbumDatasource) SaveAlbum(album entity.Album) error {
	return nil
}
