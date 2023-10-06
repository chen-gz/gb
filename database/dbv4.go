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

type DatabaseConfig struct {
	MariadbAddress       string `json:"mariadb_address"`
	MariadbUser          string `json:"mariadb_user"`
	MariadbPassword      string `json:"mariadb_password"`
	MariadbBlogDatabase  string `json:"mariadb_blog_database"`
	MariadbBlogTable     string `json:"mariadb_blog_table"`
	MariadbBlogUserTable string `json:"mariadb_blog_user_table"`
}

// blogDbConfig is a global settings for this database
// It will be initialized in when call the function InitV4
var blogDbConfig DatabaseConfig

type V4BlogUserData struct {
	Id    int           `json:"id"`
	Email string        `json:"email"`
	Name  string        `json:"name"`
	Roles set.StringSet `json:"roles"`
}
type V4PostData struct {
	Id              int           `json:"id"`
	Title           string        `json:"title"`
	Author          string        `json:"author"`
	AuthorEmail     string        `json:"author_email"`
	Url             string        `json:"url"`
	IsDraft         bool          `json:"is_draft"`
	IsDeleted       bool          `json:"is_deleted"`
	Content         string        `json:"content"`
	ContentRendered string        `json:"content_rendered"`
	Summary         string        `json:"summary"`
	Tags            string        `json:"tags"`
	Category        string        `json:"category"`
	CoverImage      string        `json:"cover_image"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
	ViewGroups      set.StringSet `json:"view_groups"`
	EditGroups      set.StringSet `json:"edit_groups"`
}

func initializeV4Table(db_blog *sql.DB) {
	// test post table exist in database or not
	var value int
	query := fmt.Sprintf(`SELECT 1 from information_schema.TABLES where TABLE_NAME='%s' and TABLE_SCHEMA='%s'`,
		blogDbConfig.MariadbBlogTable, blogDbConfig.MariadbBlogDatabase)
	err := db_blog.QueryRow(query).Scan(&value)
	if err == nil && value == 1 {
		// table exist. Do nothing
		return
	}
	// create table post
	// query = fmt.Sprintf(

	query = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
						id INT UNSIGNED AUTO_INCREMENT,
						title VARCHAR(255) NOT NULL,
						author VARCHAR(255) default '',
						author_email VARCHAR(255) default '',
						url VARCHAR(255) UNIQUE  NOT NULL,
						is_draft BOOLEAN DEFAULT FALSE,
						is_deleted BOOLEAN DEFAULT FALSE,
						content TEXT default '',
						content_rendered TEXT default '',  -- this field should be markdown, html, json, latex, etc.
						summary TEXT default '',
						tags VARCHAR(255) default '',
						category VARCHAR(255) default '',
						cover_image VARCHAR(255) default '',
						created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
						updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
						view_groups SET('admin', 'editor', 'author', 'premium', 'subscriber', 'guest') NOT NULL DEFAULT 'admin,editor,author,premium,subscriber,guest',
						edit_groups SET('admin', 'editor', 'author', 'premium', 'subscriber', 'guest') NOT NULL DEFAULT 'admin,editor,author',
						PRIMARY KEY (id))`, blogDbConfig.MariadbBlogTable)

	_, err = db_blog.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
	// if the index not exist, create it
	query = fmt.Sprintf(`ALTER TABLE %s ADD FULLTEXT INDEX idx_content (content)`, blogDbConfig.MariadbBlogTable)
	// _, err = db_blog.Exec(`ALTER TABLE post ADD FULLTEXT INDEX idx_content (content)`)
	_, err = db_blog.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
	query = fmt.Sprintf(`ALTER TABLE %s ADD FULLTEXT INDEX idx_tags (tags)`, blogDbConfig.MariadbBlogTable)
	_, err = db_blog.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
	query = fmt.Sprintf(`ALTER TABLE %s ADD FULLTEXT INDEX idx_category (category)`, blogDbConfig.MariadbBlogTable)
	_, err = db_blog.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	query = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
    					id INT UNSIGNED AUTO_INCREMENT,
    					email VARCHAR(255) UNIQUE NOT NULL,
    					name VARCHAR(255),
    					roles SET('admin', 'editor', 'author', 'premium', 'subscriber', 'guest') DEFAULT 'guest',
    					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
							PRIMARY KEY (id))`, blogDbConfig.MariadbBlogUserTable)
	_, err = db_blog.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
	return
}

// InitV4 init database for v4
// Create two tables: post and blog_users
// The content of post refer to the structure of V4PostData
// The content of blog_users refer to the structure of V4BlogUserData
// return db_blog
func InitV4(config DatabaseConfig) (db_blog *sql.DB) {
	blogDbConfig = config
	// connect to database
	sql_endpoint := config.MariadbUser + ":" + config.MariadbPassword + config.MariadbAddress
	db, err := sql.Open("mysql", sql_endpoint)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + config.MariadbBlogDatabase)
	if err != nil {
		log.Fatal(err)
	}
	db.Close()
	// connect to database eta_blog
	sql_endpoint = config.MariadbUser + ":" + config.MariadbPassword + config.MariadbAddress + "/" + config.MariadbBlogDatabase
	db_blog, err = sql.Open("mysql", sql_endpoint)
	initializeV4Table(db_blog)

	if err != nil {
		log.Fatal(err)
	}
	return db_blog
}

// v4InsertPostMigrate insert post data into database
// This function only used for migration from v3 to v4
// The create and update time will be set to the time from v3
func v4InsertPostMigrate(db_blog *sql.DB, post V4PostData) error {
	query := fmt.Sprintf(`INSERT INTO %s (title, author, author_email,url, is_draft, is_deleted, content,
                  summary, tags, category, cover_image, created_at, updated_at,
                  view_groups, edit_groups) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`, blogDbConfig.MariadbBlogTable)
	stmt, _ := db_blog.Prepare(query)
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
	viewGroups := strings.Join(post.ViewGroups.ToSlice(), ",")
	editGroups := strings.Join(post.EditGroups.ToSlice(), ",")
	_, err := stmt.Exec(post.Title, post.Author, post.AuthorEmail,
		post.Url, post.IsDraft, post.IsDeleted, post.Content,
		post.Summary, post.Tags, post.Category, post.CoverImage,
		created_at, updated_at, viewGroups, editGroups)
	if err != nil {
		log.Println("v4InsertPost error: ", err)
		return err
	}
	return nil
}

// v4InsertPost insert post to database
func v4InsertPost(db_blog *sql.DB, post V4PostData) error {
	query := fmt.Sprintf(`INSERT INTO %s (title, author, author_email, url, is_draft, is_deleted, content,
                  summary, tags, category, cover_image, view_groups, edit_groups)
                  VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)`, blogDbConfig.MariadbBlogTable)
	stmt, _ := db_blog.Prepare(query)
	var created_at, updated_at interface{}
	viewGroups := strings.Join(post.ViewGroups.ToSlice(), ",")
	editGroups := strings.Join(post.EditGroups.ToSlice(), ",")
	_, err := stmt.Exec(post.Title, post.Author, post.AuthorEmail,
		post.Url, post.IsDraft, post.IsDeleted, post.Content,
		post.Summary, post.Tags, post.Category, post.CoverImage,
		created_at, updated_at, viewGroups, editGroups)
	if err != nil {
		log.Println("v4InsertPost error: ", err)
		return err
	}
	return nil
}

func getPostByUrl(db *sql.DB, url string) (V4PostData, error) {
	var post V4PostData
	var created_at_str, updated_at_str, view_groups_str, edit_groups_str string

	query := fmt.Sprintf(`SELECT id, title, author, author_email, url, is_draft, is_deleted, content, content_rendered,
	  		   summary, tags, category, cover_image, created_at, updated_at, view_groups, edit_groups
				FROM %s WHERE url=?`, blogDbConfig.MariadbBlogTable)
	err := db.QueryRow(query, url).Scan(&post.Id, &post.Title, &post.Author, &post.AuthorEmail, &post.Url, &post.IsDraft,
		&post.IsDeleted, &post.Content, &post.ContentRendered, &post.Summary, &post.Tags, &post.Category, &post.CoverImage,
		&created_at_str, &updated_at_str, &view_groups_str, &edit_groups_str)
	post.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", created_at_str) // todo: error may need to be handled
	post.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updated_at_str) // todo: error may need to be handled
	post.ViewGroups = set.CreateStringSet(strings.Split(view_groups_str, ",")...)
	post.EditGroups = set.CreateStringSet(strings.Split(edit_groups_str, ",")...)
	if err != nil {
		log.Println("getPostByUrl err: ", err)
		return post, err
	}
	return post, nil
}

func getPostById(db *sql.DB, id int) (V4PostData, error) {
	var post V4PostData
	var created_at_str, updated_at_str, view_groups_str, edit_groups_str string
	query := fmt.Sprintf(`SELECT * FROM %s WHERE id=?`, blogDbConfig.MariadbBlogTable)
	err := db.QueryRow(query, id).Scan(&post.Id, &post.Title, &post.Author, &post.AuthorEmail, &post.Url,
		&post.IsDraft, &post.IsDeleted, &post.Content, &post.ContentRendered, &post.Summary, &post.Tags, &post.Category, &post.CoverImage,
		&created_at_str, &updated_at_str, &view_groups_str, &edit_groups_str)
	post.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", created_at_str) // todo: error may need to be handled
	post.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updated_at_str) // todo: error may need to be handled
	post.ViewGroups = set.CreateStringSet(strings.Split(view_groups_str, ",")...)
	post.EditGroups = set.CreateStringSet(strings.Split(edit_groups_str, ",")...)
	if err != nil {
		log.Println("getPostById error: ", err)
		return V4PostData{}, err
	}
	return post, nil
}

func updatePost(db *sql.DB, post V4PostData) error {
	query := fmt.Sprintf(`UPDATE %s SET title=?, author=?, author_email=?, url=?, is_draft=?, is_deleted=?, content=?,
				  summary=?, tags=?, category=?, cover_image=?, view_groups=?, edit_groups=? WHERE id=?`, blogDbConfig.MariadbBlogTable)
	stmt, _ := db.Prepare(query)
	viewGroups := strings.Join(post.ViewGroups.ToSlice(), ",")
	editGroups := strings.Join(post.EditGroups.ToSlice(), ",")
	_, err := stmt.Exec(post.Title, post.Author, post.AuthorEmail, post.Url, post.IsDraft, post.IsDeleted, post.Content,
		post.Summary, post.Tags, post.Category, post.CoverImage, viewGroups, editGroups, post.Id)
	if err != nil {
		log.Println("updatePost error: ", err, "post: ", post)
		return err
	}
	return nil
}

func v4InsertUser(db *sql.DB, user V4BlogUserData) error {
	query := fmt.Sprintf(`INSERT INTO %s (email, name, roles) VALUES (?,?,?)`, blogDbConfig.MariadbBlogUserTable)
	stmt, _ := db.Prepare(query)
	roles := strings.Join(user.Roles.ToSlice(), ",")
	_, err := stmt.Exec(user.Email, user.Name, roles)
	if err != nil {
		log.Println("v4InsertUser error: ", err)
		return err
	}
	return nil
}

func getUserRole(db *sql.DB, user User) (set.StringSet, error) {
	var roles set.StringSet
	var roles_str string
	err := db.QueryRow(`SELECT roles FROM blog_users WHERE email=?`, user.Email).Scan(&roles_str)
	if err != nil {
		log.Println("getUserRole error: ", err)
		// default role is guest
		roles = set.CreateStringSet("guest")
		if len(user.Email) > 0 && len(user.Name) > 0 && user.Id != 0 {
			v4InsertUser(db, V4BlogUserData{
				Email: user.Email,
				Name:  user.Name,
				Roles: set.CreateStringSet("guest"),
			})
			return roles, nil
		} else {
			roles = set.CreateStringSet("guest")
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
		if len(roles.ToSlice()) == 1 {
			stmt += `)`
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
		params.Limit["size"] = 100
	}

	fmt.Println(params.Limit["start"], params.Limit["size"])
	stmt += `LIMIT ` + fmt.Sprintf("%d", params.Limit["start"]) + `,` + fmt.Sprintf("%d", params.Limit["size"])
	log.Println("searchPosts: ", stmt)
	// execute sql
	rows, err := db.Query(stmt)
	if err != nil {
		log.Println("searchPosts error: ", err)
		return []V4PostData{}, err
	}
	defer rows.Close()
	var posts []V4PostData
	for rows.Next() {
		var post V4PostData
		var created_at_str, updated_at_str string
		var view_groups_str, edit_groups_str string
		err := rows.Scan(&post.Id, &post.Title, &post.Author, &post.AuthorEmail, &post.Url,
			&post.IsDraft, &post.IsDeleted, &post.Content, &post.ContentRendered, &post.Summary, &post.Tags, &post.Category, &post.CoverImage,
			&created_at_str, &updated_at_str, &view_groups_str, &edit_groups_str)
		post.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", created_at_str)
		post.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updated_at_str)
		post.ViewGroups = set.CreateStringSet(strings.Split(view_groups_str, ",")...)
		post.EditGroups = set.CreateStringSet(strings.Split(edit_groups_str, ",")...)
		if err != nil {
			log.Println("searchPosts error: ", err)
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
	old_post, err := getPostById(db, post.Id)
	if err != nil {
		log.Println("V4UpdatePosByUser: ", err)
		return err
	}
	if roles.Contains("admin") || roles.Contains("editor") || (old_post.AuthorEmail == user.Email && user.Email != "") {
		log.Println("V4UpdatePosByUser: ", roles)
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
			Title:       "New Post",
			Url:         time.Now().Format("2006-01-02 15:04:05") + "-new-post",
			ViewGroups:  set.CreateStringSet("admin,editor,author,premium,subscriber,guest"),
			EditGroups:  set.CreateStringSet("admin,editor,author"),
			Author:      user.Name,
			AuthorEmail: user.Email,
			IsDraft:     false,
			CreatedAt:   time.Now(),
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
