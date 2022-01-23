package local_datasource

import "github.com/Hickar/sound-seeker-bot/internal/entity"

type AlbumDatasource struct {
	//db
}

func (ds *AlbumDatasource) GetByQuery(string) ([]entity.Album, error) {
	return []entity.Album{}, nil
}
func (ds *AlbumDatasource) GetById(string) (entity.Album, error) {
	return entity.Album{}, nil
}

func (ds *AlbumDatasource) SaveAlbum(album entity.Album) error {
	return nil
}
