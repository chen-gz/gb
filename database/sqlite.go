package database

import (
	"database/sql"
	"math"
	"strconv"

	_ "errors"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"strings"
	"time"
)

const dbPath = "./blog.db"
const dbType = "sqlite3"
const max_post_num = 50 // max number of posts when summary is needed

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

//////////////////////////////////////////////
// database interface for v1 api

type BlogDataV1 struct {
	Id           int
	Author       string
	Title        string
	Content      string
	Tags         string
	Categories   string
	Url          string // for vue router and s3 storage. no space.
	Like         int
	Dislike      int
	CoverImg     string
	IsDraft      bool // if true, only show to author, else show to everyone
	IsDeleted    bool
	PrivateLevel int
	ViewCount    int
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Summary      string // summary of the post. This only used return the post to the front end
	Rendered     string // rendered html content, this only is used to return the post to the front end
}

const dbTypev1 = "sqlite3"
const dbPathv1 = "./blogv1.db"

func V1InitDatabase() {
	database, err := sql.Open(dbTypev1, dbPathv1)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()
	// create table if not exist
	_, err = database.Exec(`CREATE TABLE IF NOT EXISTS posts (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        author TEXT,
        title TEXT,
        content TEXT,
        tags TEXT,
        categories TEXT,
        url TEXT UNIQUE NOT NULL,
        like INTEGER,
        dislike INTEGER,
        cover_img TEXT,
        is_draft INTEGER,
        is_deleted INTEGER,
        private_level INTEGER,
        view_count INTEGER,
        created_at DATETIME,
        updated_at DATETIME
    )`)
	if err != nil {
		log.Fatal(err)
	}
}

// if success return blogdata, else return empty blogdata
func V1GetPostByUrl(url string) BlogDataV1 {
	log.Println("API V1: Get post by url: ", url)
	post := BlogDataV1{}
	if url == "" {
		log.Println("API V1: url is empty")
		return post
	}
	database, err := sql.Open(dbTypev1, dbPathv1)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()
	log.Print("sql command: SELECT id, author, title, content, tags, categories, url, like, dislike, cover_img, is_draft, " +
		"is_deleted, private_level, view_count, created_at, updated_at FROM posts WHERE url = ?")
	database.QueryRow(`SELECT id, author, title,
    content, tags, categories, url,
    like, dislike, cover_img, is_draft, is_deleted,
    private_level, view_count, created_at, updated_at
    FROM posts WHERE url = ?`, url).Scan(&post.Id, &post.Author, &post.Title, &post.Content, &post.Tags,
		&post.Categories, &post.Url, &post.Like, &post.Dislike, &post.CoverImg,
		&post.IsDraft, &post.IsDeleted,
		&post.PrivateLevel, &post.ViewCount,
		&post.CreatedAt, &post.UpdatedAt)
	return post
}

func V1SearchPost(keys map[string]string) []BlogDataV1 {
	log.Println("API V1: Search post by keys: ", keys)
	database, err := sql.Open(dbTypev1, dbPathv1)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()
	query_items := "SELECT id, author, title, url"
	query := "FROM posts "
	_ = query
	sort := ""
	if keys["sort"] != "" {
		sort = keys["sort"]
		// remove sort from keys
		delete(keys, "sort")
	}

	// if has summary, then limit to max 50 posts
	summary := false
	if keys["summary"] != "" {
		if keys["summary"] == "true" {
			summary = true
		}
		delete(keys, "summary")
	}
	limit := 50
	if keys["limit"] != "" {
		limit, err = strconv.Atoi(keys["limit"])
		if err != nil {
			limit = 50
		}
		delete(keys, "limit")
	}
	limit = int(math.Min(float64(limit), 50))

	if len(keys) != 0 {
		query += "WHERE "
		for key, value := range keys {
			if key == "tags" {
				query += key + " LIKE '%" + value + "%' AND "
			} else {
				query += key + "=" + value + " AND "

			}

		}
		if len(keys) != 0 {
			query = query[:len(query)-5]
		}
	}
	if sort != "" {
		query += " ORDER BY " + sort
	}

	if summary {
		query += " LIMIT " + strconv.Itoa(limit)
		query_items += " , content"
	}
	query = query_items + " " + query
	log.Println("API V1: Search post query: ", query)
	stmt, err := database.Prepare(query)
	result := []BlogDataV1{}
	if err != nil {
		log.Println("API V1: Search post prepare error: ", err)
		return result
	}
	rows, err := stmt.Query()
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		post := BlogDataV1{}
		if summary {
			err = rows.Scan(&post.Id, &post.Author, &post.Title, &post.Url, &post.Content)
		} else {
			err = rows.Scan(&post.Id, &post.Author, &post.Title, &post.Url)
		}
		if err != nil {
			log.Fatal(err)
		}
		result = append(result, post)
	}
	return result
}

