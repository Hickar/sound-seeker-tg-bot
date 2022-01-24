package remote_datasource

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/Hickar/sound-seeker-bot/internal/config"
	"github.com/Hickar/sound-seeker-bot/internal/entity"
)

const (
	_spotifyGetAlbumByIdEndpoint   = "https://api.spotify.com/v1/albums/%s"
	_spotifyGetAccessTokenEndpoint = "https://accounts.spotify.com/api/token"
)

type spotifyAuthResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type spotifyAlbumDto struct {
	Artists []struct {
		Name string `json:"name"`
	} `json:"artists"`
	Title       string   `json:"name,omitempty"`
	Genres      []string `json:"genres"`
	ReleaseDate string   `json:"release_date"`
}

type SpotifyAlbumDatasource struct {
	client         *http.Client
	credentials    config.SpotifyConfig
	tokenExpiresIn time.Time
	accessToken    string
}

func NewSpotifyDatasource(client *http.Client, credentials config.SpotifyConfig) *SpotifyAlbumDatasource {
	return &SpotifyAlbumDatasource{client: client, credentials: credentials}
}

func (ds *SpotifyAlbumDatasource) GetByQuery(query string) ([]entity.Album, error) {
	return []entity.Album{}, nil
}

func (ds *SpotifyAlbumDatasource) GetAlbumById(id string) (entity.Album, error) {
	var album entity.Album

	accessToken, err := ds.getSpotifyAccessToken(ds.credentials.ClientId, ds.credentials.ClientSecret)
	if err != nil {
		return album, err
	}

	req, _ := http.NewRequest("GET", fmt.Sprintf(_spotifyGetAlbumByIdEndpoint, id), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	resp, err := ds.client.Do(req)
	if err != nil {
		return album, err
	}

	if resp.StatusCode != http.StatusOK {
		return album, fmt.Errorf("unable to get spotify album by id: got code %d instead of %d", resp.StatusCode, http.StatusOK)
	}

	defer resp.Body.Close()
	respBytes, _ := ioutil.ReadAll(resp.Body)
	var respAlbum spotifyAlbumDto

	if err := json.Unmarshal(respBytes, &respAlbum); err != nil {
		return album, err
	}

	album = spotifyAlbumDtoToEntity(respAlbum)
	return album, nil
}

func (ds *SpotifyAlbumDatasource) getSpotifyAccessToken(id, secret string) (string, error) {
	if time.Now().Before(ds.tokenExpiresIn) {
		return ds.accessToken, nil
	}

	refreshToken := fmt.Sprintf("%s:%s", id, secret)
	accessToken := base64.StdEncoding.EncodeToString([]byte(refreshToken))

	formData := url.Values{}
	formData.Set("grant_type", "client_credentials")
	encodedData := formData.Encode()

	req, _ := http.NewRequest("POST", _spotifyGetAccessTokenEndpoint, strings.NewReader(encodedData))
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", accessToken))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(encodedData)))

	resp, err := ds.client.Do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unable to get spotify album by id: got code %d instead of %d", resp.StatusCode, http.StatusOK)
	}

	defer resp.Body.Close()
	respByte, _ := ioutil.ReadAll(resp.Body)
	var authInfo spotifyAuthResponse
	if err := json.Unmarshal(respByte, &authInfo); err != nil {
		return "", nil
	}

	ds.tokenExpiresIn = time.Now().Add(time.Second * time.Duration(authInfo.ExpiresIn))
	return authInfo.AccessToken, nil
}

func (ds *SpotifyAlbumDatasource) SaveAlbum(album entity.Album) error {
	return nil
}

func spotifyAlbumDtoToEntity(dto spotifyAlbumDto) entity.Album {
	var album entity.Album

	album.Title = dto.Title
	for _, genre := range dto.Genres {
		album.Genres = append(album.Genres, genre)
	}
	for _, artist := range dto.Artists {
		album.Artists = append(album.Artists, artist.Name)
	}

	releaseDate, _ := time.Parse("2006-01-02", dto.ReleaseDate)
	album.Year = strconv.Itoa(releaseDate.Year())

	return album
}
