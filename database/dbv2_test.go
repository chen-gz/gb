package database

// unit test for database
import (
	"log"
	"testing"
)

func TestV2SearchPosts(t *testing.T) {
	// get count
	param := V2SearchParams{
		Title:      "Segment Tree",
		Author:     "Guangzong",
		Tags:       "Math",
		CountsOnly: true,
	}
	result, cnt := V2SearchPosts(param)
	log.Println(result)
	log.Println(cnt)

}
func TestV2SearchPosts2(t *testing.T) {
	param := V2SearchParams{
		Tags:    "Algorithm",
		Content: "Segment Tree",
	}
	result, cnt := V2SearchPosts(param)
	log.Println(result)
	log.Println(cnt)
}

func TestV2SearchPosts3(t *testing.T) {
	param := V2SearchParams{
		Tags: "Math",
	}
	result, cnt := V2SearchPosts(param)
	log.Println(result)
	log.Println(cnt)
}

func TestV2UpdatePost(t *testing.T) {
	post := BlogDataV2Meta{
		Id: 1,
	}
	post.Title = "test"
	post.Summary = "test"
	post.VisibleGroups = "test"
	postContent := BlogDataV2Content{
		Id:      1,
		Content: "test",
	}
	params := V2UpdateParams{
		Meta:          post,
		MetaUpdate:    true,
		Content:       postContent,
		ContentUpdate: true,
	}
	if params.CommentUpdate {
		log.Fatal("comment update should be false")
	}
	V2UpdatePost(params)
}

func TestV2GetPostByUrl(t *testing.T) {
	url := "20"
	post, post_content, comments := V2GetPostByUrl(url)
	log.Println(post)
	log.Println(post_content)
	log.Println(comments)
}
