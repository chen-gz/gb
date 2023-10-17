package database

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

// each user has its own photo table

type PhotoItemV2 struct {
	Id          int    `json:"id"`
	OriHash     string `json:"ori_hash"`
	JpgMd5      string `json:"jpg_md5"`
	JpgSHA256   string `json:"jpg_sha256"`
	ThumbHash   string `json:"thumb_hash"`
	HasOriginal bool   `json:"has_original"`
	OriExt      string `json:"ori_ext"`
	Deleted     bool   `json:"deleted"`
	Tags        string `json:"tags"`
	Category    string `json:"category"`
}

func InitPhotoTableV2(photoDb *sql.DB, user User) error {
	tableName := fmt.Sprintf("photo_v2_%d", user.Id)
	stmt, err := photoDb.Prepare(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
    		id            INT UNSIGNED UNIQUE AUTO_INCREMENT,
    		ori_hash      VARCHAR(1024) NOT NULL,
            jpg_md5       VARCHAR(1024) NOT NULL,
    		jpg_sha256    VARCHAR(1024) NOT NULL,
    		thumb_hash    VARCHAR(1024) NOT NULL,
			has_original  BOOLEAN NOT NULL DEFAULT FALSE,
			ori_ext       VARCHAR(1024) NOT NULL DEFAULT "",
			deleted       BOOLEAN NOT NULL DEFAULT FALSE,
			tags          VARCHAR(2048) NOT NULL DEFAULT "",
			category      VARCHAR(2048) NOT NULL DEFAULT "",
    		created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    		updated_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			PRIMARY KEY (id))`, tableName))
	if err != nil {
		return err
	}
	//_, err = photo_db.Exec(query)
	_, err = stmt.Exec()
	return err
}

// create error, database exist
const (
	ErrPhotoExist    = "photo table exist"
	ErrPhotoNotExist = "photo does not exist"
)

func InsertPhotoUserV2(photoDb *sql.DB, user User, photo PhotoItemV2) (PhotoItemV2, error) {
	err := initPhotoTable(photoDb, user) // if photo table does not exist, create it
	if err != nil {
		return PhotoItemV2{}, err
	}
	tableName := fmt.Sprintf("photo_v2_%d", user.Id)
	photo, err = GetPhotoByJpgHash(photoDb, user, photo.JpgMd5, photo.JpgSHA256)
	if err == nil {
		return photo, errors.New(ErrPhotoExist)
	}

	query := fmt.Sprintf(`INSERT INTO %s (ori_hash, jpg_md5, jpg_sha256,  
                								thumb_hash, has_original, ori_ext, 
                								tags, category) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, tableName)
	stmt, err := photoDb.Prepare(query)
	if err != nil {
		panic(err)
	}
	_, err = stmt.Exec(photo.OriHash, photo.JpgMd5, photo.JpgSHA256, photo.ThumbHash,
		photo.HasOriginal, photo.OriExt, photo.Tags, photo.Category)
	if err != nil {
		panic(err)
	}
	return GetPhotoByJpgHash(photoDb, user, photo.JpgMd5, photo.JpgSHA256)
}

func GetPhotoByJpgHash(photoDb *sql.DB, user User, md5 string, sha256 string) (photo PhotoItemV2, err error) {
	tableName := fmt.Sprintf("photo_v2_%d", user.Id)
	query := fmt.Sprintf(`SELECT id, ori_hash, jpg_md5, jpg_sha256, thumb_hash, 
       							 has_original, ori_ext, deleted, tags, category 
                                 FROM %s WHERE jpg_md5 = ? AND jpg_sha256 = ?`, tableName)
	stmt, err := photoDb.Prepare(query)
	if err != nil {
		panic(err) // should not happen
	}
	row := stmt.QueryRow(md5, sha256)
	err = row.Scan(&photo.Id, &photo.OriHash, &photo.JpgMd5, &photo.JpgSHA256, &photo.ThumbHash,
		&photo.HasOriginal, &photo.OriExt, &photo.Deleted, &photo.Tags, &photo.Category)
	if err != nil {
		// photo does not exist
		return PhotoItemV2{}, errors.New(ErrPhotoNotExist)
	}
	return photo, nil
}

func GetPhotoById(photoDb *sql.DB, user User, id int) (photo PhotoItemV2, err error) {
	tableName := fmt.Sprintf("photo_v2_%d", user.Id)
	query := fmt.Sprintf(`SELECT id, ori_hash, jpg_md5, jpg_sha256, thumb_hash,  
		   							 has_original, ori_ext, deleted, tags, category
							         FROM %s WHERE id = ?`, tableName)
	stmt, err := photoDb.Prepare(query)
	if err != nil {
		panic(err) // should not happen
	}
	row := stmt.QueryRow(id)
	err = row.Scan(&photo.Id, &photo.OriHash, &photo.JpgMd5, &photo.JpgSHA256, &photo.ThumbHash,
		&photo.HasOriginal, &photo.OriExt, &photo.Deleted, &photo.Tags, &photo.Category)
	return photo, err
}

func UpdatePhotoMetaById(photoDb *sql.DB, user User, photo PhotoItemV2) error {
	// only some field can be updated - tags, category, deleted
	tableName := fmt.Sprintf("photo_v2_%d", user.Id)
	query := fmt.Sprintf(`UPDATE %s SET tags = ?, category = ?, deleted = ? WHERE id = ?`, tableName)
	stmt, err := photoDb.Prepare(query)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(photo.Tags, photo.Category, photo.Deleted, photo.Id)
	return err
}

func UpdatePhotoById(photoDb *sql.DB, user User, photo PhotoItemV2) error {
	tableName := fmt.Sprintf("photo_v2_%d", user.Id)
	query := fmt.Sprintf(`UPDATE %s SET ori_hash = ?, jpg_md5 = ?, jpg_sha256 = ?, thumb_hash = ?,
              has_original = ?, ori_ext = ?, deleted = ?, tags = ?, category = ? WHERE id = ?`, tableName)
	stmt, err := photoDb.Prepare(query)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(photo.OriHash, photo.JpgMd5, photo.JpgSHA256, photo.ThumbHash,
		photo.HasOriginal, photo.OriExt, photo.Deleted, photo.Tags, photo.Category, photo.Id)
	return err
}

func GetPhotoIds(photoDb *sql.DB, user User) (ids []int, err error) {
	tableName := fmt.Sprintf("photo_v2_%d", user.Id)
	query := fmt.Sprintf(`SELECT id FROM %s WHERE deleted = false`, tableName)
	stmt, err := photoDb.Prepare(query)
	if err != nil {
		return ids, err
	}
	rows, err := stmt.Query()
	if err != nil {
		return ids, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			return ids, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}
