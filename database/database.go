package database

import (
	"database/sql"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/rbc33/common"
)

type Database struct {
	MY_SQL_URL string
	Connection *sql.DB
}

// / This function gets all the posts from the current
// / database connection.
func (db Database) GetPosts() ([]common.Post, error) {
	rows, err := db.Connection.Query("SELECT title, content FROM posts")
	if err != nil {
		return make([]common.Post, 0), err
	}

	all_posts := make([]common.Post, 0)
	for rows.Next() {
		var post common.Post
		if err = rows.Scan(&post.Title, &post.Content); err != nil {
			return make([]common.Post, 0), err
		}
		all_posts = append(all_posts, post)
	}

	return all_posts, nil
}

// return post by id
func (db Database) GetPost(post_id int) (common.Post, error) {
	rows, err := db.Connection.Query("SELECT title, content FROM posts WHERE id=?;", post_id)
	if err != nil {
		return common.Post{}, err
	}
	rows.Next()
	var post common.Post

	if err = rows.Scan(&post.Title, &post.Content); err != nil {
		return common.Post{}, err
	}

	return post, nil
}

func MakeSqlConnection() (Database, error) {
	/// Checking the DB connection
	err := godotenv.Load()
	if err != nil {
		return Database{}, err
	}
	connection_str := os.Getenv("MY_SQL_URL")
	db, err := sql.Open("mysql", connection_str)
	if err != nil {
		return Database{}, err
	}
	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return Database{
		MY_SQL_URL: connection_str,
		Connection: db,
	}, nil
}
