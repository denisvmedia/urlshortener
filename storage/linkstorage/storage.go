package linkstorage

import (
	"encoding/base64"
	"github.com/denisvmedia/urlshortener/model"
	"github.com/google/uuid"
	"strings"
)

type Storage interface {
	PaginatedGetAll(pageNumber, pageSize int) (results []*model.Link, total int, err error)
	GetOne(id string) (*model.Link, error)
	GetOneByShortName(shortName string) (*model.Link, error)
	Insert(c model.Link) (*model.Link, error)
	Delete(id string) error
	Update(c model.Link) error
}

func generateShortName() string {
	hash, _ := uuid.New().MarshalBinary()
	var b strings.Builder
	encoder := base64.NewEncoder(base64.URLEncoding, &b)
	_, _ = encoder.Write(hash)
	_ = encoder.Close()

	res := strings.ReplaceAll(strings.Trim(b.String(), "_-=\n"), "_", "")
	if len(res) > 8 {
		res = res[:8]
	}

	return res
}
