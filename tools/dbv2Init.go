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
		postv2 := database.PostDataV2Meta{}
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
		postv2content := database.PostDataV2Content{}
		postv2content.Id = post.Id
		postv2content.Category = post.Categories
		postv2content.Tags = post.Tags
		postv2content.Content = post.Content
		post2 := database.PostDataV2{
			Meta:    postv2,
			Content: postv2content,
			Comment: database.PostDataV2Comment{},
		}
		err = database.V2InsertPost(post2)
		if err != nil {
			return
		}
		//postv2, postv2content, database.PostDataV2Comment{})
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
	//database.InitV2()
	//Migrate1to2()
	V2SummaryUpdate()

}
