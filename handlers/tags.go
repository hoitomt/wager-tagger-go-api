package handlers

import (
	"wager_tagger_go_api/dao"

	"github.com/ant0ine/go-json-rest/rest"
)

func GetTags(w rest.ResponseWriter, req *rest.Request) {
	tags := dao.GetTags()
	w.WriteJson(tags)
}
