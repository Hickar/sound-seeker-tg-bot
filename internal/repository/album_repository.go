package repository

import (
	"errors"
	"fmt"

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
	//spotifyAlbum, err := dataSource.GetAlbumById(id)
	//if err != nil {
	//	return entity.Album{}, err
	//}



}

func (r *AlbumRepository) GetAlbumByDiscogsReleaseID(id string) (entity.Album, error) {
	album, err := r.remoteSources["discogs"].GetAlbumById(id)
	if err != nil {
		return entity.Album{}, err
	}

	query := fmt.Sprintf("%s %s", album.Artists[0], album.Title)
	spotifyResults, err := r.remoteSources["spotify"].GetByQuery(query)
	if err != nil || len(spotifyResults) == 0 {
		if err := r.localSource.SaveAlbum(album); err != nil {
			return entity.Album{}, err
		}

		return album, nil
	}

	album = r.combineSpotifyAndDiscogsAlbumInfo(album, spotifyResults[0])
	if err := r.localSource.SaveAlbum(album); err != nil {
		return entity.Album{}, err
	}

	return album, nil
}

func (r *AlbumRepository) SaveAlbum(album entity.Album) error {
	return nil
}

func (r *AlbumRepository) combineSpotifyAndDiscogsAlbumInfo(spotify, discogs entity.Album) entity.Album {
	return entity.Album{
		Artists:   spotify.Artists,
		Title:     spotify.Title,
		Country:   "",
		Year:      discogs.Year,
		Genres:    discogs.Genres,
		Styles:    discogs.Genres,
		Link:      spotify.Link,
		SpotifyId: spotify.SpotifyId,
		DiscogsId: discogs.DiscogsId,
	}
}