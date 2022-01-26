package remote_datasource

import "errors"

var (
	ErrSpotifyBadRefreshToken = errors.New("invalid spotify refresh token")
	ErrBadOAuth               = errors.New("invalid auth credentials")
	ErrNotFound               = errors.New("nothing was found")
	ErrExceededLimit          = errors.New("too many requests to spotify api")
	ErrInternal               = errors.New("api internal error")
)
