package database

import (
	"database/sql"
	"testing"
)

func TestV2ToV3(t *testing.T) {
	// init v3 database
	InitV3()
	// migrate from v2 to v3
	db, _ := sql.Open(dbTypeV2, dbPathV2)
	defer db.Close()
	// read all post
	rows, _ := db.Query("SELECT * FROM post")
	defer rows.Close()
	for rows.Next() {
		meta := PostDataV2Meta{}
		content := PostDataV2Content{}
		comment := PostDataV2Comment{}
		err := rows.Scan(&meta.Id, &meta.Title, &meta.Author, &meta.Url, &meta.CreateTime, &meta.UpdateTime,
			&meta.PrivateLevel, &meta.Summary, &meta.VisibleGroups, &meta.IsDraft, &meta.IsDeleted)
		if err != nil {
			t.Error(err)
		}
		content.Id = meta.Id
		rows2 := db.QueryRow("SELECT * FROM post_content WHERE id = ?", meta.Id)
		//err = rows2.Scan(&content.Id, &content.Category, &content.Tags, &content.Content)
		err = rows2.Scan(&content.Id, &content.Content, &content.Category, &content.Tags)
		if err != nil {
			t.Error(err)
		}
		rows3 := db.QueryRow("SELECT * FROM post_comment WHERE id = ?", meta.Id)
		err = rows3.Scan(&comment.Id, &comment.Likes, &comment.Dislikes, &comment.ViewCount, &comment.Comments)
		if err != nil {
			t.Error(err)
		}
		metav3 := PostDataV3Meta{}
		metav3.Id = meta.Id
		metav3.Author = meta.Author
		metav3.Title = meta.Title
		metav3.Url = meta.Url
		metav3.CreateTime = meta.CreateTime
		metav3.UpdateTime = meta.UpdateTime
		metav3.PrivateLevel = meta.PrivateLevel
		metav3.Summary = meta.Summary
		metav3.VisibleGroups = meta.VisibleGroups
		metav3.IsDraft = meta.IsDraft
		metav3.IsDeleted = meta.IsDeleted
		metav3.Tags = content.Tags
		metav3.Category = content.Category
		contentv3 := PostDataV3Content{}
		contentv3.Id = meta.Id
		contentv3.Content = content.Content
		contentv3.Tags = content.Tags
		contentv3.Category = content.Category
		commentv3 := PostDataV3Comment{}
		commentv3.Id = meta.Id
		commentv3.Comments = comment.Comments
		post := PostDataV3{
			Meta:    metav3,
			Content: contentv3,
			Comment: commentv3,
		}
		err = V3InsertPost(post)
	}

}
