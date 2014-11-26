package dao

import (
	"log"
	"wager_tagger_go_api/db"
)

type Finance struct {
	TagId   int64   `json:"tag_id"`
	Won     float32 `json:"won"`     // outcome is Won
	Lost    float32 `json:"lost"`    // outcome is Lost
	Pending float32 `json:"pending"` // outcome is Pending
	Total   float32 `json:"total"`   // all outcomes
}

type dbResponse struct {
	Id       int64
	TicketId int64
	TagId    int64
	Amount   *float32
	Outcome  string
}

func GetFinances(tagId int64, startDate string, stopDate string) Finance {
	var responses []dbResponse
	daoDb := db.GetDb()
	defer daoDb.Close()

	daoDb.Table("ticket_tags").
		Select("ticket_tags.*, tickets.outcome").
		Joins("left join tickets on ticket_tags.ticket_id = tickets.id").
		Where("ticket_tags.tag_id = ? and tickets.wager_date > ? and tickets.wager_date < ?", tagId, startDate, stopDate).
		Scan(&responses)

	log.Println(responses)

	finance := Finance{TagId: tagId}
	for _, response := range responses {
		switch response.Outcome {
		case "Won":
			finance.Won += *response.Amount
		case "Lost":
			finance.Lost += *response.Amount
		case "Pending":
			finance.Pending += *response.Amount
		}
		finance.Total += *response.Amount
	}

	return finance
}

// select tt.*, t.wager_date, t.outcome
// from ticket_tags tt
// join tickets t on t.id = tt.ticket_id
// where t.wager_date > '2014-07-01' and t.wager_date < '2015-06-30'
// order by t.wager_date desc;
