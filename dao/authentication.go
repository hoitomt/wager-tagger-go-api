package dao

import "wager_tagger_go_api/db"

func ValidAccessToken(accessToken string) bool {
	queryDb := db.GetDb()
	defer queryDb.Close()

	var count int8
	queryDb.Table("api_keys").Select("access_token").Where("access_token = ?", accessToken).Count(&count)
	return count > 0
}
