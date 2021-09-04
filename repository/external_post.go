package repository

import (
	"time"

	"github.com/horri1520/hori-api/model"
	"github.com/jmoiron/sqlx"
)

func FindExternalPost(db *sqlx.DB, id int64) (*model.ExternalPost, error) {
	var externalPost model.ExternalPost

	err := db.Get(&externalPost, "select * from external_posts where id = $1", id)
	if err != nil {
		return nil, err
	}

	return &externalPost, nil
}

func AllExternalPosts(db *sqlx.DB) ([]model.ExternalPost, error) {
	var externalPosts []model.ExternalPost

	err := db.Select(&externalPosts, "select * from external_posts order by updated_at desc")
	if err != nil {
		return nil, err
	}

	return externalPosts, nil
}

func InsertExternalPost(db *sqlx.Tx, externalPost model.ExternalPost) (int64, error) {
	stmt, err := db.Preparex("insert into external_posts (title, url, thumbnail_url, created_at, updated_at, published_at) values ($1, $2, $3, $4, $5, $6) returning id")
	if err != nil {
		return 0, nil
	}

	defer func() {
		if closeErr := stmt.Close(); closeErr != nil {
			err = closeErr
		}
	}()

	var id int64
	err = stmt.QueryRow(externalPost.Title, externalPost.Url, externalPost.ThumbnailUrl, time.Now(), time.Now(), externalPost.PublishedAt).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func UpdateExternalPost(db *sqlx.Tx, externalPost *model.ExternalPost) error {
	stmt, err := db.Preparex("update external_posts set title = $1, url = $2, thumbnail_url = $3, updated_at = $4, published_at = $5 where id = $6")
	if err != nil {
		return err
	}

	defer func() {
		if closeErr := stmt.Close(); closeErr != nil {
			err = closeErr
		}
	}()

	_, err = stmt.Exec(externalPost.Title, externalPost.Url, externalPost.ThumbnailUrl, time.Now(), externalPost.PublishedAt, externalPost.Id)
	if err != nil {
		return err
	}

	return nil

}

func DeleteExternalPost(db *sqlx.Tx, id int64) error {
	stmt, err := db.Preparex("delete from external_posts where id = $1")
	if err != nil {
		return err
	}

	defer func() {
		if closeErr := stmt.Close(); closeErr != nil {
			err = closeErr
		}
	}()

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}

	return nil
}
