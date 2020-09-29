package resource

import (
	"github.com/denisvmedia/urlshortener/storage"
	"github.com/denisvmedia/urlshortener/storage/linkstorage"
	myvalidator "github.com/denisvmedia/urlshortener/validator"
	"github.com/go-extras/errors"
	"github.com/go-playground/validator/v10"
	"net/http"

	"github.com/denisvmedia/urlshortener/model"
	"github.com/go-extras/api2go"
)

// LinkResource for api2go routes
type LinkResource struct {
	LinkStorage linkstorage.Storage
	validator   *validator.Validate
}

func NewLinkResource(linkStorage linkstorage.Storage) *LinkResource {
	// Validator is not injected as a dependency, because it's actually an integral part of LinkResource
	validate := validator.New()
	err := validate.RegisterValidation("shortname", myvalidator.ValidateUrlShortName)
	if err != nil {
		panic(err) // this should never happen
	}
	err = validate.RegisterValidation("urlscheme", myvalidator.ValidateUrlScheme)
	if err != nil {
		panic(err) // this should never happen
	}

	return &LinkResource{
		LinkStorage: linkStorage,
		validator:   validate,
	}
}

// FindAll links
// @Summary List links
// @Description get links
// @Tags links
// @Accept  json-api
// @Produce  json-api
// @Param page[number] query int false "Page number" default(1)
// @Param page[size] query int false "Page size" default(10) maximum(1000)
// @Success 200 {object} jsonapi.Links
// @Router /links [get]
func (c *LinkResource) FindAll(r api2go.Request) (api2go.Responder, error) {
	pagination := parsePageArgs(r.QueryParams)

	links, total, err := c.LinkStorage.PaginatedGetAll(pagination.Number, pagination.Size)
	if err != nil {
		return nil, HttpErrorPtrWithStatus(err, internalServerError)
	}

	result := &api2go.Response{
		Res:  links,
		Code: http.StatusOK,
		Meta: map[string]interface{}{
			"links": total,
		},
		Pagination: getPagination(pagination.Number, pagination.Size, total),
	}

	return result, nil
}

// FindOne link
// @Summary Get a link
// @Description get link by ID
// @Tags links
// @Accept  json-api
// @Produce  json-api
// @Param id path string true "Link ID"
// @Success 200 {object} jsonapi.Link
// @Router /links/{id} [get]
func (c *LinkResource) FindOne(ID string, r api2go.Request) (api2go.Responder, error) {
	res, err := c.LinkStorage.GetOne(ID)
	if err != nil {
		if err == storage.ErrNotFound {
			return nil, HttpErrorPtrWithStatus(err, resourceNotFound)
		}
		return nil, HttpErrorPtrWithStatus(err, internalServerError)
	}
	return &Response{Res: res}, nil
}

// Create a new link
// @Summary Create a new link
// @Description add by link json
// @Tags links
// @Accept  json-api
// @Produce  json-api
// @Param link body jsonapi.CreateLink true "Add link"
// @Success 201 {object} jsonapi.CreatedLink
// @Router /links [post]
func (c *LinkResource) Create(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	link, ok := obj.(model.Link)
	if !ok {
		return nil, HttpErrorPtrWithStatus(errors.New("Invalid instance given"), "")
	}

	if err := c.validator.Struct(link); err != nil {
		return nil, HttpErrorPtrWithStatus(err, validationError)
	}

	newLink, err := c.LinkStorage.Insert(link)
	if err != nil {
		return nil, HttpErrorPtrWithStatus(err, errors.Cause(err).Error())
	}
	return &Response{Res: newLink, Code: http.StatusCreated}, nil
}

// Delete a link :(
// @Summary Delete a link
// @Description Delete by link ID
// @Tags links
// @Accept  json-api
// @Produce  json-api
// @Param  id path int true "Link ID"
// @Success 204
// @Router /links/{id} [delete]
func (c *LinkResource) Delete(id string, r api2go.Request) (api2go.Responder, error) {
	err := c.LinkStorage.Delete(id)
	if err != nil {
		return nil, HttpErrorPtrWithStatus(err, resourceNotFound)
	}
	return &Response{Code: http.StatusNoContent}, nil
}

// Update a link
// @Summary Update a link
// @Description Update by link json
// @Tags links
// @Accept  json-api
// @Produce  json-api
// @Param  id path int true "Link ID"
// @Param  account body jsonapi.CreateLink true "Update link"
// @Success 200 {object} jsonapi.CreatedLink
// @Router /links/{id} [patch]
func (c *LinkResource) Update(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	link, ok := obj.(model.Link)
	if !ok {
		var linkPtr *model.Link
		linkPtr, ok = obj.(*model.Link)
		if !ok {
			return nil, HttpErrorPtrWithStatus(errors.New("Invalid instance given"), "")
		}
		link = *linkPtr
	}

	if err := c.validator.Struct(link); err != nil {
		return nil, HttpErrorPtrWithStatus(err, validationError)
	}

	err := c.LinkStorage.Update(link)
	if err != nil {
		return nil, HttpErrorPtrWithStatus(err, resourceNotFound)
	}

	return &Response{Res: link, Code: http.StatusOK}, nil
}
