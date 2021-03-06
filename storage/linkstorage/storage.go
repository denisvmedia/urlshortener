package linkstorage

import (
	"github.com/denisvmedia/urlshortener/model"
)

// Storage defines an interface that must be implemented in order to be used as a backend to store the links
type Storage interface {
	PaginatedGetAll(pageNumber, pageSize int) (results []*model.Link, total int, err error)
	GetOne(id string) (*model.Link, error)
	GetOneByShortName(shortName string) (*model.Link, error)
	Insert(c model.Link) (*model.Link, error)
	Delete(id string) error
	Update(c model.Link) error
}
