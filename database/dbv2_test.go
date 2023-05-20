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
		Tags: "Misc",
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
	post := PostDataV2Meta{
		Id: 1,
	}
	post.Title = "test"
	post.Summary = "test"
	post.VisibleGroups = "test"
	postContent := PostDataV2Content{
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
	post := V2GetPostByUrl(url)
	//post, post_content, comments := V2GetPostByUrl(url)
	log.Println(post.Meta)
	log.Println(post.Content)
	log.Println(post.Comment)
}
func TestV2GetDistinct(t *testing.T) {
	// get tags
	tags, err := V2GetDistinct("author")
	log.Println(tags)
	log.Println(err)

}
