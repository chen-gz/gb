package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strings"
)

// the data structure still use from v3

func InitV4() (*sql.DB, *sql.DB) {
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
	// check the database eta_user is exist or not
	// if not exist, create a new database
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS eta_user")
	if err != nil {
		panic(err)
	}
	db.Close()
	// connect to database eta_blog
	db_blog, err := sql.Open("mysql", "zong:Connie@tcp(192.168.0.174:3306)/eta_blog")
	if err != nil {
		panic(err)
	}

	// for blog database, there are three tables： post_meta, post_content, post_comment
	// the post_content should be able to do full text search
	// check the table post_meta is exist or not
	// if not exist, create a new table
	_, err = db_blog.Exec(`CREATE TABLE IF NOT EXISTS post_meta (
    		id INT UNSIGNED AUTO_INCREMENT,
    		title VARCHAR(255) NOT NULL,
    		author VARCHAR(255),
    		url VARCHAR(255) UNIQUE  NOT NULL,
    		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    		private_level INT DEFAULT 0,
    		summary VARCHAR(255),
    		visibble_groups VARCHAR(255),
    		is_draft BOOLEAN DEFAULT FALSE,
    		is_deleted BOOLEAN DEFAULT FALSE,
    		tags VARCHAR(255),
    		category VARCHAR(255),
    		cover_image VARCHAR(255),
    		PRIMARY KEY (id))`)
	if err != nil {
		panic(err)
	}
	// check the table post_content is exist or not
	// if not exist, create a new table
	_, err = db_blog.Exec(`CREATE TABLE IF NOT EXISTS post_content (
    		id INT UNSIGNED AUTO_INCREMENT,
    		content TEXT,
    		categories VARCHAR(255),
    		tags VARCHAR(255),
    		PRIMARY KEY (id))`)

	if err != nil {
		panic(err)
	}
	// check the table post_comment is exist or not
	// if not exist, create a new table
	_, err = db_blog.Exec(`CREATE TABLE IF NOT EXISTS post_comment (
    		id INT UNSIGNED AUTO_INCREMENT,
    		likes INT DEFAULT 0,
    		dislikes INT DEFAULT 0,
    		view_count INT DEFAULT 0,
    		comments TEXT,
    		PRIMARY KEY (id))`)
	if err != nil {
		panic(err)
	}

	// connect to database eta_user
	db_user, err := sql.Open("mysql", "zong:Connie@tcp(192.168.0.174:3306)/eta_user")
	if err != nil {
		panic(err)
	}
	return db_blog, db_user
}
func V4InsertPost(db *sql.DB, post PostDataV3) error {
	meta := post.Meta
	content := post.Content
	comment := post.Comment
	stmt, _ := db.Prepare(`INSERT INTO post_meta (title, author, url, created_at,
		updated_at, private_level, summary, visibble_groups, is_draft, is_deleted,
		tags, category, cover_image) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)`)
	_, err := stmt.Exec(meta.Title, meta.Author, meta.Url, meta.CreateTime, meta.UpdateTime,
		meta.PrivateLevel, meta.Summary, meta.VisibleGroups, meta.IsDraft, meta.IsDeleted,
		meta.Tags, meta.Category, meta.CoverImg)
	if err != nil {
		log.Fatal(err)
		return err
	}
	// get the id of the post
	rows, err := db.Query(`SELECT id FROM post_meta WHERE url=?`, meta.Url)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&meta.Id)
		if err != nil {
			log.Fatal(err)
			return err
		}
	}
	// insert content
	content.Id = meta.Id
	content.Tags = strings.ReplaceAll(content.Tags, ",", " ")
	content.Category = strings.ReplaceAll(content.Category, ",", " ")
	comment.Id = meta.Id
	stmt, _ = db.Prepare(`INSERT INTO post_content (id, content, categories, tags) VALUES (?,?,?,?)`)
	_, err = stmt.Exec(content.Id, content.Content, content.Category, content.Tags)
	if err != nil {
		log.Fatal(err)
		return err
	}
	// insert comment
	stmt, _ = db.Prepare(`INSERT INTO post_comment (id, likes, dislikes, view_count, comments) VALUES (?,?,?,?,?)`)
	_, err = stmt.Exec(comment.Id, comment.Likes, comment.Dislikes, comment.ViewCount, comment.Comments)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func V4GetPostByUrl(db *sql.DB, url string) (PostDataV3, error) {
	var post PostDataV3
	var meta PostDataV3Meta
	var content PostDataV3Content
	var comment PostDataV3Comment
	// get meta
	rows, err := db.Query(`SELECT * FROM post_meta WHERE url=?`, url)
	if err != nil {
		log.Fatal(err)
		return post, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&meta.Id, &meta.Title, &meta.Author, &meta.Url, &meta.CreateTime, &meta.UpdateTime,
			&meta.PrivateLevel, &meta.Summary, &meta.VisibleGroups, &meta.IsDraft, &meta.IsDeleted,
			&meta.Tags, &meta.Category, &meta.CoverImg)
		if err != nil {
			log.Fatal(err)
			return post, err
		}
	}
	// get content
	rows, err = db.Query(`SELECT * FROM post_content WHERE id=?`, meta.Id)
	if err != nil {
		log.Fatal(err)
		return post, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&content.Id, &content.Content, &content.Category, &content.Tags)
		if err != nil {
			log.Fatal(err)
			return post, err
		}
	}
	// get comment
	rows, err = db.Query(`SELECT * FROM post_comment WHERE id=?`, meta.Id)
	if err != nil {
		log.Fatal(err)
		return post, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&comment.Id, &comment.Likes, &comment.Dislikes, &comment.ViewCount, &comment.Comments)
		if err != nil {
			log.Fatal(err)
			return post, err
		}
	}
	post.Meta = meta
	post.Content = content
	post.Comment = comment
	return post, nil
}
func V4UpdatePost(db *sql.DB, post PostDataV3) error {
	meta := post.Meta
	content := post.Content
	comment := post.Comment
	// update meta
	stmt, _ := db.Prepare(`UPDATE post_meta SET title=?, author=?, url=?, created_at=?,
		updated_at=?, private_level=?, summary=?, visibble_groups=?, is_draft=?, is_deleted=?,
		tags=?, category=?, cover_image=? WHERE id=?`)
	_, err := stmt.Exec(meta.Title, meta.Author, meta.Url, meta.CreateTime, meta.UpdateTime,
		meta.PrivateLevel, meta.Summary, meta.VisibleGroups, meta.IsDraft, meta.IsDeleted,
		meta.Tags, meta.Category, meta.CoverImg, meta.Id)
	if err != nil {
		log.Fatal(err)
		return err
	}
	// update content
	content.Tags = strings.ReplaceAll(content.Tags, ",", " ")
	content.Category = strings.ReplaceAll(content.Category, ",", " ")
	stmt, _ = db.Prepare(`UPDATE post_content SET content=?, categories=?, tags=? WHERE id=?`)
	_, err = stmt.Exec(content.Content, content.Category, content.Tags, content.Id)
	if err != nil {
		log.Fatal(err)
		return err
	}
	// update comment
	stmt, _ = db.Prepare(`UPDATE post_comment SET likes=?, dislikes=?, view_count=?, comments=? WHERE id=?`)
	_, err = stmt.Exec(comment.Likes, comment.Dislikes, comment.ViewCount, comment.Comments, comment.Id)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
func V4SearchPost(db *sql.DB, params V3SearchParams) ([]PostDataV3Meta, int) {
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
	db, err := sql.Open(dbTypeV3, dbPathV3)
	if err != nil {
		log.Println(err)
	}
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
