CREATE DATABASE gocms;
use gocms;
CREATE TABLE posts ( 
    id INT AUTO_INCREMENT PRIMARY KEY,
    content TEXT
);
INSERT INTO posts(content) VALUES("This is an example of a post");
INSERT INTO posts(content) VALUES("This is yet another example of a post");
ALTER TABLE posts ADD title TEXT;
UPDATE posts SET title="Post 1" WHERE id=1;
UPDATE posts SET title="Post 2" WHERE id=2;
INSERT INTO posts(title, content) VALUES("Incredible post","## Markdown\n content");
