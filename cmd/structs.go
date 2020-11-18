package cmd

import "errors"

// Mysql describes command-line arguments related to Mysql storage
type Mysql struct {
	Host     string `long:"mysql-host" description:"mysql storage DB host with port" default:"localhost:3306" env:"MYSQL_HOST"`
	Name     string `long:"mysql-dbname" description:"mysql storage DB name" env:"MYSQL_DBNAME"`
	User     string `long:"mysql-user" description:"mysql storage DB user" env:"MYSQL_USER"`
	Password string `long:"mysql-password" description:"mysql storage DB password" env:"MYSQL_PASSWORD"`
}

// Validate validates Mysql storage arguments
func (m Mysql) Validate() error {
	if m.Host == "" {
		return errors.New("MySQL host is not set")
	}
	if m.Name == "" {
		return errors.New("MySQL db name is not set")
	}
	if m.User == "" {
		return errors.New("MySQL user is not set")
	}
	if m.Password == "" {
		return errors.New("MySQL password is not set")
	}
	return nil
}
