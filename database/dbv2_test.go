package database

// unit test for database
import (
	"log"
	"testing"
)

func TestSearchV2(t *testing.T) {
	param := V2SearchParams{
		Tags: "Math",
	}
	result, cnt := V2SearchPosts(param)
	log.Println(result)
	log.Println(cnt)
}
func TestV2SearchPosts(t *testing.T) {
	// get count
	param := V2SearchParams{
		Tags:  "Math",
		Count: true,
	}
	result, cnt := V2SearchPosts(param)
	log.Println(result)
	log.Println(cnt)
}
