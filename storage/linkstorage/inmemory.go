package linkstorage

import (
	"fmt"
	"github.com/denisvmedia/urlshortener/storage"
	"sort"
	"sync"
	"sync/atomic"

	"github.com/denisvmedia/urlshortener/model"
	"github.com/go-extras/errors"
)

// sorting
type byID []*model.Link

func (c byID) Len() int {
	return len(c)
}

func (c byID) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c byID) Less(i, j int) bool {
	return c[i].GetID() < c[j].GetID()
}

// NewInMemoryStorage initializes the storage
func NewInMemoryStorage() Storage {
	return &InMemoryStorage{
		links:            make(map[string]*model.Link),
		linksByShortName: make(map[string]*model.Link),
		linksByID:        make([]*model.Link, 0),
		idCount:          0,
	}
}

// InMemoryStorage stores all of the links, needs to be injected into
// User and Link Resource. In the real world, you would use a database for that.
type InMemoryStorage struct {
	links            map[string]*model.Link
	linksByShortName map[string]*model.Link
	linksByID        []*model.Link
	idCount          int64
	lock             sync.RWMutex
}

// PaginatedGetAll returns a slice of links according to desired pagination and total number of items
func (s *InMemoryStorage) PaginatedGetAll(pageNumber, pageSize int) (results []*model.Link, total int, err error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	start, end := storage.SlicePaginate(pageNumber-1, pageSize, len(s.linksByID))
	results = s.linksByID[start:end]

	return results, len(s.linksByID), nil
}

// GetOne link
func (s *InMemoryStorage) GetOne(id string) (*model.Link, error) {
	s.lock.RLock()
	link, ok := s.links[id]
	s.lock.RUnlock()
	if ok {
		return link, nil
	}

	return nil, errors.Wrapf(storage.ErrNotFound, "Link for id %s not found", id)
}

// GetOneByShortName returns a link byt its short name
func (s *InMemoryStorage) GetOneByShortName(shortName string) (*model.Link, error) {
	s.lock.RLock()
	link, ok := s.linksByShortName[shortName]
	s.lock.RUnlock()
	if !ok {
		return nil, errors.Wrapf(storage.ErrNotFound, "Link for shortName %s not found", shortName)
	}

	return link, nil
}

// Insert a fresh one
func (s *InMemoryStorage) Insert(c model.Link) (*model.Link, error) {
	atomic.AddInt64(&s.idCount, 1)
	id := fmt.Sprintf("%d", atomic.LoadInt64(&s.idCount))
	c.ID = id

	s.lock.Lock()
	defer s.lock.Unlock()
	if lv, exists := s.linksByShortName[c.ShortName]; exists {
		return lv, errors.Wrapf(storage.ErrShortNameAlreadyExists, "Existing link id %s", lv.ID)
	}

	s.linksByShortName[c.ShortName] = &c
	s.links[id] = &c
	s.linksByID = append(s.linksByID, &c)

	//// the following code is commented out assuming that we always get the most recent id (although there's a slight chance to have it inaccurate)
	//s.linksByID = make([]*model.Link, 0, len(s.links))
	//for key := range s.links {
	//	s.linksByID = append(s.linksByID, s.links[key])
	//}
	//sort.Sort(byID(s.linksByID))

	return &c, nil
}

// Delete one :(
func (s *InMemoryStorage) Delete(id string) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	link, exists := s.links[id]
	if !exists {
		return errors.Wrapf(storage.ErrNotFound, "Link for id %s not found", id)
	}
	delete(s.links, id)
	delete(s.linksByShortName, link.ShortName)

	// The following is kinda heavy operation, but unavoidable (well, a possible option
	// would be storing the order index as well, and then deleting this item only by
	// a slice trick, but we don't store the item id in the slice).
	s.linksByID = make([]*model.Link, 0, len(s.links))
	for key := range s.links {
		s.linksByID = append(s.linksByID, s.links[key])
	}
	sort.Sort(byID(s.linksByID))

	return nil
}

// Update updates an existing link
func (s *InMemoryStorage) Update(c model.Link) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	_, exists := s.links[c.ID]
	if !exists {
		return errors.Wrapf(storage.ErrNotFound, "Link for id %s not found", c.ID)
	}
	if existing, exists := s.linksByShortName[c.ShortName]; exists && existing.ShortName == c.ShortName {
		return errors.Wrapf(storage.ErrShortNameAlreadyExists, "Existing link id %s", existing.ID)
	}
	s.links[c.ID] = &c

	return nil
}
