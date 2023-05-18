package database

import (
	"database/sql"
	"log"
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

type BlogDataV2 struct {
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
type BlogDataV2Content struct {
	Id       int    `json:"id"`
	Content  string `json:"content"`
	Category string `json:"category"`
	Tags     string `json:"tags"`
}
type BlogDataV2Comment struct {
	Id        int    `json:"id"`
	Likes     int    `json:"likes"`
	Dislikes  int    `json:"dislikes"`
	ViewCount int    `json:"view_count"`
	Comments  string `json:"comments"`
}

func V2InsertPost(post BlogDataV2, content BlogDataV2Content, comment BlogDataV2Comment) {
	log.Println("insert post", post)
	// the id of content and comment is the same as post.
	db, _ := sql.Open(dbTypeV2, dbPathV2)
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Println(err)
		}
	}(db)
	if post.Id != 0 {
		stmt, _ := db.Prepare(`INSERT INTO post (id, title, author, url,  create_time, update_time,
					private_level, summary, visible_groups, is_draft, is_deleted) 
				    VALUES (?,?,?,?,?,?,?,?,?,?,?)`)
		_, err := stmt.Exec(post.Id, post.Title, post.Author, post.Url, post.CreateTime, post.UpdateTime,
			post.PrivateLevel, post.Summary, post.VisibleGroups, post.IsDraft, post.IsDeleted)
		if err != nil {
			log.Println(err)
		}
	} else {
		stmt, _ := db.Prepare(`INSERT INTO post (title, author, url,  create_time, update_time,
                  					private_level, summary, visible_groups, is_draft, is_deleted)
    			  					VALUES (?,?,?,?,?,?,?,?,?,?)`)
		_, err := stmt.Exec(post.Title, post.Author, post.Url, post.CreateTime, post.UpdateTime,
			post.PrivateLevel, post.Summary, post.VisibleGroups, post.IsDraft, post.IsDeleted)
		if err != nil {
			log.Println(err)
		}
		// get the id of the post
		rows, err := db.Query(`SELECT id FROM post WHERE url = ?`, post.Url)
		if err != nil {
			log.Println(err)
		}
		for rows.Next() {
			err := rows.Scan(&post.Id)
			if err != nil {
				log.Println(err)
			}
		}
	}
	content.Id = post.Id
	comment.Id = post.Id
	stmt, _ := db.Prepare(`INSERT INTO post_content (id, content, category, tags) VALUES (?,?,?,?)`)
	_, err := stmt.Exec(content.Id, content.Content, content.Category, content.Tags)
	if err != nil {
		log.Println(err)
	}
	stmt, _ = db.Prepare(`INSERT INTO post_comment (id, likes, dislikes, view_count, comments) VALUES (?,?,?,?,?)`)
	_, err = stmt.Exec(comment.Id, comment.Likes, comment.Dislikes, comment.ViewCount, comment.Comments)
	if err != nil {
		log.Println(err)
	}
}

func V2GetPostByUrl(url string) (BlogDataV2, BlogDataV2Content, BlogDataV2Comment) {
	db, _ := sql.Open(dbTypeV2, dbPathV2)
	defer db.Close()
	row := db.QueryRow(`SELECT * FROM post WHERE url=?`, url)
	post := BlogDataV2{}
	err := row.Scan(&post.Id, &post.Title, &post.Author, &post.Url,
		&post.CreateTime, &post.UpdateTime, &post.PrivateLevel, &post.Summary, &post.VisibleGroups,
		&post.IsDraft, &post.IsDeleted)
	if err != nil {
		log.Println(err)
	}
	log.Println(post)
	// get post content
	row = db.QueryRow(`SELECT * FROM post_content WHERE id=?`, post.Id)
	content := BlogDataV2Content{}
	err = row.Scan(&content.Id, &content.Content, &content.Category, &content.Tags)
	if err != nil {
		log.Println(err)
	}
	// get post comment
	row = db.QueryRow(`SELECT * FROM post_comment WHERE id=?`, post.Id)
	comment := BlogDataV2Comment{}
	err = row.Scan(&comment.Id, &comment.Likes, &comment.Dislikes, &comment.ViewCount, &comment.Comments)
	if err != nil {
		log.Println(err)
	}
	return post, content, comment
}

