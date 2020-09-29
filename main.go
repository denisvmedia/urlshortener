package main

import (
	"github.com/denisvmedia/urlshortener/cmd"
	_ "github.com/denisvmedia/urlshortener/docs"
	"github.com/jessevdk/go-flags"
	"os"
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
	parser := flags.NewParser(nil, flags.Default)
	cmd.RegisterInitStorageCommand(parser)
	cmd.RegisterRunCommand(parser)

	//parser.CommandHandler = func(command flags.Commander, args []string) error {
	//	err := command.Execute(args)
	//
	//	return err
	//}

	_, err := parser.Parse()
	if err != nil {
		os.Exit(-1)
	}

	return
}
