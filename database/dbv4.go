package database

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/minio/minio-go/v7/pkg/set"
	"log"
	"strings"
	"time"
)

// the data structure still use from v3
type V4BlogUserData struct {
	Id    int           `json:"id"`
	Email string        `json:"email"`
	Name  string        `json:"name"`
	Roles set.StringSet `json:"roles"`
}

func InitV4() (db_blog *sql.DB) {
	// connect to database
	db, err := sql.Open("mysql", "zong:Connie@tcp(192.168.0.174:3306)/")
	if err != nil {
		panic(err)
	}
	// check the database eta_blog is exist or not
	// if not exist, create a new database
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS eta_blog")
	if err != nil {
		panic(err)
	}
	db.Close()
	// connect to database eta_blog
	db_blog, err = sql.Open("mysql", "zong:Connie@tcp(192.168.0.174:3306)/eta_blog")
	if err != nil {
		panic(err)
	}
	_, err = db_blog.Exec(`CREATE TABLE IF NOT EXISTS post (
					id INT UNSIGNED AUTO_INCREMENT,
					title VARCHAR(255) NOT NULL,
					author VARCHAR(255) default '',
					author_email VARCHAR(255) default '',
					url VARCHAR(255) UNIQUE  NOT NULL,
					is_draft BOOLEAN DEFAULT FALSE,
					is_deleted BOOLEAN DEFAULT FALSE,
					content TEXT default '',
					summary TEXT default '',
					tags VARCHAR(255) default '',
					category VARCHAR(255) default '',
					cover_image VARCHAR(255) default '',
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
					view_groups SET('admin', 'editor', 'author', 'premium', 'subscriber', 'guest') DEFAULT 'guest',
					edit_groups SET('admin', 'editor', 'author', 'premium', 'subscriber', 'guest') DEFAULT 'author',
					PRIMARY KEY (id))`)
	if err != nil {
		panic(err)
	}
	_, err = db_blog.Exec(`CREATE TABLE IF NOT EXISTS blog_users (
    					id INT UNSIGNED AUTO_INCREMENT,
    					email VARCHAR(255) UNIQUE NOT NULL,
    					name VARCHAR(255),
    					roles SET('admin', 'editor', 'author', 'premium', 'subscriber', 'guest') DEFAULT 'guest',
    					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
						PRIMARY KEY (id))`)
	if err != nil {
		panic(err)
	}

	return db_blog
}

