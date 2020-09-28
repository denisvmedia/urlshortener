package routing

import (
	"github.com/go-extras/api2go/routing"
	"net/http"

	"github.com/labstack/echo/v4"
)

type echoRouter struct {
	echo *echo.Echo
}

func (e echoRouter) Handler() http.Handler {
	return e.echo
}

func (e echoRouter) Handle(protocol, route string, handler routing.HandlerFunc) {
	// fmt.Printf("registering echo route [%s] %s\n", protocol, route)
	echoHandlerFunc := func(c echo.Context) error {
		params := map[string]string{}

		for i, p := range c.ParamNames() {
			params[p] = c.ParamValues()[i]
		}

		handler(c.Response(), c.Request(), params, make(map[string]interface{}))

		return nil
	}
	e.echo.Add(protocol, route, echoHandlerFunc)
}

// Echo created a new api2go router to use with the echo framework
func Echo(e *echo.Echo) routing.Routeable {
	return &echoRouter{echo: e}
}
