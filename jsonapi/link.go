package jsonapi

import (
	"github.com/denisvmedia/urlshortener/model"
)

// Link is an object that holds link information like name, url, etc.
type Link struct {
	// Object ID - this field is ignored for the new objects, and must match the url for the existing objects.
	ID string `json:"id" example:"1"`
	// JSON:API type
	Type       string     `json:"type" example:"links"`
	Attributes model.Link `json:"attributes"`
}

// Links is an object that holds link list information
type Links struct {
	Data []Link `json:"data"`
	Meta struct {
		Links int `json:"links" example:"1" format:"int64"`
		//Took  int `json:"took" example:"1" format:"int64"`
	} `json:"meta"`
	Links struct {
		Next  string `json:"next" example:"/api/links?page[number]=1&page[size]=10"`
		Prev  string `json:"prev" example:"/api/links?page[number]=1&page[size]=10"`
		First string `json:"first" example:"/api/links?page[number]=1&page[size]=10"`
		Last  string `json:"last" example:"/api/links?page[number]=10&page[size]=10"`
	}
}

// CreateLink is an object that holds link data information
type CreateLink struct {
	Data Link `json:"data"`
}

// CreatedLink is an object that holds link data information
type CreatedLink struct {
	Data Link `json:"data"`
	Meta struct {
		//Took int `json:"took" example:"1" format:"int64"`
	} `json:"meta"`
}
