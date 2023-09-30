package handler

//
//import (
//	"github.com/gin-gonic/gin"
//	"go_blog/database"
//	renders "go_blog/render"
//	"log"
//	"net/http"
//	"time"
//)
//
//type GetPostRequestV3 struct {
//	Url      string `json:"url"`
//	Rendered bool   `json:"rendered"` // if true, return rendered content
//}
//type GetPostResponseV3 struct {
//	Status  string              `json:"status"`
//	Message string              `json:"message"`
//	Post    database.PostDataV3 `json:"post"`
//	Html    string              `json:"html"`
//}
//
//func V3GetPost(c *gin.Context) {
//	// get by url param
//	result := GetPostResponseV3{
//		Status: "failed",
//	}
//	var postRequest GetPostRequestV3
//	if c.BindJSON(&postRequest) != nil {
//		result.Message = "invalid request"
//		c.JSON(http.StatusBadRequest, result)
//		return
//	}
//	user := GetUserByAuthHeader(c.Request.Header.Get("Authorization"))
//	//post, post_content, post_comment
//	postData := database.V3GetPostByUrl(postRequest.Url)
//
//	if user.Level < postData.Meta.PrivateLevel {
//		result.Message = "permission denied"
//		c.JSON(http.StatusForbidden, result)
//		return
//	}
//	result.Status = "success"
//	result.Message = "ok"
//	result.Post = postData
//	//database.PostDataV2{post, post_content, post_comment}
//	result.Html = string(renders.RenderMd([]byte(postData.Content.Content)))
//	c.JSON(http.StatusOK, result)
//}
//
//type SearchPostsRequestV3 database.V3SearchParams
//type SearchPostsResponseV3 struct {
//	Status        string                    `json:"status"`
//	Message       string                    `json:"message"`
//	Posts         []database.PostDataV3Meta `json:"posts"`
//	NumberOfPosts int                       `json:"number_of_posts"`
//}
//
//func V3SearchPosts(c *gin.Context) {
//	result := SearchPostsResponseV3{
//		Status: "failed",
//	}
//	var searchRequest SearchPostsRequestV3
//	log.Println(searchRequest.IsDeleted)
//	if c.BindJSON(&searchRequest) != nil {
//		result.Message = "invalid request"
//		c.JSON(http.StatusBadRequest, result)
//		return
//	}
//	user := GetUserByAuthHeader(c.Request.Header.Get("Authorization"))
//	searchRequest.PrivateLevel = user.Level
//	if user.Role != "admin" {
//		searchRequest.IsDeleted = false
//		searchRequest.IsDraft = false
//	}
//
//	posts, cnt := database.V3SearchPosts(database.V3SearchParams(searchRequest))
//	if searchRequest.Rendered {
//		for i := 0; i < len(posts); i++ {
//			posts[i].Summary = string(renders.RenderMd([]byte(posts[i].Summary)))
//		}
//	}
//	result.Status = "success"
//	result.Message = "ok"
//	result.Posts = posts
//	result.NumberOfPosts = cnt
//	c.JSON(http.StatusOK, result)
//}
//
//type UpdatePostRequestV3 database.V3UpdateParams
//type UpdatePostResponseV3 struct {
//	Status  string `json:"status"`
//	Message string `json:"message"`
//	Url     string `json:"url"` // if the url is changed, return the new url
//}
//
//func V3UpdatePost(c *gin.Context) {
//	result := UpdatePostResponseV3{
//		Status: "failed",
//	}
//	updateRequest := UpdatePostRequestV3{}
//	if c.BindJSON(&updateRequest) != nil {
//		result.Message = "invalid request"
//		c.JSON(http.StatusBadRequest, result)
//		return
//	}
//	user := GetUserByAuthHeader(c.Request.Header.Get("Authorization"))
//	if user.Level < 1 {
//		result.Message = "permission denied"
//		c.JSON(http.StatusForbidden, result)
//		return
//	}
//	if user.Role != "admin" {
//		// && (updateRequest.MetaUpdate || updateRequest.CommentUpdate) {
//		result.Message = "permission denied"
//		c.JSON(http.StatusForbidden, result)
//		return
//	}
//	// update post
//	log.Println("updateRequest", updateRequest)
//	// set update time
//	updateRequest.Meta.UpdateTime = time.Now()
//	database.V3UpdatePost(database.V3UpdateParams(updateRequest))
//	result.Status = "success"
//	result.Message = "ok"
//	result.Url = updateRequest.Meta.Url
//	c.JSON(http.StatusOK, result)
//}
//
//type NewPostResponse struct {
//	Status  string `json:"status"`
//	Message string `json:"message"`
//	Url     string `json:"url"`
//}
//
//func V3NewPost(c *gin.Context) {
//	response := NewPostResponse{
//		Status: "failed",
//	}
//	auth := c.Request.Header.Get("Authorization")
//	user := GetUserByAuthHeader(auth)
//	if user.Role != "admin" {
//		response.Message = "permission denied"
//		c.JSON(http.StatusForbidden, response)
//	}
//	post := database.PostDataV3{
//		database.PostDataV3Meta{
//			Url:        time.Now().String(),
//			CreateTime: time.Now(),
//			UpdateTime: time.Now(),
//			IsDraft:    true,
//		},
//		database.PostDataV3Content{},
//		database.PostDataV3Comment{},
//	}
//
//	err := database.V3InsertPost(post)
//	if err != nil {
//		response.Message = "internal error"
//		c.JSON(http.StatusInternalServerError, response)
//		return
//	}
//	response.Status = "success"
//	response.Message = "ok"
//	response.Url = post.Meta.Url
//	c.JSON(http.StatusOK, response)
//}
//
//type GetDistinctRequest struct {
//	Column string `json:"column"`
//}
//type GetDistinctResponse struct {
//	Status  string   `json:"status"`
//	Message string   `json:"message"`
//	Values  []string `json:"values"`
//	Length  int      `json:"length"`
//}
//
//func V3GetDistinct(c *gin.Context) {
//	response := GetDistinctResponse{
//		Status: "failed",
//	}
//	var request GetDistinctRequest
//	if c.BindJSON(&request) != nil {
//		response.Message = "invalid request"
//		c.JSON(http.StatusBadRequest, response)
//		return
//	}
//	values, err := database.V3GetDistinct(request.Column)
//	if err != nil {
//		log.Println(err)
//		response.Message = "internal error"
//		c.JSON(http.StatusInternalServerError, response)
//		return
//	}
//	response.Status = "success"
//	response.Message = "ok"
//	response.Values = values
//	response.Length = len(values)
//	c.JSON(http.StatusOK, response)
//}
//
//type LoginResponseV3 struct {
//	Status  string `json:"status"`
//	Message string `json:"message"`
//	Email   string `json:"email"`
//	Token   string `json:"token"`
//	Name    string `json:"name"`
//}
//
//func V3Login(c *gin.Context) {
//	res := LoginResponseV3{Status: "failed"}
//	auth := c.GetHeader("Authorization")
//	user := GetUserByAuthHeader(auth)
//	if user.Email != "" && auth[0:6] == "Basic " {
//		res = LoginResponseV3{
//			Status:  "success",
//			Message: "log in success",
//			Email:   user.Email,
//			Token:   V1GenerateToken(user.Email),
//			Name:    user.Name,
//		}
//		c.JSON(http.StatusOK, res)
//	} else if user.Email != "" && auth[0:7] == "Bearer " {
//		res = LoginResponseV3{
//			Status:  "success",
//			Message: "log in success",
//			Email:   user.Email,
//			Token:   auth[7:],
//			Name:    user.Name,
//		}
//		c.JSON(http.StatusOK, res)
//	} else {
//		res = LoginResponseV3{
//			Status:  "failed",
//			Message: "log in failed",
//			Email:   "",
//			Token:   "",
//			Name:    "",
//		}
//		c.JSON(http.StatusUnauthorized, res)
//	}
//}
//
//func passwordLogin(email string, password string) bool {
//	if email == "chen-gz@outlook.com" && password == "Connie" {
//		return true
//	}
//	return false
//}
//
//type RenderMdRequestV3 struct {
//	Content string `json:"content"`
//}
//type RenderMdResponseV3 struct {
//	Status  string `json:"status"`
//	Message string `json:"message"`
//	Html    string `json:"html"`
//}
//
////func V2RenderMdV3(c *gin.Context) {
////	// get auth header
////	res := RenderMdResponseV3{
////		Status: "failed",
////	}
////
////	auth := c.Request.Header.Get("Authorization")
////	user := GetUserByAuthHeader(auth)
////
////	if user.Role != "admin" {
////		res.Message = "permission denied"
////		c.JSON(http.StatusForbidden, res)
////		return
////	}
////	// get body to string
////	body := make([]byte, 1024)
////	c.Request.Body.Read(body)
////	c.JSON(http.StatusOK, gin.H{
////		"status": "success",
////		"html":   string(renders.RenderMd(body)),
////	})
////}
