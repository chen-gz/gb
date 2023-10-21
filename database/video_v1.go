package database

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type VideoItem struct {
	Id      int    `json:"id"`
	UserId  int    `json:"user_id"`
	Md5     string `json:"md5"`
	Sha256  string `json:"sha256"`
	Deleted bool   `json:"deleted"`
	Ext     string `json:"ext"` // Extension of video file. e.g. mp4, mov, avi, etc.
}

const videoTableName = "videos"

func InitVideoTableV1(videoDb *sql.DB) {
	stmt, err := videoDb.Prepare(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS( %s
    		id            INT UNSIGNED UNIQUE AUTO_INCREMENT,
			user_id       INT UNSIGNED NOT NULL,
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

func InsertVideoUser(videoDb *sql.DB, video VideoItem) (err error) {
	stmt, err := videoDb.Prepare(fmt.Sprintf(`INSERT INTO %s (
                		user_id, md5, sha256, ext ) VALUES (?, ?, ?, ?)`, videoTableName))
	if err != nil {
		panic(err)
	}
	_, err = stmt.Exec(video.UserId, video.Md5, video.Sha256)
	return err
}

func GetVideoUser(videoDb *sql.DB, user User) (video []VideoItem, err error) {
	stmt, err := videoDb.Prepare(fmt.Sprintf(`SELECT (id, user_id, md5, sha256, ext) FROM %s WHERE user_id = ?`, videoTableName))
	if err != nil {
		panic(err)
	}
	rows, err := stmt.Query(user.Id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var videoItem VideoItem
		err := rows.Scan(&videoItem.Id, &videoItem.UserId, &videoItem.Md5, &videoItem.Sha256)
		if err != nil {
			panic(err)
		}
		video = append(video, videoItem)
	}
	return video, nil
}

func GetVideoByMd5Sha256(videoDb *sql.DB, md5 string, sha256 string) (video VideoItem) {
	stmt, err := videoDb.Prepare(fmt.Sprintf(`SELECT (id, user_id, md5, sha256) FROM %s WHERE md5 = ? AND sha256 = ?`, videoTableName))
	if err != nil {
		panic(err)
	}
	row := stmt.QueryRow(md5, sha256)
	err = row.Scan(&video.Id, &video.UserId, &video.Md5, &video.Sha256)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return VideoItem{}
		}
		panic(err)
	}
	return video
}