func V1InsertPost(blog BlogDataV1) error {
	database, err := sql.Open(dbTypev1, dbPathv1)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()
	stmt, err := database.Prepare(`INSERT INTO posts (author, title, content, tags, categories, url,
    														like, dislike, cover_img, is_draft, is_deleted,
    														private_level, view_count, created_at, updated_at) VALUES
    														(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	_, err = stmt.Exec(blog.Author, blog.Title, blog.Content, blog.Tags, blog.Categories, blog.Url,
		blog.Like, blog.Dislike, blog.CoverImg, blog.IsDraft, blog.IsDeleted,
		blog.PrivateLevel, blog.ViewCount, blog.CreatedAt, blog.UpdatedAt)

	if err != nil {
		log.Fatal(err)
	}
	return err
}

type BlogTags struct {
	Name  string
	Count int
}

func V1GetTags() map[string]int {
	database, err := sql.Open(dbTypev1, dbPathv1)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()
	rows, err := database.Query("select distinct tags from posts")
	if err != nil {
		log.Fatal(err)
	}
	tags := []string{}
	for rows.Next() {
		var tag string
		err = rows.Scan(&tag)
		if err != nil {
			log.Fatal(err)
		}
		if tag != "" {
			// seperate by comma and remove space in front and back
			split_tags := strings.Split(tag, ",")
			for _, tag := range split_tags {
				tags = append(tags, strings.TrimSpace(tag))
			}
			// tags = append(tags, strings.Split(tag, ",")...)
		}
	}
	result := map[string]int{}
	// get the number of post for each tag
	for _, tag := range tags {
		log.Println("select count(*) from posts where tags like '%" + tag + "%'")
		rows := database.QueryRow("select count(*) from posts where tags like '%" + tag + "%'")
		if err != nil {
			log.Fatal(err)
		}
		var count int
		err = rows.Scan(&count)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(tag, count)
		result[tag] += count
	}
	defer rows.Close()
	return result

}

type BlogCategories struct {
	Name  string
	Count int
	Posts []BlogDataV1
}

func V1GetCategories() []BlogCategories {
	database, err := sql.Open(dbTypev1, dbPathv1)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()
	rows, err := database.Query("select distinct categories from posts")
	if err != nil {
		log.Fatal(err)
	}
	cates := []string{}
	for rows.Next() {
		var cate string
		err = rows.Scan(&cate)
		if err != nil {
			log.Fatal(err)
		}
		if cate != "" {
			cates = append(cates, strings.TrimSpace(cate))
		}
	}
	// list 5 post title for each category
	result := []BlogCategories{}
	for _, cate := range cates {
		rows := database.QueryRow("select count(*) from posts where categories='" + cate + "'")
		if err != nil {
			log.Fatal(err)
		}
		var count int
		err = rows.Scan(&count)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(cate, count)
		params := map[string]string{
			"categories": "'" + cate + "'",
			"limit":      "5",
			"sort":       "updated_at DESC",
			"summary":    "true",
		}
		posts := V1SearchPost(params)
		result = append(result, BlogCategories{Name: cate, Count: count, Posts: posts})
	}
	defer rows.Close()
	return result
}

func V1UpdatePost(blogData BlogDataV1) {
	database, err := sql.Open(dbTypev1, dbPathv1)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()
	stmt, err := database.Prepare(`UPDATE posts SET author=?, title=?, content=?, tags=?, categories=?, url=?,
                 															like=?, dislike=?, cover_img=?, is_draft=?, is_deleted=?,
                 															private_level=?, view_count=?, created_at=?, updated_at=? WHERE id=?`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(blogData.Author, blogData.Title, blogData.Content, blogData.Tags, blogData.Categories, blogData.Url,
		blogData.Like, blogData.Dislike, blogData.CoverImg, blogData.IsDraft, blogData.IsDeleted,
		blogData.PrivateLevel, blogData.ViewCount, blogData.CreatedAt, blogData.UpdatedAt, blogData.Id)
	if err != nil {
		log.Fatal(err)
	}
}
