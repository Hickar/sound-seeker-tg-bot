package entity

type Post struct {
	Text     string `json:"text"`
	AuthorID int64  `json:"author_id"`
	AlbumID  int64  `json:"album_id"`
}