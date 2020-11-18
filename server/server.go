package server

import (
	"github.com/denisvmedia/urlshortener/model"
	"github.com/denisvmedia/urlshortener/resource"
	"github.com/denisvmedia/urlshortener/routing"
	"github.com/denisvmedia/urlshortener/shortener"
	"github.com/denisvmedia/urlshortener/storage/linkstorage"
	"github.com/go-extras/api2go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// NewEcho create a new API router
func NewEcho(linkStorage linkstorage.Storage) *echo.Echo {
	e := echo.New()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	api := api2go.NewAPIWithRouting(
		"api",
		api2go.NewStaticResolver("/"),
		routing.Echo(e),
	)

	api.AddResource(model.Link{}, resource.NewLinkResource(linkStorage))

	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
	e.GET("/swagger/*any", echoSwagger.EchoWrapHandler(echoSwagger.URL("/swagger/doc.json")))
	e.GET("/*", shortener.Handler(linkStorage))

	return e
}
