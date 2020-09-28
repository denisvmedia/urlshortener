/*
Package examples shows how to implement a basic CRUD for two data structures with the api2go server functionality.
To play with this example server you can run some of the following curl requests

In order to demonstrate dynamic baseurl handling for requests, apply the --header="REQUEST_URI:https://www.your.domain.example.com" parameter to any of the commands.

Create a new user:
	curl -X POST http://localhost:31415/v0/users -d '{"data" : {"type" : "users" , "attributes": {"user-name" : "marvin"}}}'

List users:
	curl -X GET http://localhost:31415/v0/users

List paginated users:
	curl -X GET 'http://localhost:31415/v0/users?page\[offset\]=0&page\[limit\]=2'
OR
	curl -X GET 'http://localhost:31415/v0/users?page\[number\]=1&page\[size\]=2'

Update:
	curl -vX PATCH http://localhost:31415/v0/users/1 -d '{ "data" : {"type" : "users", "id": "1", "attributes": {"user-name" : "better marvin"}}}'

Delete:
	curl -vX DELETE http://localhost:31415/v0/users/2

Create a chocolate with the name sweet
	curl -X POST http://localhost:31415/v0/chocolates -d '{"data" : {"type" : "chocolates" , "attributes": {"name" : "Ritter Sport", "taste": "Very Good"}}}'

Create a user with a sweet
	curl -X POST http://localhost:31415/v0/users -d '{"data" : {"type" : "users" , "attributes": {"user-name" : "marvin"}, "relationships": {"sweets": {"data": [{"type": "chocolates", "id": "1"}]}}}}'

List a users sweets
	curl -X GET http://localhost:31415/v0/users/1/sweets

Replace a users sweets
	curl -X PATCH http://localhost:31415/v0/users/1/relationships/sweets -d '{"data" : [{"type": "chocolates", "id": "2"}]}'

Add a sweet
	curl -X POST http://localhost:31415/v0/users/1/relationships/sweets -d '{"data" : [{"type": "chocolates", "id": "2"}]}'

Remove a sweet
	curl -X DELETE http://localhost:31415/v0/users/1/relationships/sweets -d '{"data" : [{"type": "chocolates", "id": "2"}]}'
*/
package main

import (
	"fmt"
	"github.com/denisvmedia/urlshortener/model"
	"github.com/denisvmedia/urlshortener/resource"
	"github.com/denisvmedia/urlshortener/routing"
	"github.com/denisvmedia/urlshortener/shortener"
	"github.com/denisvmedia/urlshortener/storage"
	myvalidator "github.com/denisvmedia/urlshortener/validator"
	"github.com/go-extras/api2go"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "github.com/denisvmedia/urlshortener/docs"

	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title URL Shortener Example
// @version 1.0
// @description This is an example url-shortener server.

// @contact.name API Support
// @contact.url https://github.com/denisvmedia/urlshortener/issues
// @contact.email ask@artprima.cz

// @license.name MIT
// @license.url https://github.com/denisvmedia/urlshortener/blob/master/LICENSE

// @BasePath /api

func main() {
	port := 31415

	e := echo.New()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	api := api2go.NewAPIWithRouting(
		"api",
		api2go.NewStaticResolver("/"),
		routing.Echo(e),
	)

	validate := validator.New()
	err := validate.RegisterValidation("shortname", myvalidator.ValidateUrlShortName)
	if err != nil {
		panic(err) // this should never happen
	}
	err = validate.RegisterValidation("urlscheme", myvalidator.ValidateUrlScheme)
	if err != nil {
		panic(err) // this should never happen
	}
	linkStorage := storage.NewLinkStorage()
	api.AddResource(model.Link{}, resource.LinkResource{
		LinkStorage: linkStorage,
		Validator:   validate,
	})

	fmt.Printf("Listening on :%d\n", port)
	e.GET("/swagger/*any", echoSwagger.EchoWrapHandler())
	e.GET("/*", shortener.Handler(linkStorage))

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
}
