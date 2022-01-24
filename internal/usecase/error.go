package usecase

import "errors"

var (
	ErrInvalidSpotifyURL = errors.New("invalid spotify URI was provided")
	ErrInvalidDiscogsURL = errors.New("invalid discogs URI was provided")
)
