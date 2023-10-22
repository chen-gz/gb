package database

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"go_blog/interfaces"
)

const videoTableName = "videos"

func initVideoTableV1(videoDb *sql.DB) {
	stmt, err := videoDb.Prepare(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
    		id            INT UNSIGNED UNIQUE AUTO_INCREMENT,
			user_id       INT UNSIGNED NOT NULL,
			title   	  VARCHAR(1024) NOT NULL default "",
			tag  		  VARCHAR(1024) NOT NULL default "",
			md5 		  VARCHAR(1024) NOT NULL,
			sha256 		  VARCHAR(1024) NOT NULL,
			ext 		  VARCHAR(1024) NOT NULL,
			deleted       BOOLEAN NOT NULL DEFAULT FALSE,
    		created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    		updated_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			PRIMARY KEY (id))`, videoTableName))
	if err != nil {
		panic(err)
	}
	_, err = stmt.Exec()
	if err != nil {
		panic(err)
	}
}

var dbConfig interfaces.DbConfig

// InitVideoDb creates a database if not exist, and return a sql.DB object.
// Any error will cause panic.
func InitVideoDb(_config interfaces.DbConfig) *sql.DB {
	dbConfig = _config
	sqlEndpoint := fmt.Sprintf("%s:%s@%s/", dbConfig.User, dbConfig.Password, dbConfig.Address)
	db, err := sql.Open("mysql", sqlEndpoint)
	if err != nil {
		panic(err)
	}
	query := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbConfig.DatabaseName)
	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}
	err = db.Close()
	if err != nil {
		panic(err)
	}

	sqlEndpoint = fmt.Sprintf("%s:%s@%s/%s", dbConfig.User, dbConfig.Password, dbConfig.Address, dbConfig.DatabaseName)
	db, err = sql.Open("mysql", sqlEndpoint)
	if err != nil {
		panic(err)
	}
	initVideoTableV1(db)
	return db
}

func InsertVideoUser(videoDb *sql.DB, video interfaces.VideoItem) (err error) {
	stmt, err := videoDb.Prepare(fmt.Sprintf(`INSERT INTO %s 
               		(user_id, md5, sha256, ext, title) VALUES (?, ?, ?, ?, ?)`, videoTableName))
	if err != nil {
		panic(err)
	}
	_, err = stmt.Exec(video.UserId, video.Md5, video.Sha256, video.Ext, video.Title)
	return err
}

func GetVideoUser(videoDb *sql.DB, user User) (video []interfaces.VideoItem) {
	stmt, err := videoDb.Prepare(fmt.Sprintf(`SELECT id, user_id, md5, sha256, ext FROM %s WHERE user_id = ?`, videoTableName))
	if err != nil {
		panic(err)
	}
	rows, err := stmt.Query(user.Id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var videoItem interfaces.VideoItem
		err := rows.Scan(&videoItem.Id, &videoItem.UserId, &videoItem.Md5, &videoItem.Sha256, &videoItem.Ext)
		if err != nil {
			panic(err)
		}
		video = append(video, videoItem)
	}
	return video
}

func GetVideoByMd5Sha256(videoDb *sql.DB, user_id int, md5 string, sha256 string) (video interfaces.VideoItem) {
	stmt, err := videoDb.Prepare(fmt.Sprintf(`SELECT id, user_id, md5, sha256 FROM %s WHERE md5 = ? AND sha256 = ? AND user_id = ?`, videoTableName))
	if err != nil {
		panic(err)
	}
	row := stmt.QueryRow(md5, sha256, user_id)
	err = row.Scan(&video.Id, &video.UserId, &video.Md5, &video.Sha256)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return interfaces.VideoItem{}
		}
		panic(err)
	}
	return video
}

func GetVideoById(videoDb *sql.DB, id int, user_id int) (video interfaces.VideoItem) {
	stmt, err := videoDb.Prepare(fmt.Sprintf(`SELECT id, user_id, md5, sha256 FROM %s WHERE id = ? AND user_id = ?`, videoTableName))
	if err != nil {
		panic(err)
	}
	row := stmt.QueryRow(id, user_id)
	err = row.Scan(&video.Id, &video.UserId, &video.Md5, &video.Sha256)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return interfaces.VideoItem{}
		}
		panic(err)
	}
	return video
}
