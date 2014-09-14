package dao

import (
	"fmt"
	"wager_tagger_go_api/db"
)

type TicketTag struct {
	Id       int64   `json:"id"`
	TicketId int64   `json:"ticket_id"`
	TagId    int     `json:"tag_id"`
	Amount   float64 `json:"amount"`
	Name     string  `json:"name"`
}

func (t *TicketTag) Create() {
	fmt.Println("Ticket Tag Create: ", t)
	daoDb := db.GetDb()
	defer daoDb.Close()

	query := daoDb.Raw("INSERT INTO ticket_tags (ticket_id, tag_id, amount) VALUES (?, ?, ?) RETURNING id", &t.TicketId, &t.TagId, &t.Amount)
	row := query.Row()
	row.Scan(&t.Id)
	// db.Table("ticket_tags").Insert
	// return t
}

func DeleteTicketTag(ticketTagId int64) {
	daoDb := db.GetDb()
	defer daoDb.Close()

	daoDb.Where("id = ?", ticketTagId).Delete(TicketTag{})
}

func GetTicketTags(ticketIds []int64) map[int64][]TicketTag {
	var ticketTags []TicketTag
	ticketTagsMap := make(map[int64][]TicketTag)

	daoDb := db.GetDb()
	defer daoDb.Close()

	daoDb.Table("ticket_tags").
		Select("tags.name, ticket_tags.*").
		Joins("left join tags on ticket_tags.tag_id = tags.id").
		Scan(&ticketTags)

	for i := range ticketTags {
		ticketTag := ticketTags[i]

		ticketTagArray := ticketTagsMap[ticketTag.TicketId]
		ticketTagsMap[ticketTag.TicketId] = append(ticketTagArray, ticketTag)
	}

	return ticketTagsMap
}
