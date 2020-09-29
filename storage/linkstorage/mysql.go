package linkstorage

import (
	"fmt"
	"github.com/denisvmedia/urlshortener/model"
	"github.com/denisvmedia/urlshortener/storage"
	"github.com/go-extras/errors"
	_ "github.com/go-sql-driver/mysql"
	"time"

	"github.com/jmoiron/sqlx"
)

// NewMysqlStorage initializes the storage
func NewMysqlStorage(db *sqlx.DB) Storage {
	return &MysqlStorage{
		db: db,
	}
}

type MysqlStorage struct {
	db *sqlx.DB
}

func (m *MysqlStorage) countAll() (count int, err error) {
	query := "SELECT COUNT(ID) FROM links"
	stmt, err := m.db.Prepare(query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	if !rows.Next() {
		return 0, storage.ErrStorageFailure
	}

	err = rows.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (m *MysqlStorage) PaginatedGetAll(pageNumber, pageSize int) (results []*model.Link, total int, err error) {
	offset := (pageNumber - 1) * pageSize
	limit := pageSize

	cnt, err := m.countAll()
	if err != nil {
		return nil, 0, err
	}

	query := "SELECT id, short_name, original_url, comment FROM links LIMIT ?, ?"
	stmt, err := m.db.Prepare(query)
	if err != nil {
		return nil, 0, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(offset, limit)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var idNew int
		var shortNameNew, originalUrl, comment string
		err = rows.Scan(&idNew, &shortNameNew, &originalUrl, &comment)
		if err != nil {
			return nil, 0, err
		}

		results = append(results, &model.Link{
			ID:          fmt.Sprint(idNew),
			ShortName:   shortNameNew,
			OriginalUrl: originalUrl,
			Comment:     comment,
		})
	}

	return results, cnt, nil
}

func (m *MysqlStorage) GetOne(id string) (*model.Link, error) {
	query := "SELECT id, short_name, original_url, comment FROM links WHERE id=?"
	stmt, err := m.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, storage.ErrNotFound
	}

	var idNew int
	var shortNameNew, originalUrl, comment string
	err = rows.Scan(&idNew, &shortNameNew, &originalUrl, &comment)
	if err != nil {
		return nil, err
	}

	return &model.Link{
		ID:          fmt.Sprint(idNew),
		ShortName:   shortNameNew,
		OriginalUrl: originalUrl,
		Comment:     comment,
	}, nil
}

func (m *MysqlStorage) GetOneByShortName(shortName string) (*model.Link, error) {
	query := "SELECT id, short_name, original_url, comment FROM links WHERE short_name=?"
	stmt, err := m.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(shortName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, storage.ErrNotFound
	}

	var id int
	var shortNameNew, originalUrl, comment string
	err = rows.Scan(&id, &shortNameNew, &originalUrl, &comment)
	if err != nil {
		return nil, err
	}

	return &model.Link{
		ID:          fmt.Sprint(id),
		ShortName:   shortNameNew,
		OriginalUrl: originalUrl,
		Comment:     comment,
	}, nil
}

func (m *MysqlStorage) Insert(c model.Link) (*model.Link, error) {
	if c.ShortName == "" {
		c.ShortName = generateShortName()
	}

	existing, err := m.GetOneByShortName(c.ShortName)
	if err != nil && err != storage.ErrNotFound {
		return nil, err
	}
	if existing != nil && existing.ID != c.ID {
		return existing, errors.Wrapf(storage.ErrShortNameAlreadyExists, "Existing link id %s", existing.ID)
	}

	query := "INSERT INTO links (short_name, original_url, comment, created_at, updated_at) VALUES (?, ?, ?, ?, ?)"
	stmt, err := m.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	created := time.Now()
	result, err := stmt.Exec(c.ShortName, c.OriginalUrl, c.Comment, created, created)
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()
	if id <= 0 {
		return nil, errors.Wrapf(storage.ErrStorageFailure, "Got non-positive last insert id")
	}

	c.ID = fmt.Sprint(id)

	return &c, nil
}

func (m *MysqlStorage) Delete(id string) error {
	query := "DELETE FROM links WHERE id = ?"
	stmt, err := m.db.Prepare(query)
	if err != nil {
		panic(err) // here we intentionally panic so that the program is forced to restart
	}
	defer stmt.Close()

	result, err := stmt.Exec(id)
	if err != nil {
		return err
	}

	cnt, _ := result.RowsAffected()
	if cnt == 0 {
		return storage.ErrNotFound
	}

	return nil
}

func (m *MysqlStorage) Update(c model.Link) error {
	_, err := m.GetOne(c.ID)
	if err != nil {
		return err
	}

	existing, err := m.GetOneByShortName(c.ShortName)
	if err != nil && err != storage.ErrNotFound {
		return err
	}
	if existing != nil && existing.ID != c.ID {
		return errors.Wrapf(storage.ErrShortNameAlreadyExists, "Existing link id %s", existing.ID)
	}

	query := "UPDATE links SET short_name = ?, original_url = ?, comment = ?, updated_at = ? WHERE id = ?"
	stmt, err := m.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(c.ShortName, c.OriginalUrl, c.Comment, time.Now(), c.ID)
	if err != nil {
		return err
	}

	cnt, _ := result.RowsAffected()
	if cnt == 0 {
		return errors.Wrapf(storage.ErrStorageFailure, "DB reported no rows changed")
	}

	return nil
}

func mysqlCreateDB(dbUser, dbPassword, dbHost, dbName string) error {
	dbh, err := sqlx.Connect("mysql",
		fmt.Sprintf("%s:%s@(%s)/?parseTime=true",
			dbUser, dbPassword, dbHost))
	if err != nil {
		return err
	}
	defer dbh.Close()

	_, err = dbh.Queryx("DROP DATABASE IF EXISTS " + dbName)
	if err != nil {
		return err
	}
	_, err = dbh.Queryx("CREATE DATABASE IF NOT EXISTS " + dbName)
	if err != nil {
		return err
	}

	return nil
}

func MysqlConnect(dbUser, dbPassword, dbHost, dbName string) (*sqlx.DB, error) {
	return sqlx.Connect("mysql",
		fmt.Sprintf("%s:%s@(%s)/%s?parseTime=true",
			dbUser, dbPassword, dbHost, dbName))
}

func MysqlInitStorage(dbUser, dbPassword, dbHost, dbName string, createDb bool) error {
	if createDb {
		err := mysqlCreateDB(dbUser, dbPassword, dbHost, dbName)
		if err != nil {
			return err
		}
	}

	dbh, err := sqlx.Connect("mysql",
		fmt.Sprintf("%s:%s@(%s)/%s?parseTime=true",
			dbUser, dbPassword, dbHost, dbName))
	if err != nil {
		return err
	}

	_, err = dbh.Queryx("CREATE TABLE `links` (`id` INT NOT NULL AUTO_INCREMENT, " +
		"`short_name` VARCHAR(255) NOT NULL, " +
		"`original_url` TEXT NOT NULL, " +
		"`comment` VARCHAR(255) NOT NULL, " +
		"`created_at` DATETIME NOT NULL, " +
		"`updated_at` DATETIME NOT NULL, " +
		"PRIMARY KEY (`id`), " +
		"UNIQUE INDEX `short_name` (`short_name`)) " +
		"COLLATE='utf8_general_ci'")
	if err != nil {
		return err
	}

	return nil
}

func MysqlDropDB(dbUser, dbPassword, dbHost, dbName string) error {
	dbh, err := sqlx.Connect("mysql",
		fmt.Sprintf("%s:%s@(%s)/?parseTime=true",
			dbUser, dbPassword, dbHost))
	if err != nil {
		return err
	}
	defer dbh.Close()

	_, err = dbh.Queryx("DROP DATABASE IF EXISTS " + dbName)
	if err != nil {
		return err
	}

	return nil
}
