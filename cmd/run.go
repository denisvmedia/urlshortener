package cmd

import (
	"fmt"
	"github.com/denisvmedia/urlshortener/metrics"
	"github.com/denisvmedia/urlshortener/model"
	"github.com/denisvmedia/urlshortener/resource"
	"github.com/denisvmedia/urlshortener/routing"
	"github.com/denisvmedia/urlshortener/shortener"
	"github.com/denisvmedia/urlshortener/storage/linkstorage"
	"github.com/go-extras/api2go"
	"github.com/jessevdk/go-flags"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// RegisterInitStorageCommand registers `run` command
func RegisterRunCommand(parser *flags.Parser) *RunCommand {
	cmd := &RunCommand{}
	_, err := parser.AddCommand("run", "runs url shortener web server daemon", "", cmd)
	if err != nil {
		panic(err)
	}
	return cmd
}

// InitStorageCommand defines `run` command
type RunCommand struct {
	BindAddress string `long:"bind-address" description:"http bind address" default:":31456" env:"BIND_ADDRESS"`
	Storage     string `long:"storage" description:"storage to use" choice:"mysql" choice:"inmemory" default:"inmemory" env:"STORAGE"`
	Mysql
}

// Execute implements `run` command
func (cmd *RunCommand) Execute(_ []string) error {
	var linkStorage linkstorage.Storage
	if cmd.Storage == "mysql" {
		if err := cmd.Mysql.Validate(); err != nil {
			return err
		}
		fmt.Println("Storing all data in MySQL.")
		dbh, err := linkstorage.MysqlConnect(cmd.Mysql.User, cmd.Mysql.Password, cmd.Mysql.Host, cmd.Mysql.Name)
		if err != nil {
			return err
		}
		linkStorage = linkstorage.NewMysqlStorage(dbh)
	} else {
		fmt.Println("Storing all data in memory. All your activity will be lost after you stop the application.")
		linkStorage = linkstorage.NewInMemoryStorage()
	}

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

	metrics.RegisterAll()
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
	e.GET("/swagger/*any", echoSwagger.EchoWrapHandler(echoSwagger.URL("/swagger/doc.json")))
	e.GET("/*", shortener.Handler(linkStorage))

	fmt.Printf("Listening on %s\n", cmd.BindAddress)
	e.Logger.Fatal(e.Start(cmd.BindAddress))
	return nil
}
