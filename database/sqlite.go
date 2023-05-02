package database

import (
    "database/sql"
    "fmt"
    _ "github.com/mattn/go-sqlite3"
    "log"
    "strings"
    "time"
    // "os"
)

func Init() {
    fmt.Println("hello")
    // create sqlite database
    database, err := sql.Open("sqlite3", "./blog.db")
    if err != nil {
        log.Fatal(err)
    }
    defer database.Close()
    database.Exec(`CREATE TABLE IF NOT EXISTS posts (id INTEGER PRIMARY KEY AUTOINCREMENT, title TEXT, author TEXT, content TEXT, tags TEXT, categories TEXT, datetime DATETIME)`)
}
var dbPath = "./blog.db"
var dbType = "sqlite3"

type BlogData struct {
    Id         int
    Name       string // auto generate for s3 storage
    Author     string
    Title      string
    Content    string
    Tags       []string
    Categories []string
    Datetime   time.Time
}

func InsertPost(blog BlogData) {
    log.Println("insert post")
    database, err := sql.Open("sqlite3", "./blog.db")
    if err != nil {
        log.Fatal(err)
    }
    defer database.Close()
    // check if the data exist
    query := database.QueryRow("SELECT * FROM posts WHERE Id = ?", blog.Id)

    // print the query result
    post := BlogData{}
    err = query.Scan(&post.Id, &post.Title, &post.Content, &post.Tags, &post.Categories, &post.Datetime)

    stmt, error := database.Prepare("INSERT INTO posts(title, author, content, tags, categories, datetime) VALUES(?, ?, ?, ?, ?, ?)")
    if error != nil {
        log.Fatal(error)
    }
    // connect Tags to a single string and surround with '[]'
    Tags := "[" + strings.Join(blog.Tags, ",") + "]"
    Categories := "[" + strings.Join(blog.Categories, ",") + "]"
    tmp, err := stmt.Exec( blog.Title, blog.Author, blog.Content, Tags, Categories, blog.Datetime) 
    if error != nil {
        log.Fatal(error)
    }
    if tmp != nil {
        log.Println("insert success")
    }else {
        log.Println("insert failed")
    }
}
func UpdatePost(blog BlogData) {
    log.Println("update post")
    database, err := sql.Open(dbType, dbPath)
    if err != nil {
        log.Fatal(err)
    }
    defer database.Close()
    // update the post based on the id
    stmt, error := database.Prepare("UPDATE posts SET title = ?, content = ?, tags = ?, categories = ?, datetime = ?, author = ?, WHERE id = ?")
    if error != nil {
        log.Fatal(error)
    }
    // connect Tags to a single string and connect use ',' to connect
    Tags := strings.Join(blog.Tags, ",")
    Categories := strings.Join(blog.Categories, "/")
    
    tmp, err := stmt.Exec(blog.Title, blog.Content, Tags, Categories, blog.Datetime, blog.Author, blog.Id)
    if error != nil {
        log.Fatal(error)
    }
    if tmp != nil {
        log.Println("update success")
    }else {
        log.Println("update failed")
    }
}

func GetAllPostIdAndName() (map[int]string, error) {
    // get all post id and name
    database, err := sql.Open(dbType, dbPath)
    stmt, err := database.Prepare("SELECT id, title FROM posts")
    if err != nil {
        log.Fatal(err)
    }
    defer database.Close()
    rows, err := stmt.Query()
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()
    var id int
    var title string
    result := make(map[int]string)
    for rows.Next() {
        err := rows.Scan(&id, &title)
        if err != nil {
            log.Fatal(err)
        }
        // log.Println(id, title)
        result[id] = title
    }
    return result, err
}

func GetPostByIndex(index int) (BlogData, error) {
    database, err := sql.Open(dbType, dbPath)
    if err != nil {
        log.Fatal(err)
    }
    defer database.Close()
    query := database.QueryRow("SELECT * FROM posts WHERE Id = ?", index)
    post := BlogData{}
    err = query.Scan(&post.Id, &post.Title, &post.Author,&post.Content, &post.Tags, &post.Categories, &post.Datetime)
    return post, err
}


