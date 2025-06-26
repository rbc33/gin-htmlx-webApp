package database

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	// "os"
	"time"

	// "github.com/joho/godotenv"
	"github.com/google/uuid"
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
	// GetCard(uuid string) (common.Card, error)
	AddCard(image string, schema string, content string) (string, error)
	GetCards(schema_uuid string, limit int, page int) ([]common.Card, error)
	AddCardSchema(json_schema string, json_title string) (string, error)
	GetCardSchemas(offset int, limit int) ([]common.CardSchema, error)
	GetCardSchema(uuid string) (common.CardSchema, error)
	DeleteCardSchema(uuid string) error
	ChangeCard(uuid string, image_location string, json_data string, schema_name string) error
	DeleteCard(uuid string) error
	GetPages(offset int, limit int) ([]common.Page, error)
	AddPage(title string, content string, link string) (int, error)
	GetPage(link string) (common.Page, error)
	ChangePage(id int, title string, content string, link string) error
	DeletePage(link string) error
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

// / This function gets a post from the database
// / with the given ID.
func (db SqlDatabase) GetPost(post_id int) (post common.Post, err error) {
	rows, err := db.Connection.Query("SELECT id, title, content, excerpt FROM posts WHERE id=?;", post_id)
	if err != nil {
		return common.Post{}, err
	}
	defer func() {
		err = errors.Join(rows.Close())
	}()

	rows.Next()
	if err = rows.Scan(&post.Id, &post.Title, &post.Content, &post.Excerpt); err != nil {
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

// // / This function gets a post from the database
// // / with the given ID.
// func (db *SqlDatabase) GetCard(uuid string) (card common.Card, err error) {
// 	rows, err := db.Connection.Query("SELECT title, image, schema, content FROM cards WHERE uuid=?;", uuid)
// 	if err != nil {
// 		return common.Card{}, err
// 	}
// 	defer func() {
// 		err = errors.Join(rows.Close())
// 	}()

// 	rows.Next()
// 	if err = rows.Scan(&card.Title, &card.Image, &card.Schema, &card.Content ); err != nil {
// 		return common.Card{}, err
// 	}

// 	// Validate the json
// 	validateJson(card., card.SchemaName)

// 	return card, nil
// }

// DeletePost changes a post based on the values
// provided. Note that empty strings will mean that
// the value will not be updated.
func (db *SqlDatabase) DeleteImage(uuid string) error {
	if _, err := db.Connection.Exec("DELETE FROM image WHERE uuid=?;", uuid); err != nil {
		return err
	}

	return nil
}

// / This function adds the card metadata to the cards table.
// / Returns the uuid as a string if successful, otherwise error
// / won't be null
func (db SqlDatabase) AddCard(image string, schema_uuid string, content string) (string, error) {

	schema, err := db.GetCardSchema(schema_uuid)
	if err != nil {
		return "", err
	}

	tx, err := db.Connection.Begin()
	if err != nil {
		return "", err
	}
	defer func() {
		if commit_err := tx.Commit(); commit_err != nil {
			err = errors.Join(err, tx.Rollback(), commit_err)
		}
	}()

	uuid := uuid.New().String()

	_, err = tx.Exec("INSERT INTO cards(uuid, image_location, json_data, json_schema) VALUES(UuidToBin(?), ?, ?, ?)", uuid, image, content, schema_uuid)
	if err != nil {
		return "", err
	}

	schema.Cards = append(schema.Cards, uuid)
	cards_string, err := json.Marshal(schema.Cards)
	fmt.Printf("UPDATE card_schemas SET card_ids = %s WHERE uuid = UuidToBin(%s);", cards_string, schema_uuid)

	_, err = tx.Exec("UPDATE card_schemas SET card_ids = ? WHERE uuid = UuidToBin(?);", cards_string, schema_uuid)
	if err != nil {
		return "", err
	}

	return uuid, nil
}

func (db SqlDatabase) GetCards(schema_uuid string, limit int, page int) (all_cards []common.Card, err error) {

	// TODO : need to get the card within the current schema
	schema, err := db.GetCardSchema(schema_uuid)
	if err != nil {
		return []common.Card{}, err
	}

	// Some hand-rolled logic to paginate
	// we want the first "limit" cards after the "page * limit"
	// cards, if available.

	if len(schema.Cards) <= (page * limit) {
		return []common.Card{}, fmt.Errorf("no cards available for the pagination given")
	}

	window := schema.Cards[(page * limit):min(len(schema.Cards), (page+1)*limit)]

	// gets the right UuidToBin(?), UuidToBin(?), ..., UuidToBin(?) for the IN clause
	ids_query := strings.Repeat("UuidToBin(?), ", len(window))
	ids_query = ids_query[:len(ids_query)-2]

	// Convert slice to interface{} slice for passing to Query
	args := make([]interface{}, len(window))
	for i, v := range window {
		args[i] = v
	}

	query := fmt.Sprintf("SELECT UuidFromBin(uuid), image_location, json_data, json_schema FROM cards WHERE uuid IN (%s);", ids_query)
	rows, err := db.Connection.Query(query, args...)

	if err != nil {
		return []common.Card{}, err
	}
	defer func() {
		err = errors.Join(rows.Close())
	}()

	for rows.Next() {
		var card common.Card
		if err = rows.Scan(&card.Id, &card.Image, &card.Content, &card.Schema); err != nil {
			return []common.Card{}, err
		}
		all_cards = append(all_cards, card)
	}

	return all_cards, nil
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
		_, err = tx.Exec("UPDATE cards SET image_location = ? WHERE uuid = UuidFromBin(?9;", image_location, uuid)
		if err != nil {
			return err
		}
	}

	if json_data != "" && schema_name != "" {

		err = validateJson(json_data, schema_name)
		if err != nil {
			return err
		}
		_, err = tx.Exec("UPDATE cards SET json_data = ?, json_schema = ? WHERE uuid = UuidFromBin(?);", json_data, schema_name, uuid)
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
	if _, err := db.Connection.Exec("DELETE FROM cards WHERE uuid=UuidFromBin(?);", uuid); err != nil {
		return err
	}

	return nil
}

func (db SqlDatabase) AddCardSchema(json_schema string, json_title string) (string, error) {
	uuid := uuid.New().String()

	_, err := db.Connection.Exec(
		"INSERT INTO card_schemas(uuid, json_id, json_schema, json_title, card_ids) VALUES(UuidToBin(?), ?, ?, ?, ?)",
		uuid,
		"some_id",
		json_schema,
		json_title,
		"[]")

	if err != nil {
		return "", err
	}

	return uuid, nil
}

func (db SqlDatabase) GetCardSchema(id string) (schema common.CardSchema, err error) {
	rows, err := db.Connection.Query("SELECT json_schema, json_title, card_ids FROM card_schemas WHERE uuid=UuidToBin(?);", id)
	if err != nil {
		return common.CardSchema{}, err
	}

	defer func() {
		err = errors.Join(rows.Close())
	}()

	rows.Next()
	var card_ids_string string
	if err = rows.Scan(&schema.Schema, &schema.Title, &card_ids_string); err != nil {
		return common.CardSchema{}, err
	}

	// We need to parse the schemas here
	err = json.Unmarshal([]byte(card_ids_string), &schema.Cards)
	if err != nil {
		return common.CardSchema{}, fmt.Errorf("can't parse card ids json: %v", err)
	}

	return schema, nil
}

func (db *SqlDatabase) DeleteCardSchema(uuid string) error {
	if _, err := db.Connection.Exec("DELETE FROM card_schemas WHERE uuid=UuidFromBin(?);", uuid); err != nil {
		return err
	}

	return nil
}

func (db SqlDatabase) GetCardSchemas(offset int, limit int) (schemas []common.CardSchema, err error) {
	all_schemas := []common.CardSchema{}
	var rows *sql.Rows

	query := "SELECT UuidFromBin(uuid), json_schema, json_title, card_ids FROM card_schemas"
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
		var schema common.CardSchema
		var cardsRaw []byte

		if err = rows.Scan(&schema.Uuid, &schema.Schema, &schema.Title, &cardsRaw); err != nil {
			return []common.CardSchema{}, err
		}

		// Aquí parseas el JSON de la base a tu array de Go:
		if len(cardsRaw) > 0 {
			if err = json.Unmarshal(cardsRaw, &schema.Cards); err != nil {
				return []common.CardSchema{}, fmt.Errorf("error unmarshaling card_ids: %w", err)
			}
		}
		all_schemas = append(all_schemas, schema)

	}
	return all_schemas, rows.Err()
}

func (db *SqlDatabase) AddPage(title string, content string, link string) (int, error) {
	res, err := db.Connection.Exec("INSERT INTO pages(content, title, link) VALUES(?, ?, ?);", content, title, link)
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

func (db *SqlDatabase) ChangePage(id int, title string, content string, link string) (err error) {
	tx, err := db.Connection.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if len(title) > 0 {
		_, err := tx.Exec("UPDATE pages SET title = ? WHERE id = ?;", title, id)
		if err != nil {
			return err
		}
	}

	if len(link) > 0 {
		_, err := tx.Exec("UPDATE pages SET link = ? WHERE id = ?;", link, id)
		if err != nil {
			return err
		}
	}

	if len(content) > 0 {
		_, err := tx.Exec("UPDATE pages SET content = ? WHERE id = ?;", content, id)
		if err != nil {
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (db *SqlDatabase) DeletePage(link string) error {
	if _, err := db.Connection.Exec("DELETE FROM pages WHERE link=?;", link); err != nil {
		return err
	}

	return nil
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
