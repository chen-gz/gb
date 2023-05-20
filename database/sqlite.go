package database

import (
	"database/sql"
	_ "errors"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"
)

//////////////////////////////////////////////
// database interface for v1 api

type BlogDataV1 struct {
	Id           int       `json:"id"`
	Author       string    `json:"author"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	Tags         string    `json:"tags"`
	Categories   string    `json:"categories"`
	Url          string    `json:"url"`
	Like         int       `json:"like"`
	Dislike      int       `json:"dislike"`
	CoverImg     string    `json:"cover_img"`
	IsDraft      bool      `json:"is_draft"`
	IsDeleted    bool      `json:"is_deleted"`
	PrivateLevel int       `json:"private_level"`
	ViewCount    int       `json:"view_count"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Summary      string    `json:"summary"`
	Rendered     string    `json:"rendered"`
}

const DbTypev1 = "sqlite3"
const DbPathv1 = "./blogv1.db"

func V1DeletePost(url string) error {
	database, err := sql.Open(DbTypev1, DbPathv1)
	if err != nil {
		log.Println(err)
		return err
	}
	defer database.Close()
	stmt, err := database.Prepare(`DELETE FROM posts WHERE url=?`)
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = stmt.Exec(url)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
