package database

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"gocms/common"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type Database struct {
	MY_SQL_URL string
	Connection *sql.DB
}

func GetDatabaseURL() (string, error) {
	// En producción (Clever Cloud)
	if os.Getenv("ENVIRONMENT") == "production" {
		return buildCleverCloudMySQLURL(), nil
	}

	// En desarrollo (tu configuración local)
	err := godotenv.Load()
	if err != nil {
		return "", err
	}
	return os.Getenv("MY_SQL_URL"), nil
}

func buildCleverCloudMySQLURL() string {
	host := os.Getenv("MYSQL_ADDON_HOST")
	port := os.Getenv("MYSQL_ADDON_PORT")
	user := os.Getenv("MYSQL_ADDON_USER")
	password := os.Getenv("MYSQL_ADDON_PASSWORD")
	database := os.Getenv("MYSQL_ADDON_DB")

	if host == "" || user == "" || password == "" || database == "" {
		// Fallback a variable personalizada si existe
		if customURL := os.Getenv("MY_SQL_URL"); customURL != "" {
			return customURL
		}
		panic("Database configuration not found")
	}

	// Formato directo para Go MySQL driver: user:password@tcp(host:port)/database
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		user, password, host, port, database)
}

// / GetPosts gets all the posts from the current
// / database connection.
func (db Database) GetPosts() ([]common.Post, error) {
	rows, err := db.Connection.Query("SELECT title, excerpt, id FROM posts")
	if err != nil {
		return make([]common.Post, 0), err
	}
	defer rows.Close()

	all_posts := make([]common.Post, 0)
	for rows.Next() {
		var post common.Post
		if err = rows.Scan(&post.Title, &post.Excerpt, &post.Id); err != nil {
			return make([]common.Post, 0), err
		}
		all_posts = append(all_posts, post)
	}

	return all_posts, nil
}

// return post by id
func (db *Database) GetPost(post_id int) (common.Post, error) {
	rows, err := db.Connection.Query("SELECT title, content FROM posts WHERE id=?;", post_id)
	if err != nil {
		return common.Post{}, err
	}
	defer rows.Close()
	rows.Next()
	var post common.Post

	if err = rows.Scan(&post.Title, &post.Content); err != nil {
		return common.Post{}, err
	}

	return post, nil
}

// AddPost adds a post to the database
func (db *Database) AddPost(title string, excerpt string, content string) (int, error) {
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
func (db *Database) ChangePost(id int, title string, excerpt string, content string) error {
	tx, err := db.Connection.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if rbErr := tx.Rollback(); rbErr != nil && err == nil {
			err = fmt.Errorf("tx rollback error: %v", rbErr)
		}
	}()

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
func (db *Database) DeletePost(id int) error {
	if _, err := db.Connection.Exec("DELETE FROM posts WHERE id=?;", id); err != nil {
		return err
	}

	return nil
}

// AddImage will add the image metadata to the db.
// name - file name saved to disk
// alt - alternative text
// return(uuid, nil) if succeeded, ("", err) otherwise
func (db *Database) AddImage(uuid string, name string, alt string) error {
	tx, err := db.Connection.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if rbErr := tx.Rollback(); rbErr != nil && err == nil {
			err = fmt.Errorf("tx rollback error: %v", rbErr)
		}
	}()

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

func MakeSqlConnection() (Database, error) {
	/// Checking the DB connection

	connection_str, err := GetDatabaseURL()
	if err != nil {
		log.Error().Msgf("%s", err)
	}
	db, err := sql.Open("mysql", connection_str)
	if err != nil {
		return Database{}, err
	}
	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Second * 5)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return Database{
		MY_SQL_URL: connection_str,
		Connection: db,
	}, nil
}