type V4PostData struct {
	Id          int           `json:"id"`
	Title       string        `json:"title"`
	Author      string        `json:"author"`
	AuthorEmail string        `json:"author_email"`
	Url         string        `json:"url"`
	IsDraft     bool          `json:"is_draft"`
	IsDeleted   bool          `json:"is_deleted"`
	Content     string        `json:"content"`
	Summary     string        `json:"summary"`
	Tags        string        `json:"tags"`
	Category    string        `json:"category"`
	CoverImage  string        `json:"cover_image"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	ViewGroups  set.StringSet `json:"view_groups"`
	EditGroups  set.StringSet `json:"edit_groups"`
}

func v4InsertPost(db_blog *sql.DB, post V4PostData) error {
	stmt, _ := db_blog.Prepare(`INSERT INTO post (title, author, author_email,url, is_draft, is_deleted, content, 
                  summary, tags, category, cover_image, view_groups, edit_groups) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)`)
	// convert viewGroup as a string. e.g. {"admin", "editor"} -> "admin,editor"
	viewGroups := strings.Join(post.ViewGroups.ToSlice(), ",")
	editGroups := strings.Join(post.EditGroups.ToSlice(), ",")
	_, err := stmt.Exec(post.Title, post.Author, post.AuthorEmail,
		post.Url, post.IsDraft, post.IsDeleted, post.Content,
		post.Summary, post.Tags, post.Category, post.CoverImage, viewGroups, editGroups)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func getPostByUrl(db *sql.DB, url string) (V4PostData, error) {
	var post V4PostData
	var created_at_str, updated_at_str, view_groups_str, edit_groups_str string

	err := db.QueryRow(`SELECT id, title, author, author_email, url, is_draft, is_deleted, content,
	  		   summary, tags, category, cover_image, created_at, updated_at, view_groups, edit_groups
				FROM post WHERE url=?`, url).Scan(&post.Id, &post.Title, &post.Author, &post.AuthorEmail, &post.Url,
		&post.IsDraft, &post.IsDeleted, &post.Content, &post.Summary, &post.Tags, &post.Category, &post.CoverImage,
		&created_at_str, &updated_at_str, &view_groups_str, &edit_groups_str)
	post.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", created_at_str)
	post.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updated_at_str)
	post.ViewGroups = set.CreateStringSet(strings.Split(view_groups_str, ",")...)
	post.EditGroups = set.CreateStringSet(strings.Split(edit_groups_str, ",")...)

	if err != nil {
		log.Println("getPostByUrl error")
		log.Fatal(err)
		return post, err
	}
	log.Println(post.Url)
	return post, nil
}
func getPostById(db *sql.DB, id int) (V4PostData, error) {
	var post V4PostData
	err := db.QueryRow(`SELECT * FROM post WHERE id=?`, id).Scan(&post.Id, &post.Title, &post.Author, &post.AuthorEmail, &post.Url,
		&post.IsDraft, &post.IsDeleted, &post.Content, &post.Summary, &post.Tags, &post.Category, &post.CoverImage,
		&post.CreatedAt, &post.UpdatedAt, &post.ViewGroups, &post.EditGroups)
	if err != nil {
		log.Fatal(err)
		return post, err
	}
	return post, nil
}
func updatePost(db *sql.DB, post V4PostData) error {
	stmt, _ := db.Prepare(`UPDATE post SET title=?, author=?,author_email=?, url=?, is_draft=?, is_deleted=?, content=?, 
				  summary=?, tags=?, category=?, cover_image=?, view_groups=?, edit_groups=? WHERE id=?`)
	_, err := stmt.Exec(post.Title, post.Author, post.AuthorEmail, post.Url, post.IsDraft, post.IsDeleted, post.Content,
		post.Summary, post.Tags, post.Category, post.CoverImage, post.ViewGroups, post.EditGroups, post.Id)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
func v4InsertUser(db *sql.DB, user V4BlogUserData) error {
	stmt, _ := db.Prepare(`INSERT INTO blog_users (email, name, roles) VALUES (?,?,?)`)
	roles := strings.Join(user.Roles.ToSlice(), ",")
	_, err := stmt.Exec(user.Email, user.Name, roles)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
func getUserRole(db *sql.DB, user User) (set.StringSet, error) {
	var roles set.StringSet
	err := db.QueryRow(`SELECT roles FROM blog_users WHERE email=?`, user.Email).Scan(&roles)
	if err != nil {
		return set.StringSet{}, err
	}
	return roles, nil
}

// following is public function
func V4InsertPosByUser(db *sql.DB, post V4PostData, user User) error {
	// get user role
	roles, _ := getUserRole(db, user)
	if roles.Contains("admin") || roles.Contains("editor") || roles.Contains("author") {
		// if user is admin, editor or author, insert post
		v4InsertPost(db, post)
	} else {
		return errors.New("permission denied")
	}
	return nil
}
func V4UpdatePosByUser(db *sql.DB, post V4PostData, user User) error {
	roles, _ := getUserRole(db, user)
	// get old_post
	old_post, _ := getPostById(db, post.Id)

	if roles.Contains("admin") || roles.Contains("editor") || old_post.AuthorEmail == user.Email {
		updatePost(db, post)
	} else {
		return errors.New("permission denied")
	}
	return nil
}
func V4GetPostByUrlUser(db *sql.DB, url string, user User) (V4PostData, error) {
	roles, _ := getUserRole(db, user)
	// get post
	post, _ := getPostByUrl(db, url)
	if post.ViewGroups.Contains("guest") {
		return post, nil
	}
	if roles.Contains("admin") || roles.Contains("editor") ||
		(len(post.AuthorEmail) > 0 && post.AuthorEmail == user.Email) {
		return post, nil
	}
	return V4PostData{}, errors.New("permission denied")
}
