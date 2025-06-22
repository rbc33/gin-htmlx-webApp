package database

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	// "os"
	"time"

	// "github.com/joho/godotenv"
	"github.com/rbc33/gocms/common"
	"github.com/rs/zerolog/log"
	"github.com/xeipuuv/gojsonschema"
)

type Database interface {
	GetPosts(offset int, limit int) ([]common.Post, error)
	GetPost(post_id int) (common.Post, error)
	AddPost(title string, excerpt string, content string) (int, error)
	ChangePost(id int, title string, excerpt string, content string) error
	DeletePost(id int) error
	AddImage(uuid string, name string, alt string) error
	DeleteImage(uuid string) error
	AddCard(uuid string, image_location string, json_data string, schema_name string) error
	GetCard(uuid string) (common.Card, error)
	ChangeCard(uuid string, image_location string, json_data string, schema_name string) error
	DeleteCard(uuid string) error
	GetPages(offset int, limit int) ([]common.Page, error)
	AddPage(title string, content string, link string) (int, error)
	GetPage(link string) (common.Page, error)
}

type SqlDatabase struct {
	MY_SQL_URL string
	Connection *sql.DB
}

// / GetPosts gets all the posts from the current
// / database connection.
func (db *SqlDatabase) GetPosts(limit int, offset int) ([]common.Post, error) {
	all_posts := make([]common.Post, 0)
	var rows *sql.Rows
	var err error

	query := "SELECT title, excerpt, id FROM posts"
	args := make([]interface{}, 0)

	// A limit of 0 or less means no limit.
	if limit > 0 {
		query += " LIMIT ?"
		args = append(args, limit)
		// OFFSET is only valid with LIMIT
		if offset > 0 {
			query += " OFFSET ?"
			args = append(args, offset)
		}
	}
	query += ";"

	rows, err = db.Connection.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post common.Post
		if err = rows.Scan(&post.Title, &post.Excerpt, &post.Id); err != nil {
			return make([]common.Post, 0), err
		}
		all_posts = append(all_posts, post)
	}

	return all_posts, rows.Err()
}

// return post by id
func (db *SqlDatabase) GetPost(post_id int) (post common.Post, err error) {
	rows, err := db.Connection.Query("SELECT id, title, content FROM posts WHERE id=?;", post_id)
	if err != nil {
		return common.Post{}, err
	}
	defer rows.Close() // QueryRow doesn't need this, but it doesn't hurt

	if err = rows.Scan(&post.Id, &post.Title, &post.Content); err != nil {
		return common.Post{}, err
	}

	return post, nil
}

// AddPost adds a post to the database
func (db *SqlDatabase) AddPost(title string, excerpt string, content string) (Id int, err error) {
	res, err := db.Connection.Exec("INSERT INTO posts(content, title, excerpt) VALUES(?, ?, ?)", content, title, excerpt)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Warn().Msgf("could not get last ID: %v", err)
		return -1, nil
	}

	// TODO : possibly unsafe int conv,
	// make sure all IDs are i64 in the
	// future
	return int(id), nil
}

