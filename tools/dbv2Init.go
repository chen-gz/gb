package main

import (
	"database/sql"
	"go_blog/database"
	"log"
)

func Migrate1to2() {
	// open v1 database
	log.Println("migrate v1 to v2")
	db, err := sql.Open(database.DbTypev1, database.DbPathv1)
	if err != nil {
		log.Println(err)
	}
	// reade all post
	rows, _ := db.Query("SELECT * FROM posts")
	defer rows.Close()
	for rows.Next() {
		post := database.BlogDataV1{}
		err := rows.Scan(&post.Id, &post.Author, &post.Title, &post.Content, &post.Tags, &post.Categories,
			&post.Url, &post.Like, &post.Dislike, &post.CoverImg, &post.IsDraft, &post.IsDeleted, &post.PrivateLevel, &post.ViewCount,
			&post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			log.Println(err)
		}
		// insert to v2 database
		postv2 := database.BlogDataV2Meta{}
		postv2.Id = post.Id
		postv2.Author = post.Author
		postv2.Title = post.Title
		postv2.Url = post.Url
		postv2.CreateTime = post.CreatedAt
		postv2.UpdateTime = post.UpdatedAt
		postv2.PrivateLevel = post.PrivateLevel
		postv2.Summary = ""
		postv2.VisibleGroups = ""
		postv2.IsDraft = post.IsDraft
		postv2.IsDeleted = post.IsDeleted
		postv2content := database.BlogDataV2Content{}
		postv2content.Id = post.Id
		postv2content.Category = post.Categories
		postv2content.Tags = post.Tags
		postv2content.Content = post.Content
		database.V2InsertPost(postv2, postv2content, database.BlogDataV2Comment{})
	}
}

func main() {
	database.InitV2()
	Migrate1to2()

}
