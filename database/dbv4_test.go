package database

import (
	"fmt"
	"github.com/minio/minio-go/v7/pkg/set"
	"strings"
	"testing"
	"unicode/utf8"
)

func TestInitV4(t *testing.T) {
	InitV4()
}
func TestV3toV4(t *testing.T) {
	// read Post from V3
	// get all url from V3
	urls, _ := V3GetDistinct("url")

	db_blog := InitV4()
	// get a post from V3
	for _, url := range urls {
		post := V3GetPostByUrl(url)
		// check summary is encoder valid in utf8 or not
		if !utf8.ValidString(post.Meta.Summary) {
			fmt.Println(url)
			post.Meta.Summary = ""
		}
		v4_post := V4PostData{
			Id:          post.Meta.Id,
			Title:       post.Meta.Title,
			Author:      post.Meta.Author,
			AuthorEmail: "",
			Url:         post.Meta.Url,
			IsDraft:     post.Meta.IsDraft,
			IsDeleted:   post.Meta.IsDeleted,
			Content:     post.Content.Content,
			Summary:     post.Meta.Summary,
			Tags:        post.Meta.Tags,
			Category:    post.Meta.Category,
			CoverImage:  post.Meta.CoverImg,
			CreatedAt:   post.Meta.CreateTime,
			UpdatedAt:   post.Meta.UpdateTime,
			ViewGroups:  set.CreateStringSet("guest"),
			EditGroups:  set.CreateStringSet("admin"),
		}
		// if content has lead empty line, remove it
		v4_post.Content = strings.TrimLeft(v4_post.Content, "\n")
		v4_post.Summary = strings.TrimLeft(v4_post.Summary, "\n")
		v4InsertPost(db_blog, v4_post)
	}
}
