package database

import (
	"database/sql"
	"errors"
	"log"
	"strings"
	"time"
)

// add tags and categories, cover_img to meta table

// initialize v3 database
const dbPathV3 = "v3.db"
const dbTypeV3 = "sqlite3"

func InitV3() {
	db, _ := sql.Open(dbTypeV3, dbPathV3)
	defer db.Close()
	// create post table
	stmt, _ := db.Prepare(`CREATE TABLE IF NOT EXISTS post_meta (
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
								 is_deleted integer,
                                 tags text,
                                 category text,
                                 cover_img text
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

type PostDataV3Meta struct {
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
	Tags          string    `json:"tags"`
	Category      string    `json:"category"`
	CoverImg      string    `json:"cover_img"`
}
type PostDataV3Content struct {
	Id       int    `json:"id"`
	Content  string `json:"content"`
	Category string `json:"category"`
	Tags     string `json:"tags"`
}
type PostDataV3Comment struct {
	Id        int    `json:"id"`
	Likes     int    `json:"likes"`
	Dislikes  int    `json:"dislikes"`
	ViewCount int    `json:"view_count"`
	Comments  string `json:"comments"`
}
type PostDataV3 struct {
	Meta    PostDataV3Meta    `json:"meta"`
	Content PostDataV3Content `json:"content"`
	Comment PostDataV3Comment `json:"comment"`
}

func V3InsertPost(postData PostDataV3) error {
	meta := postData.Meta
	content := postData.Content
	comment := postData.Comment
	// the id of content and comment should be the same as post.
	db, _ := sql.Open(dbTypeV3, dbPathV3)
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Println(err)
		}
	}(db)
	stmt, _ := db.Prepare(`INSERT INTO post_meta (title, author, url,  create_time, update_time,
                  					private_level, summary, visible_groups, is_draft, is_deleted, tags, category, cover_img)
    			  					VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)`)
	_, err := stmt.Exec(meta.Title, meta.Author, meta.Url, meta.CreateTime, meta.UpdateTime,
		meta.PrivateLevel, meta.Summary, meta.VisibleGroups, meta.IsDraft, meta.IsDeleted, meta.Tags, meta.Category, meta.CoverImg)
	if err != nil {
		log.Println(err)
		return err
	}
	// get the id of the post
	rows, err := db.Query(`SELECT id FROM post_meta WHERE url = ?`, meta.Url)
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
	content.Tags = strings.ReplaceAll(content.Tags, ",", " ")
	content.Category = strings.ReplaceAll(content.Category, ",", " ")
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

func V3GetPostByUrl(url string) PostDataV3 {
	db, _ := sql.Open(dbTypeV3, dbPathV3)
	defer db.Close()
	row := db.QueryRow(`SELECT * FROM post_meta WHERE url=?`, url)
	post := PostDataV3Meta{}
	err := row.Scan(&post.Id, &post.Title, &post.Author, &post.Url,
		&post.CreateTime, &post.UpdateTime, &post.PrivateLevel, &post.Summary, &post.VisibleGroups,
		&post.IsDraft, &post.IsDeleted, &post.Tags, &post.Category, &post.CoverImg)
	if err != nil {
		log.Println(err)
	}
	log.Println(post)
	// get post content
	row = db.QueryRow(`SELECT * FROM post_content WHERE id=?`, post.Id)
	content := PostDataV3Content{}
	err = row.Scan(&content.Id, &content.Content, &content.Category, &content.Tags)
	if err != nil {
		log.Println(err)
	}
	// get post comment
	row = db.QueryRow(`SELECT * FROM post_comment WHERE id=?`, post.Id)
	comment := PostDataV3Comment{}
	err = row.Scan(&comment.Id, &comment.Likes, &comment.Dislikes, &comment.ViewCount, &comment.Comments)
	if err != nil {
		log.Println(err)
	}
	return PostDataV3{post, content, comment}
}

type V3UpdateParams struct {
	Id            int               `json:"id"`
	Meta          PostDataV3Meta    `json:"meta"`
	MetaUpdate    bool              `json:"meta_update"`
	Content       PostDataV3Content `json:"content"`
	ContentUpdate bool              `json:"content_update"`
	Comment       PostDataV3Comment `json:"comment"`
	CommentUpdate bool              `json:"comment_update"`
}

func V3UpdatePost(params V3UpdateParams) {
	db, _ := sql.Open(dbTypeV3, dbPathV3)
	defer db.Close()
	if params.MetaUpdate {
		stmt, _ := db.Prepare(`UPDATE post_meta SET title=?, author=?, url=?,  create_time=?, update_time=?,
                					private_level=?, summary=?, visible_groups=?, is_draft=?, is_deleted=?, tags=?, 
                					category=?, cover_img=? WHERE id=?`)
		_, err := stmt.Exec(params.Meta.Title, params.Meta.Author, params.Meta.Url,
			params.Meta.CreateTime, params.Meta.UpdateTime, params.Meta.PrivateLevel, params.Meta.Summary, params.Meta.VisibleGroups,
			params.Meta.IsDraft, params.Meta.IsDeleted, params.Meta.Tags, params.Meta.Category, params.Meta.CoverImg, params.Meta.Id)
		stmt, _ = db.Prepare(`UPDATE post_content SET category=?, tags=? WHERE id=?`)
		_, err = stmt.Exec(params.Meta.Category, params.Meta.Tags, params.Meta.Id)
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

type V3SearchParams struct {
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

func V3SearchPosts(params V3SearchParams) ([]PostDataV3Meta, int) {
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
	sqlPrepare := `SELECT * FROM post_meta `
	wherePrepare := ""
	var prepareParams []any

	wherePrepare += `WHERE private_level <= ? `
	prepareParams = append(prepareParams, params.PrivateLevel)

	if params.IsDraft {
		wherePrepare += ` AND is_draft = 1 `
	} else {
		wherePrepare += ` AND is_draft = false `
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

	sqlPrepare += wherePrepare
	if contentCondition != "" {
		sqlPrepare += ` AND id IN ` + contentCondition
		prepareParams = append(prepareParams, contentParams...)
	}

	if params.Sort != "" {
		sqlPrepare += ` ORDER BY ` + params.Sort + " "
	}
	if params.Limit != nil {
		sqlPrepare += ` LIMIT ?,? `
		prepareParams = append(prepareParams, params.Limit["start"], params.Limit["size"])
	}
	// make a query
	db, _ := sql.Open(dbTypeV3, dbPathV3)
	defer db.Close()
	log.Println(sqlPrepare)
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
	var result []PostDataV3Meta
	resultCount := 0
	for rows.Next() {
		resultCount++
		post := PostDataV3Meta{}
		err := rows.Scan(&post.Id, &post.Title, &post.Author, &post.Url, &post.CreateTime, &post.UpdateTime,
			&post.PrivateLevel, &post.Summary, &post.VisibleGroups, &post.IsDraft, &post.IsDeleted, &post.Tags, &post.Category, &post.CoverImg)
		if err != nil {
			log.Println(err)
		}
		result = append(result, post)
	}
	if params.CountsOnly {
		return []PostDataV3Meta{}, resultCount
	}
	return result, resultCount
}

type V3GetDistinctParams struct {
	Column string `json:"column"`
}

func V3GetDistinct(col string) ([]string, error) {
	db, _ := sql.Open(dbTypeV3, dbPathV3)
	// if col is a valid column in meta, then search from post
	// if col is a valid column in content, then search from content
	// if col is a valid column in comment, then search from comment
	sqlquery := ""
	if col == "author" || col == "title" || col == "url" || col == "create_time" ||
		col == "update_time" || col == "private_level" || col == "summary" ||
		col == "visible_groups" || col == "is_draft" || col == "is_deleted" {
		sqlquery = `SELECT DISTINCT ` + col + ` FROM post_meta`
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
