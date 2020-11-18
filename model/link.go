package model

// Link defines a link structure that is used for redirects
type Link struct {
	ID string `json:"-" swaggerignore:"true"`
	// Link short name as user requires, if empty will be generated, must be unique
	ShortName string `json:"shortName" example:"link-short-name" validate:"shortname"`
	// Original URL where to redirect the visitor
	OriginalURL string `json:"originalUrl" example:"https://example.com/my-cool-url-path" validate:"required,url,urlscheme"`
	// User comment
	Comment string `json:"comment" example:"Free text comment"`
}

// GetID to satisfy jsonapi.MarshalIdentifier interface
func (c Link) GetID() string {
	return c.ID
}

// SetID to satisfy jsonapi.UnmarshalIdentifier interface
func (c *Link) SetID(id string) error {
	c.ID = id
	return nil
}

// FillDefaults sets defaults values for those that fields are not set
// Currently sets only Link.ShortName (builds a pseudo-random string up to 8 chars len)
func (c *Link) FillDefaults() {
	if c.ShortName == "" {
		c.ShortName = generateShortName()
	}
}