// ChangePost changes a post based on the values
// provided. Note that empty strings will mean that
// the value will not be updated.
func (db *SqlDatabase) ChangePost(id int, title string, excerpt string, content string) (err error) {
	tx, err := db.Connection.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if len(title) > 0 {
		_, err := tx.Exec("UPDATE posts SET title = ? WHERE id = ?;", title, id)
		if err != nil {
			return err
		}
	}

	if len(excerpt) > 0 {
		_, err := tx.Exec("UPDATE posts SET excerpt = ? WHERE id = ?;", excerpt, id)
		if err != nil {
			return err
		}
	}

	if len(content) > 0 {
		_, err := tx.Exec("UPDATE posts SET content = ? WHERE id = ?;", content, id)
		if err != nil {
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

// DeletePost changes a post based on the values
// provided. Note that empty strings will mean that
// the value will not be updated.
func (db *SqlDatabase) DeletePost(id int) error {
	if _, err := db.Connection.Exec("DELETE FROM posts WHERE id=?;", id); err != nil {
		return err
	}

	return nil
}

// AddImage will add the image metadata to the db.
// name - file name saved to disk
// alt - alternative text
// return(uuid, nil) if succeeded, ("", err) otherwise
func (db *SqlDatabase) AddImage(uuid string, name string, alt string) error {
	tx, err := db.Connection.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if name == "" {
		return fmt.Errorf("cannot have empty names")
	}
	if alt == "" {
		return fmt.Errorf("cannot have empty alt text")
	}

	query := "INSERT INTO images(uuid, name, alt) VALUES (?, ?, ?);"
	_, err = tx.Exec(query, uuid, name, alt)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

// / This function gets a post from the database
// / with the given ID.
func (db *SqlDatabase) GetCard(uuid string) (card common.Card, err error) {
	rows, err := db.Connection.Query("SELECT image_location, json_data, json_schema FROM cards WHERE uuid=?;", uuid)
	if err != nil {
		return common.Card{}, err
	}
	defer func() {
		err = errors.Join(rows.Close())
	}()

	rows.Next()
	if err = rows.Scan(&card.ImageLocation, &card.JsonData, &card.SchemaName); err != nil {
		return common.Card{}, err
	}

	// Validate the json
	validateJson(card.JsonData, card.SchemaName)

	return card, nil
}

// DeletePost changes a post based on the values
// provided. Note that empty strings will mean that
// the value will not be updated.
func (db *SqlDatabase) DeleteImage(uuid string) error {
	if _, err := db.Connection.Exec("DELETE FROM image WHERE uuid=?;", uuid); err != nil {
		return err
	}

	return nil
}

func (db *SqlDatabase) AddCard(uuid string, image_location string, json_data string, schema_name string) error {

	log.Info().Msgf("adding card to the DB")

	// Check the file exist and is a file
	// not a directory
	if image_location == "" {
		return fmt.Errorf("cannot have empty image")
	}
	image_stat, err := os.Stat(image_location)
	if errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("file does not exist: %s", image_location)
	}
	image_location = "/" + image_location
	if err != nil {
		return err
	}
	if image_stat.IsDir() {
		return fmt.Errorf("given path is a directory: %s", image_stat)
	}

	if json_data == "" {
		return fmt.Errorf("cannot have empty data")
	}

	// Load schema
	if schema_name == "" {
		return fmt.Errorf("cannot have empty data")
	}

	err = validateJson(json_data, schema_name)
	if err != nil {
		return err
	}

	tx, err := db.Connection.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if commit_err := tx.Commit(); commit_err != nil {
			err = errors.Join(err, tx.Rollback(), commit_err)
		}
	}()

	query := "INSERT INTO cards(uuid, image_location, json_data, json_schema) VALUES (?, ?, ?, ?);"
	_, err = tx.Exec(query, uuid, image_location, json_data, schema_name)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (db *SqlDatabase) ChangeCard(uuid string, image_location string, json_data string, schema_name string) error {

	log.Info().Msgf("changing card")

	tx, err := db.Connection.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if commit_err := tx.Commit(); commit_err != nil {
			err = errors.Join(err, tx.Rollback(), commit_err)
		}
	}()

	if image_location != "" {
		image_stat, err := os.Stat(image_location)
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("file does not exist: %s", image_location)
		}
		if err != nil {
			return err
		}
		if image_stat.IsDir() {
			return fmt.Errorf("given path is a directory: %s", image_stat)
		}
		_, err = tx.Exec("UPDATE cards SET image_location = ? WHERE uuid = ?;", image_location, uuid)
		if err != nil {
			return err
		}
	}

	if json_data != "" && schema_name != "" {

		err = validateJson(json_data, schema_name)
		if err != nil {
			return err
		}
		_, err = tx.Exec("UPDATE cards SET json_data = ?, json_schema = ? WHERE uuid = ?;", json_data, schema_name, uuid)
		if err != nil {
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (db *SqlDatabase) DeleteCard(uuid string) error {
	if _, err := db.Connection.Exec("DELETE FROM cards WHERE uuid=?;", uuid); err != nil {
		return err
	}

	return nil
}

func (db *SqlDatabase) AddPage(title string, content string, link string) (int, error) {
	res, err := db.Connection.Exec("INSERT INTO pages(content, title, link) VALUES(?, ?, ?)", content, title, link)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Warn().Msgf("could not get last ID: %v", err)
		return -1, nil
	}

	// TODO : possibly unsafe int conv,
	// make sure all IDs are i64 in the
	// future
	return int(id), nil
}

func (db *SqlDatabase) GetPages(limit int, offset int) ([]common.Page, error) {
	all_pages := make([]common.Page, 0)
	var rows *sql.Rows
	var err error

	query := "SELECT title, content, link, id FROM pages"
	args := make([]interface{}, 0)

	// A limit of 0 or less means no limit.
	if limit > 0 {
		query += " LIMIT ?"
		args = append(args, limit)
		// OFFSET is only valid with LIMIT
		if offset > 0 {
			query += " OFFSET ?"
			args = append(args, offset)
		}
	}
	query += ";"

	rows, err = db.Connection.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var page common.Page
		if err = rows.Scan(&page.Title, &page.Content, &page.Link, &page.Id); err != nil {
			return nil, err
		}
		all_pages = append(all_pages, page)
	}

	return all_pages, rows.Err()
}

func (db *SqlDatabase) GetPage(link string) (common.Page, error) {
	query := "SELECT id, title, content, link FROM pages WHERE link=?;"
	row := db.Connection.QueryRow(query, link)
	var page common.Page
	if err := row.Scan(&page.Id, &page.Title, &page.Content, &page.Link); err != nil {
		return common.Page{}, err
	}

	return page, nil
}

func validateJson(json_data string, schema_name string) error {
	schema_data, err := os.ReadFile(filepath.Join("schemas", schema_name+".json"))
	if err != nil {
		return fmt.Errorf("%v, %v", schema_name, err)
	}

	schemaLoader := gojsonschema.NewBytesLoader([]byte(schema_data))
	documentLoader := gojsonschema.NewStringLoader(json_data) // Changed from NewReferenceLoader to NewStringLoader
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return fmt.Errorf("could not validate json_data: %v", err)
	}
	if !result.Valid() {
		var errors []string
		for _, err := range result.Errors() {
			errors = append(errors, err.String())
		}
		return fmt.Errorf("invalid card json: %s", strings.Join(errors, "; "))
	}
	return nil
}

func MakeSqlConnection(appSettings common.AppSettings) (database SqlDatabase, err error) {
	/// Checking the DB connection
	// err := godotenv.Load()
	// if err != nil {
	// 	return SqlDatabase{}, err
	// }
	connection_str := appSettings.DatabaseUri
	db, err := sql.Open("mysql", connection_str)
	if err != nil {
		return SqlDatabase{}, err
	}
	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Second * 5)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	log.Info().Msg(connection_str)

	return SqlDatabase{
		MY_SQL_URL: connection_str,
		Connection: db,
	}, nil
}
