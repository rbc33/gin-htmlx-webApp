package database

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/rbc33/gocms/common"
	"github.com/rs/zerolog/log"
)

type Database interface {
	GetPosts() ([]common.Post, error)
	GetPost(post_id int) (common.Post, error)
	AddPost(title string, excerpt string, content string) (int, error)
	ChangePost(id int, title string, excerpt string, content string) error
	DeletePost(id int) error
	AddImage(uuid string, name string, alt string) error
}

type SqlDatabase struct {
	MY_SQL_URL string
	Connection *sql.DB
}

// / GetPosts gets all the posts from the current
// / database connection.
func (db SqlDatabase) GetPosts() ([]common.Post, error) {
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
func (db *SqlDatabase) GetPost(post_id int) (common.Post, error) {
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
func (db *SqlDatabase) AddPost(title string, excerpt string, content string) (int, error) {
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
func (db *SqlDatabase) ChangePost(id int, title string, excerpt string, content string) error {
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

func MakeSqlConnection() (SqlDatabase, error) {
	/// Checking the DB connection
	err := godotenv.Load()
	if err != nil {
		return SqlDatabase{}, err
	}
	connection_str := os.Getenv("MY_SQL_URL")
	db, err := sql.Open("mysql", connection_str)
	if err != nil {
		return SqlDatabase{}, err
	}
	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Second * 5)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return SqlDatabase{
		MY_SQL_URL: connection_str,
		Connection: db,
	}, nil
}
