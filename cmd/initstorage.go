package cmd

import (
	"github.com/denisvmedia/urlshortener/storage/linkstorage"
	"github.com/jessevdk/go-flags"
)

// RegisterInitStorageCommand registers `init-storage` command
func RegisterInitStorageCommand(parser *flags.Parser) *InitStorageCommand {
	cmd := &InitStorageCommand{}
	_, err := parser.AddCommand("init-storage", "initializes the selected storage", "", cmd)
	if err != nil {
		panic(err)
	}
	return cmd
}

// InitStorageCommand defines `init-storage` command
type InitStorageCommand struct {
	Storage        string `long:"storage" description:"storage to use" choice:"mysql" default:"mysql" env:"STORAGE"`
	CreateDatabase bool   `long:"create-database" description:"will DROP (!) and create the database, use CAREFULLY!"`
	Mysql
}

// Execute implements `init-storage` command
func (cmd *InitStorageCommand) Execute(_ []string) error {
	if err := cmd.Mysql.Validate(); err != nil {
		return err
	}

	err := linkstorage.MysqlInitStorage(cmd.Mysql.User, cmd.Mysql.Password, cmd.Mysql.Host, cmd.Mysql.Name, cmd.CreateDatabase)
	if err != nil {
		return err
	}

	return nil
}
