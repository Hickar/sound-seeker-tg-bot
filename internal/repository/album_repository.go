package repository

import (
	"errors"

	"github.com/Hickar/sound-seeker-bot/internal/entity"
)

type AlbumRepository struct {
	localSource   AlbumDatasource
	remoteSources map[string]AlbumDatasource
}

func NewAlbumRepo(local, discogs, spotify AlbumDatasource) *AlbumRepository {
	return &AlbumRepository{
		localSource: local,
		remoteSources: map[string]AlbumDatasource{
			"discogs": discogs,
			"spotify": spotify,
		},
	}
}

func (r *AlbumRepository) GetAlbumsByQuery(query string) ([]entity.Album, error) {
	dataSource, ok := r.remoteSources["spotify"]
	if !ok {
		return []entity.Album{}, errors.New("spotify data source is not defined")
	}

	return dataSource.GetByQuery(query)
}

func (r *AlbumRepository) GetAlbumBySpotifyAlbumID(id string) (entity.Album, error) {
	dataSource, ok := r.remoteSources["spotify"]
	if !ok {
		return entity.Album{}, errors.New("spotify data source is not defined")
	}

	return dataSource.GetAlbumById(id)
}

func (r *AlbumRepository) GetAlbumByDiscogsReleaseID(id string) (*entity.Album, error) {
	return nil, nil
}

func (r *AlbumRepository) SaveAlbum(album entity.Album) error {
	return nil
}
