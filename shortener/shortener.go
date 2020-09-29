package shortener

import (
	"github.com/denisvmedia/urlshortener/metrics"
	"github.com/denisvmedia/urlshortener/storage/linkstorage"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	// using this package because of unsolved issue in golang.org/x/net/http: https://github.com/golang/go/issues/19307
	"github.com/golang/gddo/httputil"
)

const resourceNotFound = "resource not found"
const pageError = `
<!doctype html>
<html class="no-js" lang="en">
    <head>
        <meta charset="utf-8">
        <meta http-equiv="x-ua-compatible" content="IE=edge,chrome=1">
        <meta name="viewport" content="width=device-width, initial-scale=1">

        <title>Error 404 - Not Found!</title>
        <meta name="robots" content="noindex">

        <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/twitter-bootstrap/3.3.7/css/bootstrap.min.css">

        <style>
          h1.error {
            margin-top: 1em;
            font-size: 7em;
            font-weight: 500;
          }
        </style>

        <script src="https://cdnjs.cloudflare.com/ajax/libs/modernizr/2.8.3/modernizr.min.js"></script>
    </head>
    <body>

        <div class="container">
            <div class="row">
                <div class="col-md-10 col-md-offset-2">
                    <h1 class="error">Oops!</h1>
                    <h2>Error 404 - Not Found</h2>
                    <p class="lead">Sorry, an error has occured. The resource you requested has not been found!</p>
                    <a href="#" class="btn btn-primary btn-md"><span class="glyphicon glyphicon-home"></span> Home </a> <a href="#" class="btn btn-default btn-md"><span class="glyphicon glyphicon-envelope"></span> Contact </a>
                </div>
            </div>
        </div>

    </body>
</html>`

// Handler Handle short link redirection
func Handler(linkStorage linkstorage.Storage) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		shortName := strings.Trim(ctx.Param("*"), "/ ")
		link, err := linkStorage.GetOneByShortName(shortName)
		if err == nil {
			metrics.RequestProcessed.WithLabelValues("301").Inc()
			return ctx.Redirect(http.StatusMovedPermanently, link.OriginalUrl)
		}

		contentType := httputil.NegotiateContentType(ctx.Request(), []string{"text/plain", "text/html", "application/json", "application/vnd.api+json"}, "")
		metrics.RequestProcessed.WithLabelValues("404").Inc()

		switch contentType {
		case "application/json", "application/vnd.api+json":
			return ctx.JSON(404, map[string]interface{}{
				"errors": map[string]interface{}{
					"status": 404,
					"title":  resourceNotFound,
				},
			})
		case "text/html":
			return ctx.HTML(404, pageError)
		}

		return ctx.String(404, resourceNotFound)
	}
}
