package database

import (
	"database/sql"
	"errors"
	"log"
	"strings"
	"time"
)

// initialize v2 database
const dbPathV2 = "v2.db"
const dbTypeV2 = "sqlite3"

func InitV2() {
	db, _ := sql.Open(dbTypeV2, dbPathV2)
	defer db.Close()
	// create post table
	stmt, _ := db.Prepare(`CREATE TABLE IF NOT EXISTS post (
    							 id INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
    							 title text,
    							 author text,
    							 url text UNIQUE NOT NULL,
								 create_time DateTime,
								 update_time DateTime,
								 private_level integer,
								 summary text,
								 visible_groups text,
								 is_draft integer,
								 is_deleted integer
								 );`)
	// create post_content table
	_, err := stmt.Exec()
	if err != nil {
		return
	}

	stmt, _ = db.Prepare(`CREATE VIRTUAL TABLE post_content USING fts5(id, content, category, tags);`)
	_, err = stmt.Exec()
	if err != nil {
		log.Println(err)
	}
	stmt, _ = db.Prepare(`CREATE TABLE IF NOT EXISTS post_comment (
    							id INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
    							likes integer,
    							dislikes integer,
    							view_count integer,
    							comments text);`)
	stmt.Exec()
}

type PostDataV2Meta struct {
	Id            int       `json:"id"`
	Title         string    `json:"title"`
	Author        string    `json:"author"`
	Url           string    `json:"url"`
	CreateTime    time.Time `json:"create_time"`
	UpdateTime    time.Time `json:"update_time"`
	PrivateLevel  int       `json:"private_level"`
	Summary       string    `json:"summary"`
	VisibleGroups string    `json:"visible_groups"`
	IsDraft       bool      `json:"is_draft"`
	IsDeleted     bool      `json:"is_deleted"`
}
type PostDataV2Content struct {
	Id       int    `json:"id"`
	Content  string `json:"content"`
	Category string `json:"category"`
	Tags     string `json:"tags"`
}
type PostDataV2Comment struct {
	Id        int    `json:"id"`
	Likes     int    `json:"likes"`
	Dislikes  int    `json:"dislikes"`
	ViewCount int    `json:"view_count"`
	Comments  string `json:"comments"`
}
type PostDataV2 struct {
	Meta    PostDataV2Meta    `json:"meta"`
	Content PostDataV2Content `json:"content"`
	Comment PostDataV2Comment `json:"comment"`
}

func V2InsertPost(postData PostDataV2) error {
	meta := postData.Meta
	content := postData.Content
	comment := postData.Comment
	// the id of content and comment should be the same as post.
	db, _ := sql.Open(dbTypeV2, dbPathV2)
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Println(err)
		}
	}(db)
	stmt, _ := db.Prepare(`INSERT INTO post (title, author, url,  create_time, update_time,
                  					private_level, summary, visible_groups, is_draft, is_deleted)
    			  					VALUES (?,?,?,?,?,?,?,?,?,?)`)
	_, err := stmt.Exec(meta.Title, meta.Author, meta.Url, meta.CreateTime, meta.UpdateTime,
		meta.PrivateLevel, meta.Summary, meta.VisibleGroups, meta.IsDraft, meta.IsDeleted)
	if err != nil {
		log.Println(err)
		return err
	}
	// get the id of the post
	rows, err := db.Query(`SELECT id FROM post WHERE url = ?`, meta.Url)
	if err != nil {
		log.Println(err) // this should not happen
		return err
	}
	for rows.Next() {
		err := rows.Scan(&meta.Id)
		if err != nil {
			log.Println(err) // this should not happen
			return err
		}
	}
	content.Id = meta.Id
	comment.Id = meta.Id
	stmt, err = db.Prepare(`INSERT INTO post_content (id, content, category, tags) VALUES (?,?,?,?)`)
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = stmt.Exec(content.Id, content.Content, content.Category, content.Tags)
	if err != nil {
		log.Println(err)
		return err
	}
	stmt, err = db.Prepare(`INSERT INTO post_comment (id, likes, dislikes, view_count, comments) VALUES (?,?,?,?,?)`)
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = stmt.Exec(comment.Id, comment.Likes, comment.Dislikes, comment.ViewCount, comment.Comments)
	if err != nil {
		log.Println(err)
	}
	return nil
}

