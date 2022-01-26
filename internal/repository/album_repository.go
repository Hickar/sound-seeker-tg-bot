package repository

import (
	"errors"
	"fmt"
	"math"

	"github.com/Hickar/sound-seeker-bot/internal/entity"
)

type AlbumRepository struct {
	localSource   AlbumDatasource
	remoteSources map[string]AlbumDatasource
}

func NewAlbumRepo(local, discogs, spotify, musicBrainz AlbumDatasource) *AlbumRepository {
	return &AlbumRepository{
		localSource: local,
		remoteSources: map[string]AlbumDatasource{
			"discogs":     discogs,
			"spotify":     spotify,
			"musicbrainz": musicBrainz,
		},
	}
}

func (r *AlbumRepository) GetAlbumsByQuery(query string, limit int) ([]entity.Album, error) {
	var albums []entity.Album

	mbAlbumResults, err := r.remoteSources["musicbrainz"].GetByQuery(query, limit)
	if err != nil {
		return albums, err
	}

	dsgsAlbumResults, err := r.remoteSources["discogs"].GetByQuery(query, limit)
	if err != nil {
		return albums, err
	}

	stfyAlbumResults, err := r.remoteSources["spotify"].GetByQuery(query, limit)
	if err != nil {
		return albums, err
	}

	if areEmpty(mbAlbumResults, dsgsAlbumResults, stfyAlbumResults)  {
		return albums, nil
	}

	return r.composeAlbumSlicesInfo(mbAlbumResults, dsgsAlbumResults, stfyAlbumResults), nil
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
	spotifyResults, err := r.remoteSources["spotify"].GetByQuery(query, 1)
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

func (r *AlbumRepository) composeAlbumSlicesInfo(albumsResults ...[]entity.Album) []entity.Album {
	var albums []entity.Album

	if areEmpty(albumsResults...) {
		return albums
	}

	totalMinLen := totalSlicesMinLen(albumsResults...)
	for i := 0; i < totalMinLen; i++ {
		var composables []entity.Album
		for _, result := range albumsResults {
			if len(result) > 0 {
				composables = append(composables, result[i])
			}
		}

		composedAlbum := r.composeAlbumsInfo(composables...)
		albums = append(albums, composedAlbum)
	}

	return albums
}

func (r *AlbumRepository) composeAlbumsInfo (albums ...entity.Album) entity.Album {
	var resultAlbum entity.Album

	for _, album := range albums {
		if len(resultAlbum.Artists) == 0 {
			resultAlbum.Artists = album.Artists
		}

		if resultAlbum.Title == "" {
			resultAlbum.Title = album.Title
		}

		if resultAlbum.Country == "" {
			resultAlbum.Country = album.Country
		}

		if resultAlbum.Year == "" {
			resultAlbum.Year = album.Year
		}

		if len(resultAlbum.Genres) == 0 {
			resultAlbum.Genres = album.Genres
		}

		if len(resultAlbum.Styles) == 0 {
			resultAlbum.Styles = album.Styles
		}

		if resultAlbum.SpotifyLink == "" {
			resultAlbum.SpotifyLink = album.SpotifyLink
		}

		if resultAlbum.SpotifyId == "" {
			resultAlbum.SpotifyId = album.SpotifyId
		}

		if resultAlbum.DiscogsId == "" {
			resultAlbum.DiscogsId = album.DiscogsId
		}
	}

	return resultAlbum
}

func (r *AlbumRepository) combineSpotifyAndDiscogsAlbumInfo(spotify, discogs entity.Album) entity.Album {
	return entity.Album{
		Artists:     spotify.Artists,
		Title:       spotify.Title,
		Country:     "",
		Year:        discogs.Year,
		Genres:      discogs.Genres,
		Styles:      discogs.Genres,
		SpotifyLink: spotify.SpotifyLink,
		SpotifyId:   spotify.SpotifyId,
		DiscogsId:   discogs.DiscogsId,
	}
}

func totalSlicesMinLen(slices ...[]entity.Album) int {
	min := math.MaxInt

	for _, slice := range slices {
		if len(slice) < min {
			min = len(slice)
		}
	}

	return min
}

func areEmpty(slices ...[]entity.Album) bool {
	for _, slice := range slices {
		if len(slice) != 0 {
			return false
		}
	}

	return true
}