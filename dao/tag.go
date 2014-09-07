package dao

import "wager_tagger_go_api/db"

type Tag struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

func GetTags() []Tag {
	var tags []Tag
	daoDb := db.GetDb()
	defer daoDb.Close()

	daoDb.Find(&tags)

	return tags
}