func V2GetPostByUrl(url string) PostDataV2 {
	db, _ := sql.Open(dbTypeV2, dbPathV2)
	defer db.Close()
	row := db.QueryRow(`SELECT * FROM post WHERE url=?`, url)
	post := PostDataV2Meta{}
	err := row.Scan(&post.Id, &post.Title, &post.Author, &post.Url,
		&post.CreateTime, &post.UpdateTime, &post.PrivateLevel, &post.Summary, &post.VisibleGroups,
		&post.IsDraft, &post.IsDeleted)
	if err != nil {
		log.Println(err)
	}
	log.Println(post)
	// get post content
	row = db.QueryRow(`SELECT * FROM post_content WHERE id=?`, post.Id)
	content := PostDataV2Content{}
	err = row.Scan(&content.Id, &content.Content, &content.Category, &content.Tags)
	if err != nil {
		log.Println(err)
	}
	// get post comment
	row = db.QueryRow(`SELECT * FROM post_comment WHERE id=?`, post.Id)
	comment := PostDataV2Comment{}
	err = row.Scan(&comment.Id, &comment.Likes, &comment.Dislikes, &comment.ViewCount, &comment.Comments)
	if err != nil {
		log.Println(err)
	}
	return PostDataV2{post, content, comment}
}

type V2UpdateParams struct {
	Id            int               `json:"id"`
	Meta          PostDataV2Meta    `json:"meta"`
	MetaUpdate    bool              `json:"meta_update"`
	Content       PostDataV2Content `json:"content"`
	ContentUpdate bool              `json:"content_update"`
	Comment       PostDataV2Comment `json:"comment"`
	CommentUpdate bool              `json:"comment_update"`
}

func V2UpdatePost(params V2UpdateParams) {
	db, _ := sql.Open(dbTypeV2, dbPathV2)
	defer db.Close()
	if params.MetaUpdate {
		stmt, _ := db.Prepare(`UPDATE post SET title=?, author=?, url=?,  create_time=?, update_time=?,
                					private_level=?, summary=?, visible_groups=?, is_draft=?, is_deleted=? WHERE id=?`)
		_, err := stmt.Exec(params.Meta.Title, params.Meta.Author, params.Meta.Url,
			params.Meta.CreateTime, params.Meta.UpdateTime, params.Meta.PrivateLevel, params.Meta.Summary, params.Meta.VisibleGroups,
			params.Meta.IsDraft, params.Meta.IsDeleted, params.Meta.Id)
		if err != nil {
			log.Println(err)
		}
	}
	if params.ContentUpdate {
		stmt, _ := db.Prepare(`UPDATE post_content SET content=? WHERE id=?`)
		_, err := stmt.Exec(params.Content.Content, params.Content.Id)
		if err != nil {
			log.Println(err)
		}
	}
	if params.CommentUpdate {
		stmt, _ := db.Prepare(`UPDATE post_comment SET likes=?, dislikes=?, view_count=?, comments=? WHERE id=?`)
		_, err := stmt.Exec(params.Comment.Likes, params.Comment.Dislikes, params.Comment.ViewCount, params.Comment.Comments, params.Comment.Id)
		if err != nil {
			log.Println(err)
		}
	}
}

type V2SearchParams struct {
	Author       string         `json:"author"`        // exact match
	Title        string         `json:"title"`         // use like to search
	Limit        map[string]int `json:"limit"`         // two values: start, size the number of post to return
	Sort         string         `json:"sort"`          // directly apply to sql
	Rendered     bool           `json:"rendered"`      // if true, rendered content will be returned, default false;
	CountsOnly   bool           `json:"counts_only"`   // if true, only return the count of the result, default false;
	Content      string         `json:"content"`       // use match to search
	Tags         string         `json:"tags"`          // use match to search
	Categories   string         `json:"categories"`    // use match to search
	PrivateLevel int            `json:"private_level"` // 0: public, 1: private, 2: group
	IsDraft      bool           `json:"is_draft"`      // if true, only return the draft post, default false;
	IsDeleted    bool           `json:"is_deleted"`    // if true, only return the deleted post, default false;
}

