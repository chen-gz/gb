package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"strings"
	"time"
)

const dbPath = "./blog.db"
const dbType = "sqlite3"

type BlogData struct {
	Id         int
	Author     string
	Title      string
	Content    string
	Tags       []string
	Categories []string
	Datetime   time.Time
	Url        string // for vue router and s3 storage. no space.
}

func Init() {
	// create sqlite database
	database, err := sql.Open(dbType, dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	database.Exec(`CREATE TABLE IF NOT EXISTS posts 
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT, 
    author TEXT, 
    content TEXT, 
    tags TEXT, 
    categories TEXT, 
    datetime DATETIME,
    url TEXT UNIQUE NOT NULL
    `)
}

func InsertPost(blog BlogData) error {
	// log the title and author
	log.Println("insert post, title", blog.Title, "\t author", blog.Author)
	database, err := sql.Open(dbType, dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	stmt, err := database.Prepare(`INSERT INTO posts
    (title, author, content, tags, categories, datetime, url) 
    VALUES(?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		log.Fatal(err)
	}
	// connect Tags to a single string and surround with '[]'
	Tags := strings.Join(blog.Tags, ",")
	Categories := strings.Join(blog.Categories, ",")

	_, err = stmt.Exec(blog.Title, blog.Author,
		blog.Content, Tags, Categories,
		blog.Datetime, blog.Url)
	if err != nil {
		log.Fatal(err)
	}
	return err
}
func UpdatePost(blog BlogData) error {
	log.Println("update post, title", blog.Title, "\t author", blog.Author)
	database, err := sql.Open(dbType, dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()
	// update the post based on the id
	stmt, err := database.Prepare(`UPDATE posts SET 
                                    title = ?, content = ?, 
                                    tags = ?, categories = ?, 
                                    datetime = ?, author = ?, url = ?, 
                                    WHERE id = ?`)
	if err != nil {
		log.Fatal(err)
	}
	// connect Tags to a single string and connect use ',' to connect
	Tags := strings.Join(blog.Tags, ",")
	Categories := strings.Join(blog.Categories, "/")

	_, err = stmt.Exec(blog.Title, blog.Content,
		Tags, Categories,
		blog.Datetime, blog.Author, blog.Url,
		blog.Id)

	if err != nil {
		log.Fatal(err)
	}
	return err
}

// get all post id and title from database and return a map
func GetAllPostIdAndTitle() (map[int]string, error) {
	database, err := sql.Open(dbType, dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()
	rows, err := database.Query("SELECT id, title FROM posts")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	result := make(map[int]string)
	for rows.Next() {
		var id int
		var title string
		err := rows.Scan(&id, &title)
		if err != nil {
			log.Fatal(err)
		}
		result[id] = title
	}
	return result, err
}

func GetPostById(index int) (BlogData, error) {
	database, err := sql.Open(dbType, dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()
	query := database.QueryRow("SELECT * FROM posts WHERE Id = ?", index)
	post := BlogData{}
	tag := ""
	category := ""
	err = query.Scan(&post.Id, &post.Title, &post.Author,
		&post.Content, &tag, &category,
		&post.Datetime, &post.Url)
	post.Tags = strings.Split(tag, ",")
	post.Categories = strings.Split(category, ",")
	if err != nil {
		log.Println("error in get post by id")
		log.Fatal(err)
		// todo: id not found should be handled and test
	}
	return post, err
}

func GetRecentPosts(num int) ([]BlogData, error) {
	// sort by datetime and get the first `num` posts
	database, err := sql.Open(dbType, dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()
	rows, err := database.Query("SELECT * FROM posts ORDER BY datetime DESC LIMIT ?", num)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	result := []BlogData{}
	tag := ""
	category := ""
	for rows.Next() {
		post := BlogData{}
		err := rows.Scan(&post.Id, &post.Title, &post.Author,
			&post.Content, &tag, &category,
			&post.Datetime, &post.Url)
		post.Tags = strings.Split(tag, ",")
		post.Categories = strings.Split(category, ",")
		if err != nil {
			log.Fatal(err)
		}
		result = append(result, post)
	}

	return result, err
}
