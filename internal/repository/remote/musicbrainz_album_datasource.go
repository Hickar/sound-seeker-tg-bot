package remote

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Hickar/sound-seeker-bot/internal/entity"
)

const (
	_musicBrainzSearchAlbumsEndpoint = "https://musicbrainz.org/ws/2/release?query=%s&limit=%d"
	_musicBrainzGetArtistByIdEndpoint = "https://musicbrainz.org/ws/2/artist/%s"
)

type musicBrainzGetArtistResponse struct {
	Country string `json:"country"`
}

type musicBrainzSearchResponse struct {
	Releases []musicBrainzAlbumDto `json:"releases"`
}

type musicBrainzAlbumDto struct {
	Title         string `json:"title"`
	ArtistCredits []struct {
		Artist struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"artist"`
	} `json:"artist-credit"`
	Date string `json:"date"`
}

type MusicBrainzAlbumDatasource struct {
	client *http.Client
}

func NewMusicBrainzDatasource(client *http.Client) *MusicBrainzAlbumDatasource {
	return &MusicBrainzAlbumDatasource{client: client}
}

func (ds *MusicBrainzAlbumDatasource) GetByQuery(query string, limit int) ([]entity.Album, error) {
	var albums []entity.Album

	query = strings.Replace(query, " ", "+", -1)

	endpoint := fmt.Sprintf(_musicBrainzSearchAlbumsEndpoint, strings.ToLower(query), limit)
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return albums, err
	}

	req.Header.Set("Accept", "application/json")

	resp, err := ds.client.Do(req)
	if err != nil {
		return albums, err
	}

	if resp.StatusCode != http.StatusOK {
		return albums, ErrInternal
	}

	defer resp.Body.Close()
	respBytes, _ := ioutil.ReadAll(resp.Body)
	var respResults musicBrainzSearchResponse

	if err := json.Unmarshal(respBytes, &respResults); err != nil {
		return albums, err
	}

	for _, result := range respResults.Releases {
		country, err := ds.getArtistCountry(result.ArtistCredits[0].Artist.ID)
		if err != nil {
			continue
		}

		album := ds.musicBrainzAlbumDtoToEntity(result, country)
		albums = append(albums, album)
	}

	return albums, nil
}

func (ds *MusicBrainzAlbumDatasource) GetAlbumById(id string) (entity.Album, error) {
	return entity.Album{}, nil
}

func (ds *MusicBrainzAlbumDatasource) musicBrainzAlbumDtoToEntity(dto musicBrainzAlbumDto, country string) entity.Album {
	var album entity.Album

	for _, artist := range dto.ArtistCredits {
		album.Artists = append(album.Artists, artist.Artist.Name)
	}

	album.Title = dto.Title
	album.Country = country

	releaseDate, _ := time.Parse("2006-01-02", dto.Date)
	album.Year = strconv.Itoa(releaseDate.Year())

	return album
}

func (ds *MusicBrainzAlbumDatasource) getArtistCountry(artistId string) (string, error) {
	var country string

	req, err := http.NewRequest("GET", fmt.Sprintf(_musicBrainzGetArtistByIdEndpoint, artistId), nil)
	if err != nil {
		return country, err
	}

	req.Header.Set("Accept", "application/json")

	resp, err := ds.client.Do(req)
	if err != nil {
		return country, err
	}

	if resp.StatusCode != http.StatusOK {
		return country, ErrInternal
	}

	defer resp.Body.Close()
	respBytes, _ := ioutil.ReadAll(resp.Body)
	var respArtist musicBrainzGetArtistResponse

	if err := json.Unmarshal(respBytes, &respArtist); err != nil {
		return country, err
	}

	return respArtist.Country, nil
}

func (ds *MusicBrainzAlbumDatasource) SaveAlbum(album entity.Album) error {
	return nil
}
