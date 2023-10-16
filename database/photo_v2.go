package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

// each user has its own photo table

type PhotoItemV2 struct {
	Id          int    `json:"id"`
	OriHash     string `json:"ori_hash"`
	JpgHash     string `json:"jpg_hash"`
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
    		jpg_hash      VARCHAR(1024) NOT NULL,
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

func getPhotoById(photoDb *sql.DB, user User, id int) (photo PhotoItemV2, err error) {
	tableName := fmt.Sprintf("photo_v2_%d", user.Id)
	stmt, err := photoDb.Prepare(fmt.Sprintf(`SELECT id, ori_hash, jpg_hash, has_original, ori_ext, deleted, tags, category FROM %s WHERE id = ?`, tableName))
	if err != nil {
		return PhotoItemV2{}, err
	}
	row := stmt.QueryRow(id)
	err = row.Scan(&photo.Id, &photo.OriHash, &photo.JpgHash, &photo.HasOriginal, &photo.OriExt, &photo.Deleted, &photo.Tags, &photo.Category)
	return photo, err
}
func updatePhotoFileById(photoDb *sql.DB, user User, photo PhotoItemV2) error {
	tableName := fmt.Sprintf("photo_v2_%d", user.Id)
	query := fmt.Sprintf(`UPDATE %s SET jpg_hash = ?, has_original = ?, ori_ext = ? WHERE id = ?`, tableName)
	stmt, err := photoDb.Prepare(query)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(photo.JpgHash, photo.HasOriginal, photo.OriExt, photo.Id)
	return err
}

func InsertPhotoUserV2(photoDb *sql.DB, user User, photo PhotoItemV2) (PhotoItemV2, error) {
	err := initPhotoTable(photoDb, user)
	if err != nil {
		return PhotoItemV2{}, err
	}
	tableName := fmt.Sprintf("photo_v2_%d", user.Id)
	// check if jpg_hash exists
	query := fmt.Sprintf(`select id from %s where jpg_hash = ?`, tableName)
	stmt, err := photoDb.Prepare(query)
	if err != nil {
		return PhotoItemV2{}, err
	}
	row := stmt.QueryRow(photo.JpgHash)
	var id int
	err = row.Scan(&id)
	if err == nil {
		// jpg_hash exists. return error
		return PhotoItemV2{}, fmt.Errorf("jpg_hash exists")
	}

	query = fmt.Sprintf(`INSERT INTO %s (ori_hash, jpg_hash, thumb_hash, has_original, ori_ext, tags, category) VALUES (?, ?, ?, ?, ?, ?, ?)`, tableName)
	stmt, err = photoDb.Prepare(query)
	if err != nil {
		return PhotoItemV2{}, err
	}
	_, err = stmt.Exec(photo.OriHash, photo.JpgHash, photo.ThumbHash, photo.HasOriginal, photo.OriExt, photo.Tags, photo.Category)
	if err != nil {
		return PhotoItemV2{}, err
	}
	return GetPhotoByJpgHash(photoDb, user, photo.JpgHash)
}

func GetPhotoByJpgHash(photoDb *sql.DB, user User, hash string) (photo PhotoItemV2, err error) {
	tableName := fmt.Sprintf("photo_v2_%d", user.Id)
	query := fmt.Sprintf(`SELECT id, ori_hash, jpg_hash, thumb_hash, has_original, ori_ext, deleted, tags, category FROM %s WHERE jpg_hash = ?`, tableName)
	stmt, err := photoDb.Prepare(query)
	if err != nil {
		return PhotoItemV2{}, err
	}
	row := stmt.QueryRow(hash)
	err = row.Scan(&photo.Id, &photo.OriHash, &photo.JpgHash, &photo.ThumbHash, &photo.HasOriginal, &photo.OriExt, &photo.Deleted, &photo.Tags, &photo.Category)
	return photo, err
}

func GetPhotoById(photoDb *sql.DB, user User, id int) (photo PhotoItemV2, err error) {
	tableName := fmt.Sprintf("photo_v2_%d", user.Id)
	query := fmt.Sprintf(`SELECT id, ori_hash, jpg_hash, thumb_hash, has_original, ori_ext, deleted, tags, category FROM %s WHERE id = ?`, tableName)
	//stmt, err := photoDb.Prepare(`SELECT id, ori_hash, jpg_hash, has_original, ori_ext, deleted, tags, category FROM ? WHERE id = ?`)
	stmt, err := photoDb.Prepare(query)
	if err != nil {
		return PhotoItemV2{}, err
	}
	row := stmt.QueryRow(id)
	err = row.Scan(&photo.Id, &photo.OriHash, &photo.JpgHash, &photo.ThumbHash, &photo.HasOriginal, &photo.OriExt, &photo.Deleted, &photo.Tags, &photo.Category)
	return photo, err
}

func UpdatePhotoMetaById(photoDb *sql.DB, user User, photo PhotoItemV2) error {
	// only some field can be updated - tags, category, deleted
	tableName := fmt.Sprintf("photo_v2_%d", user.Id)
	query := fmt.Sprintf(`UPDATE %s SET tags = ?, category = ?, deleted = ? WHERE id = ?`, tableName)
	//stmt, err := photoDb.Prepare(`UPDATE ? SET tags = ?, category = ?, deleted = ? WHERE id = ?`)
	stmt, err := photoDb.Prepare(query)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(photo.Tags, photo.Category, photo.Deleted, photo.Id)
	return err
}

func UpdatePhotoById(photoDb *sql.DB, user User, photo PhotoItemV2) error {
	tableName := fmt.Sprintf("photo_v2_%d", user.Id)
	query := fmt.Sprintf(`UPDATE %s SET ori_hash = ?, jpg_hash = ?, has_original = ?, ori_ext = ?, deleted = ?, tags = ?, category = ? WHERE id = ?`, tableName)
	//stmt, err := photoDb.Prepare(`UPDATE ? SET ori_hash = ?, jpg_hash = ?, has_original = ?, ori_ext = ?, deleted = ?, tags = ?, category = ? WHERE id = ?`)
	stmt, err := photoDb.Prepare(query)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(photo.OriHash, photo.JpgHash, photo.HasOriginal, photo.OriExt, photo.Deleted, photo.Tags, photo.Category, photo.Id)
	return err
}
