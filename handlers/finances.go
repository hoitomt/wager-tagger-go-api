package handlers

import (
	"log"
	"net/http"
	"wager_tagger_go_api/dao"

	"github.com/ant0ine/go-json-rest/rest"
)

func GetFinances(w rest.ResponseWriter, req *rest.Request) {
	queryParams := req.URL.Query()
	var (
		startDate string
		stopDate  string
	)

	if len(queryParams["start_date"]) > 0 {
		startDate = queryParams["start_date"][0]
	} else {
		rest.Error(w, "start_date is missing", http.StatusBadRequest)
		return
	}

	if len(queryParams["stop_date"]) > 0 {
		stopDate = queryParams["stop_date"][0]
	} else {
		rest.Error(w, "stop_date is missing", http.StatusBadRequest)
		return
	}

	tags := dao.GetTags()

	finances := []dao.Finance{}

	for _, tag := range tags {
		finance := dao.GetFinances(tag.Id, startDate, stopDate)
		log.Println("Finance", finance)
		finances = append(finances, finance)
	}

	w.WriteJson(finances)
}
