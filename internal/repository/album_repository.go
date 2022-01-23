package repository

import "github.com/Hickar/sound-seeker-bot/internal/entity"

type AlbumRepository struct {
	localSource   AlbumDatasource
	remoteSources map[string]AlbumDatasource
}

func NewAlbumRepo(local, discogs, spotify AlbumDatasource) *AlbumRepository {
	return &AlbumRepository{
		localSource:   local,
		remoteSources: map[string]AlbumDatasource{
			"discogs": discogs,
			"spotify": spotify,
		},
	}
}

func (r *AlbumRepository) GetAlbumsByQuery(query string) ([]entity.Album, error) {
	return []entity.Album{}, nil
}

func (r *AlbumRepository) GetAlbumBySpotifyAlbumID(id string) (entity.Album, error) {
	return entity.Album{}, nil
}

func (r *AlbumRepository) GetAlbumByDiscogsReleaseID(id string) (entity.Album, error) {
	return entity.Album{}, nil
}

func (r *AlbumRepository) SaveAlbum(album entity.Album) error {
	return nil
}
