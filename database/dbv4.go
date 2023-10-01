package database

import (
	"database/sql"
	"errors"
	"fmt"
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

// InitV4 init database for v4
// Create two tables: post and blog_users
// The content of post refer to the structure of V4PostData
// The content of blog_users refer to the structure of V4BlogUserData
// return db_blog
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
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
					view_groups SET('admin', 'editor', 'author', 'premium', 'subscriber', 'guest') DEFAULT 'admin,editor,author,premium,subscriber,guest',
					edit_groups SET('admin', 'editor', 'author', 'premium', 'subscriber', 'guest') DEFAULT 'admin,editor,author',
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

// v4InsertPost insert post to database
func v4InsertPost(db_blog *sql.DB, post V4PostData) error {
	stmt, _ := db_blog.Prepare(`INSERT INTO post (title, author, author_email,url, is_draft, is_deleted, content, 
                  summary, tags, category, cover_image, created_at, updated_at,
                  view_groups, edit_groups) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`)
	var created_at, updated_at interface{}
	if post.CreatedAt.IsZero() {
		created_at = nil
	} else {
		created_at = post.CreatedAt
	}
	if post.UpdatedAt.IsZero() {
		updated_at = nil
	} else {
		updated_at = post.UpdatedAt
	}

	// convert viewGroup as a string. e.g. {"admin", "editor"} -> "admin,editor"
	viewGroups := strings.Join(post.ViewGroups.ToSlice(), ",")
	editGroups := strings.Join(post.EditGroups.ToSlice(), ",")
	_, err := stmt.Exec(post.Title, post.Author, post.AuthorEmail,
		post.Url, post.IsDraft, post.IsDeleted, post.Content,
		post.Summary, post.Tags, post.Category, post.CoverImage,
		created_at, updated_at,
		viewGroups, editGroups)
	if err != nil {
		log.Fatal("v4InsertPost error: ", err)
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
		log.Fatal("getPostByUrl error: ", err)
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
		log.Fatal("getPostById error: ", err)
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
		log.Fatal("updatePost error: ", err)
		return err
	}
	return nil
}
func v4InsertUser(db *sql.DB, user V4BlogUserData) error {
	stmt, _ := db.Prepare(`INSERT INTO blog_users (email, name, roles) VALUES (?,?,?)`)
	roles := strings.Join(user.Roles.ToSlice(), ",")
	_, err := stmt.Exec(user.Email, user.Name, roles)
	if err != nil {
		log.Fatal("v4InsertUser error: ", err)
		return err
	}
	return nil
}
func getUserRole(db *sql.DB, user User) (set.StringSet, error) {
	var roles set.StringSet
	var roles_str string
	err := db.QueryRow(`SELECT roles FROM blog_users WHERE email=?`, user.Email).Scan(&roles_str)
	if err != nil {
		// default role is guest
		roles = set.CreateStringSet("guest")

		if len(user.Email) > 0 && len(user.Name) > 0 && user.Id != 0 {
			v4InsertUser(db, V4BlogUserData{
				Email: user.Email,
				Name:  user.Name,
				Roles: set.CreateStringSet("guest"),
			})
		}
		return roles, err
	}
	roles = set.CreateStringSet(strings.Split(roles_str, ",")...)
	return roles, nil
}

type SearchParams struct {
	Author     string         `json:"author"`      // exact match
	Title      string         `json:"title"`       // use like to search
	Limit      map[string]int `json:"limit"`       // two values: start, size the number of post to return
	Sort       string         `json:"sort"`        // directly apply to sql
	Rendered   bool           `json:"rendered"`    // if true, rendered content will be returned, default false;
	CountsOnly bool           `json:"counts_only"` // if true, only return the count of the result, default false;
	Content    string         `json:"content"`     // use match to search
	Tags       string         `json:"tags"`        // use match to search
	Categories string         `json:"categories"`  // use match to search
	IsDraft    bool           `json:"is_draft"`    // if true, only return the draft post, default false;
	IsDeleted  bool           `json:"is_deleted"`  // if true, only return the deleted post, default false;
}

func searchPosts(db *sql.DB, params SearchParams, user User) ([]V4PostData, error) {
	roles, _ := getUserRole(db, user)
	stmt := `SELECT * from post WHERE `
	for index, item := range roles.ToSlice() {
		if index == 0 {
			stmt += `(FIND_IN_SET("` + item + `", view_groups) `
		} else if index < len(roles.ToSlice())-1 {
			stmt += `OR FIND_IN_SET("` + item + `", view_groups) `
		} else if index == len(roles.ToSlice())-1 {
			stmt += `OR FIND_IN_SET("` + item + `", view_groups)) `
		}
	}
	if params.Author != "" {
		stmt += `AND author="` + params.Author + `"`
	}
	if params.Title != "" {
		stmt += `AND title LIKE "%` + params.Title + `%" `
	}
	if params.Content != "" {
		stmt += `AND MATCH (content) AGAINST ("` + params.Content + `") `
	}
	if params.Tags != "" {
		stmt += `AND MATCH (tags) AGAINST ("` + params.Tags + `") `
	}
	if params.Categories != "" {
		stmt += `AND MATCH (category) AGAINST ("` + params.Categories + `") `
	}
	// add limit and sort
	if params.Sort != "" {
		stmt += `ORDER BY ` + params.Sort + ` `
	}
	// limit should have two values: start, size. if not, use default value
	if params.Limit == nil {
		params.Limit = make(map[string]int)
	}
	if params.Limit["start"] == 0 {
		params.Limit["start"] = 0
	}
	if params.Limit["size"] == 0 {
		params.Limit["size"] = 10
	}

	fmt.Println(params.Limit["start"], params.Limit["size"])
	stmt += `LIMIT ` + fmt.Sprintf("%d", params.Limit["start"]) + `,` + fmt.Sprintf("%d", params.Limit["size"])
	// execute sql
	log.Println(stmt)
	rows, err := db.Query(stmt)
	if err != nil {
		log.Fatal("searchPosts error: ", err)
		return []V4PostData{}, err
	}
	defer rows.Close()
	var posts []V4PostData
	for rows.Next() {
		var post V4PostData
		var created_at_str, updated_at_str string
		var view_groups_str, edit_groups_str string
		err := rows.Scan(&post.Id, &post.Title, &post.Author, &post.AuthorEmail, &post.Url,
			&post.IsDraft, &post.IsDeleted, &post.Content, &post.Summary, &post.Tags, &post.Category, &post.CoverImage,
			&created_at_str, &updated_at_str, &view_groups_str, &edit_groups_str)
		post.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", created_at_str)
		post.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updated_at_str)
		post.ViewGroups = set.CreateStringSet(strings.Split(view_groups_str, ",")...)
		post.EditGroups = set.CreateStringSet(strings.Split(edit_groups_str, ",")...)
		if err != nil {
			log.Fatal("searchPosts error: ", err)
			log.Println("searchPosts error")
			return []V4PostData{}, err
		}
		posts = append(posts, post)
	}
	return posts, nil

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
func V4SearchPostUser(db *sql.DB, params SearchParams, user User) ([]V4PostData, error) {
	return searchPosts(db, params, user)
}

func V4NewPostUser(db *sql.DB, user User) (string, error) {
	roles, _ := getUserRole(db, user)
	if roles.Contains("admin") || roles.Contains("editor") || roles.Contains("author") {
		post := V4PostData{
			Title:     "New Post",
			Url:       time.Now().Format("2006-01-02") + "-new-post",
			CreatedAt: time.Now(),
		}
		v4InsertPost(db, post)
		return post.Url, nil
	} else {
		return "", errors.New("permission denied")
	}
}
func V4GetDistinctUser(db *sql.DB, field string, user User) ([]string, error) {
	//roles, _ := getUserRole(db, user)
	var result []string
	// field can only be author, tags, category, title, url
	if field != "author" && field != "tags" && field != "category" && field != "title" && field != "url" {
		return result, errors.New("invalid field")
	}
	rows, err := db.Query(`SELECT DISTINCT ` + field + ` FROM post`)
	if err != nil {
		log.Fatal("V4GetDistinctUser error: ", err)
		return result, err
	}
	defer rows.Close()
	for rows.Next() {
		var field string
		err := rows.Scan(&field)
		if err != nil {
			log.Fatal("V4GetDistinctUser error: ", err)
			return result, err
		}
		result = append(result, field)
	}
	return result, nil
}