type V2UpdateParams struct {
	Id            int
	Post          BlogDataV2
	PostUpdate    bool
	Content       BlogDataV2Content
	ContentUpdate bool
	Comment       BlogDataV2Comment
	CommentUpdate bool
}

func V2UpdatePost(params V2UpdateParams) {
	db, _ := sql.Open(dbTypeV2, dbPathV2)
	defer db.Close()
	if params.PostUpdate {
		stmt, _ := db.Prepare(`UPDATE post SET title=?, author=?, url=?,  create_time=?, update_time=?,
                					private_level=?, summary=?, visible_groups=?, is_draft=?, is_deleted=? WHERE id=?`)
		_, err := stmt.Exec(params.Post.Title, params.Post.Author, params.Post.Url,
			params.Post.CreateTime, params.Post.UpdateTime, params.Post.PrivateLevel, params.Post.Summary, params.Post.VisibleGroups,
			params.Post.IsDraft, params.Post.IsDeleted, params.Post.Id)
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
	Author     string         `json:"author"`     // exact match
	Title      string         `json:"title"`      // use like to search
	Limit      map[string]int `json:"limit"`      // two values: start, size the number of post to return
	Sort       string         `json:"sort"`       // directly apply to sql
	Rendered   bool           `json:"rendered"`   // if true, rendered content will be returned, default false;
	Count      bool           `json:"count"`      // if true, only return the count of the result
	Content    string         `json:"content"`    // use match to search
	Tags       string         `json:"tags"`       // use match to search
	Categories string         `json:"categories"` // use match to search
}

func V2SearchPosts(params V2SearchParams) ([]BlogDataV2, int) {
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

	sqlPrepare := ""
	//sqlFront := `SELECT * FROM post `
	sqlPrepare += `SELECT * FROM post `
	wherePrepare := ""
	var prepareParams []any
	hasCondition := false

	if params.Author != "" {
		if hasCondition {
			wherePrepare += ` AND `
		}
		wherePrepare += `author = ? `
		prepareParams = append(prepareParams, params.Author)
		hasCondition = true
	}
	if params.Title != "" {
		if hasCondition {
			wherePrepare += ` AND `
		}
		wherePrepare += `title LIKE ? `
		prepareParams = append(prepareParams, "%"+params.Title+"%")
		hasCondition = true
	}
	if params.Limit != nil {
		wherePrepare += `LIMIT ?,? `
		prepareParams = append(prepareParams, params.Limit["start"], params.Limit["size"])
	}
	if params.Sort != "" {
		wherePrepare += `ORDER BY ? `
		prepareParams = append(prepareParams, params.Sort)
	}
	if contentCondition != "" {
		if wherePrepare != "" {
			sqlPrepare += `WHERE ` + wherePrepare + ` AND id IN ` + contentCondition
			prepareParams = append(prepareParams, contentParams...)
		} else {
			sqlPrepare += `WHERE id IN ` + contentCondition
			prepareParams = append(prepareParams, contentParams...)
		}
	} else {
		if wherePrepare != "" {
			sqlPrepare += `WHERE ` + wherePrepare
		}
	}
	// make a query
	db, _ := sql.Open(dbTypeV2, dbPathV2)
	defer db.Close()
	log.Println(sqlPrepare)
	log.Println(prepareParams)
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
	var result []BlogDataV2
	resultCount := 0
	for rows.Next() {
		resultCount++
		post := BlogDataV2{}
		err := rows.Scan(&post.Id, &post.Title, &post.Author, &post.Url, &post.CreateTime, &post.UpdateTime,
			&post.PrivateLevel, &post.Summary, &post.VisibleGroups, &post.IsDraft, &post.IsDeleted)
		if err != nil {
			log.Println(err)
		}
		result = append(result, post)
	}
	if params.Count {
		return []BlogDataV2{}, resultCount
	}
	return result, resultCount
}
