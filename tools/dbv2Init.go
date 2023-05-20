package main

import (
	"database/sql"
	"go_blog/database"
	"log"
	"math"
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
		meta := database.PostDataV2Meta{}
		meta.Id = post.Id
		meta.Author = post.Author
		meta.Title = post.Title
		meta.Url = post.Url
		meta.CreateTime = post.CreatedAt
		meta.UpdateTime = post.UpdatedAt
		meta.PrivateLevel = post.PrivateLevel
		meta.Summary = ""
		meta.VisibleGroups = ""
		meta.IsDraft = post.IsDraft
		meta.IsDeleted = post.IsDeleted
		content := database.PostDataV2Content{}
		content.Id = post.Id
		content.Category = post.Categories
		content.Tags = post.Tags
		content.Content = post.Content
		post2 := database.PostDataV2{
			Meta:    meta,
			Content: content,
			Comment: database.PostDataV2Comment{},
		}
		err = database.V2InsertPost(post2)
		if err != nil {
			return
		}
		//meta, content, database.PostDataV2Comment{})
	}
}
func V2SummaryUpdate() {
	param := database.V2SearchParams{
		PrivateLevel: 10,
	}
	posts, _ := database.V2SearchPosts(param)
	for _, post := range posts {
		np := database.V2GetPostByUrl(post.Url)
		np.Meta.Summary = np.Content.Content[:int(math.Min(300, float64(len(np.Content.Content))))]
		updateParam := database.V2UpdateParams{
			Id:            np.Meta.Id,
			Meta:          np.Meta,
			MetaUpdate:    true,
			ContentUpdate: false,
			CommentUpdate: false,
		}
		database.V2UpdatePost(updateParam)
	}

}

func main() {
	database.InitV2()
	Migrate1to2()
	V2SummaryUpdate()
}
