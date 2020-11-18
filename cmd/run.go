package cmd

import (
	"fmt"
	"github.com/denisvmedia/urlshortener/metrics"
	"github.com/denisvmedia/urlshortener/server"
	"github.com/denisvmedia/urlshortener/storage/linkstorage"
	"github.com/jessevdk/go-flags"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// RegisterRunCommand registers `run` command
func RegisterRunCommand(parser *flags.Parser) *RunCommand {
	cmd := &RunCommand{}
	_, err := parser.AddCommand("run", "runs url shortener web server daemon", "", cmd)
	if err != nil {
		panic(err)
	}
	return cmd
}

// RunCommand defines `run` command
type RunCommand struct {
	BindAddress string `long:"bind-address" description:"http bind address" default:":31456" env:"BIND_ADDRESS"`
	Storage     string `long:"storage" description:"storage to use" choice:"mysql" choice:"inmemory" default:"inmemory" env:"STORAGE"`
	Mysql
}

func setUpGracefulExit(server *http.Server) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		// graceful exit on ctrl-c
		<-c
		_ = server.Close()
	}()
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

	metrics.RegisterAll()
	e := server.NewEcho(linkStorage)
	fmt.Printf("Listening on %s\n", cmd.BindAddress)
	setUpGracefulExit(e.Server)
	e.Logger.Fatal(e.Start(cmd.BindAddress))

	return nil
}
