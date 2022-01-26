package remote_datasource

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Hickar/sound-seeker-bot/internal/entity"
)

const (
	_discogsGetAlbumByIdEndpoint = "https://api.discogs.com/masters/%s"
	_discogsSearchAlbumsEndpoint = "https://api.discogs.com/database/search?q=%s&type=release"
)

type discogsSearchResponse struct {
	Results []struct {
		ID uint64 `json:"master_id"`
	} `json:"results"`
}

type discogsAlbumDto struct {
	ID      uint64   `json:"id"`
	Genres  []string `json:"genres,genre"`
	Styles  []string `json:"styles,style"`
	Year    int      `json:"year"`
	Artists []struct {
		Name string `json:"name"`
	} `json:"artists"`
	Title string `json:"title"`
}

type DiscogsCredentials struct {
	ConsumerKey   string
	ConsumerToken string
	OAuthToken    string
	OAuthSecret   string
	VerifyKey     string
}

type DiscogsAlbumDatasource struct {
	client      *http.Client
	credentials DiscogsCredentials
}

func NewDiscogsDatasource(client *http.Client, credentials DiscogsCredentials) *DiscogsAlbumDatasource {
	return &DiscogsAlbumDatasource{client: client, credentials: credentials}
}

func (ds *DiscogsAlbumDatasource) GetByQuery(query string, limit int) ([]entity.Album, error) {
	var albums []entity.Album

	query = strings.Replace(query, " ", "+", -1)
	req, err := http.NewRequest("GET", fmt.Sprintf(_discogsSearchAlbumsEndpoint, query), nil)
	if err != nil {
		return albums, err
	}

	authHeader := ds.buildAuthHeader(ds.credentials)
	req.Header.Set("Authorization", authHeader)

	resp, err := ds.client.Do(req)
	if err != nil {
		return albums, err
	}

	if resp.StatusCode != http.StatusOK {
		return albums, ErrInternal
	}

	defer resp.Body.Close()
	respBytes, _ := ioutil.ReadAll(resp.Body)
	var respResults discogsSearchResponse

	if err := json.Unmarshal(respBytes, &respResults); err != nil {
		return albums, err
	}

	for i, result := range respResults.Results {
		if i == limit {
			break
		}

		id := strconv.FormatUint(result.ID, 10)
		album, err := ds.GetAlbumById(id)
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				continue
			} else {
				return albums, err
			}
		}

		albums = append(albums, album)
	}

	return albums, nil
}

func (ds *DiscogsAlbumDatasource) GetAlbumById(id string) (entity.Album, error) {
	var album entity.Album

	req, _ := http.NewRequest("GET", fmt.Sprintf(_discogsGetAlbumByIdEndpoint, id), nil)

	authHeader := ds.buildAuthHeader(ds.credentials)
	req.Header.Set("Authorization", authHeader)

	resp, err := ds.client.Do(req)
	if err != nil {
		return album, err
	}

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusForbidden:
			err = ErrBadOAuth
		case http.StatusNotFound:
			err = ErrNotFound
		case http.StatusTooManyRequests:
			err = ErrExceededLimit
		case http.StatusInternalServerError:
			err = ErrInternal
		}
		return album, err
	}

	defer resp.Body.Close()
	respBytes, _ := ioutil.ReadAll(resp.Body)
	var respAlbum discogsAlbumDto

	if err := json.Unmarshal(respBytes, &respAlbum); err != nil {
		return album, err
	}

	album = ds.discogsAlbumDtoToEntity(respAlbum)
	return album, nil
}

func (ds *DiscogsAlbumDatasource) discogsAlbumDtoToEntity(dto discogsAlbumDto) entity.Album {
	var album entity.Album

	album.DiscogsId = strconv.FormatUint(dto.ID, 10)
	album.Title = dto.Title
	album.Year = strconv.Itoa(dto.Year)
	album.Genres = dto.Genres
	album.Styles = dto.Styles

	for _, artist := range dto.Artists {
		album.Artists = append(album.Artists, artist.Name)
	}

	return album
}

func (ds *DiscogsAlbumDatasource) buildAuthHeader(credentials DiscogsCredentials) string {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	return fmt.Sprintf(`OAuth oauth_consumer_key="%s",oauth_token="%s",oauth_signature_method="PLAINTEXT",oauth_timestamp="%s",oauth_nonce="cSrDRTKTPJw",oauth_version="1.0",oauth_verifier="%s",oauth_signature="%s"`,
		credentials.ConsumerKey,
		credentials.OAuthToken,
		timestamp,
		credentials.VerifyKey,
		fmt.Sprintf("%s&%s", credentials.ConsumerToken, credentials.OAuthSecret))
}

func (ds *DiscogsAlbumDatasource) SaveAlbum(album entity.Album) error {
	return nil
}
