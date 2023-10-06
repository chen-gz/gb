package database

import (
	"fmt"
	"github.com/minio/minio-go/v7/pkg/set"
	"github.com/stretchr/testify/assert"
	"log"
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
	// drop the table blog
	_, err := db_blog.Exec(`DROP TABLE IF EXISTS post`)
	db_blog.Close()

	db_blog = InitV4()

	if err != nil {
		log.Fatal(err)
	}

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
			ViewGroups:  set.CreateStringSet("admin", "editor", "author", "premium", "subscriber", "guest"),
			EditGroups:  set.CreateStringSet("admin", "editor", "author"),
		}
		// if content has lead empty line, remove it
		v4_post.Content = strings.TrimLeft(v4_post.Content, "\n")
		v4_post.Summary = strings.TrimLeft(v4_post.Summary, "\n")
		v4InsertPostMigrate(db_blog, v4_post)
	}
}

func TestSearchPosts(t *testing.T) {
	db_blog := InitV4()
	searchRequest := SearchParams{
		Author: "Guangzong",
	}
	db_user, _ := UserDbInit()
	// add user to blog_user

	user := GetUserByEmail(db_user, "chen-gz@outlook.com")

	post, err := searchPosts(db_blog, searchRequest, user)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(post)
}

func TestV4InsertPosByUser(t *testing.T) {

	db_blog := InitV4()
	db_user, _ := UserDbInit()
	user := GetUserByEmail(db_user, "chen-gz@outlook.com")
	fmt.Println(user)
	v4InsertUser(db_blog, V4BlogUserData{
		Email: user.Email,
		Name:  user.Name,
		Roles: set.CreateStringSet("admin", "editor", "author"),
	})
}

func TestGetUserRole(t *testing.T) {

	db_blog := InitV4()
	db_user, _ := UserDbInit()
	user := GetUserByEmail(db_user, "chen-gz@outlook.com")
	fmt.Println(user)
	roles_str := ""
	err := db_blog.QueryRow(`SELECT roles FROM blog_users WHERE email=?`, user.Email).Scan(&roles_str)
	fmt.Println(err)
	assert.Nil(t, err, "should be nil")

}