func V2SearchPosts(params V2SearchParams) ([]PostDataV2Meta, int) {
	// make query for content first to get the id
	contentCondition := ""
	contentParams := []any{}

	if params.Content != "" {
		contentCondition += `content MATCH ? `
		contentParams = append(contentParams, params.Content)
	}
	if params.Tags != "" {
		if contentCondition != "" {
			contentCondition += `AND `
		}
		contentCondition += `tags MATCH ? `

		contentParams = append(contentParams, params.Tags)
	}
	if params.Categories != "" {
		if contentCondition != "" {
			contentCondition += `AND `
		}
		contentCondition += `category MATCH ? `
		contentParams = append(contentParams, params.Categories)
	}
	if contentCondition != "" {
		contentCondition = `(SELECT id FROM post_content WHERE ` + contentCondition + `)`
	}

	log.Println(contentCondition)
	// make a query based on the id
	// select * from post where id in (select id from post_content where content match 'test')
	sqlPrepare := `SELECT * FROM post `
	wherePrepare := ""
	var prepareParams []any

	wherePrepare += `WHERE private_level <= ? `
	prepareParams = append(prepareParams, params.PrivateLevel)

	if params.IsDraft {
		wherePrepare += ` AND is_draft = 1 `
	} else {
		wherePrepare += ` AND is_draft = 0 `
	}
	if params.IsDeleted {
		wherePrepare += ` AND is_deleted = 1 `
	} else {
		wherePrepare += ` AND is_deleted = 0 `
	}
	if params.Author != "" {
		wherePrepare += ` AND author = ? `
		prepareParams = append(prepareParams, params.Author)
	}
	if params.Title != "" {
		wherePrepare += ` AND title LIKE ? `
		prepareParams = append(prepareParams, "%"+params.Title+"%")
	}
	if params.Sort != "" {
		wherePrepare += ` ORDER BY ` + params.Sort + " "
	}

	sqlPrepare += wherePrepare
	if contentCondition != "" {
		sqlPrepare += ` AND id IN ` + contentCondition
		prepareParams = append(prepareParams, contentParams...)
	}
	if params.Limit != nil {
		sqlPrepare += ` LIMIT ?,? `
		prepareParams = append(prepareParams, params.Limit["start"], params.Limit["size"])
	}
	// make a query
	db, _ := sql.Open(dbTypeV2, dbPathV2)
	defer db.Close()
	stmt, err := db.Prepare(sqlPrepare)
	if err != nil {
		log.Println(err)
	}
	rows, err := stmt.Query(prepareParams...)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	// get the result
	var result []PostDataV2Meta
	resultCount := 0
	for rows.Next() {
		resultCount++
		post := PostDataV2Meta{}
		err := rows.Scan(&post.Id, &post.Title, &post.Author, &post.Url, &post.CreateTime, &post.UpdateTime,
			&post.PrivateLevel, &post.Summary, &post.VisibleGroups, &post.IsDraft, &post.IsDeleted)
		if err != nil {
			log.Println(err)
		}
		result = append(result, post)
	}
	if params.CountsOnly {
		return []PostDataV2Meta{}, resultCount
	}
	return result, resultCount
}

type V2GetDistinctParams struct {
	Column string `json:"column"`
}

func V2GetDistinct(col string) ([]string, error) {
	db, _ := sql.Open(dbTypeV2, dbPathV2)
	// if col is a valid column in meta, then search from post
	// if col is a valid column in content, then search from content
	// if col is a valid column in comment, then search from comment
	sqlquery := ""
	if col == "author" || col == "title" || col == "url" || col == "create_time" ||
		col == "update_time" || col == "private_level" || col == "summary" ||
		col == "visible_groups" || col == "is_draft" || col == "is_deleted" {
		sqlquery = `SELECT DISTINCT ` + col + ` FROM post`
	} else if col == "content" || col == "tags" || col == "category" {
		sqlquery = `SELECT DISTINCT ` + col + ` FROM post_content`
	} else if col == "like" || col == "dislike" || col == "comment" {
		sqlquery = `SELECT DISTINCT ` + col + ` FROM post_comment`
	} else {
		return []string{}, errors.New("invalid column")
	}
	defer db.Close()
	stmt, err := db.Prepare(sqlquery)
	if err != nil {
		log.Println(err)
		return []string{}, err
	}
	rows, err := stmt.Query()
	if err != nil {
		log.Println(err)
		return []string{}, err
	}
	defer rows.Close()
	var result []string
	for rows.Next() {
		var col string
		err := rows.Scan(&col)
		if err != nil {
			log.Println(err)
		}
		result = append(result, col)
	}
	if col == "tags" {
		// seperate the tags by comma and remove the space and return the unique tags
		var uniqueTags []string
		for _, tag := range result {
			if tag == "" {
				continue
			}
			tags := strings.Split(tag, ",")
			for _, t := range tags {
				t = strings.TrimSpace(t)
				if t != "" {
					uniqueTags = append(uniqueTags, t)
				}
			}
		}
		result = uniqueTags
	}
	return result, nil
}
