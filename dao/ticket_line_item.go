package dao

import (
	"time"
	"wager_tagger_go_api/db"
)

type TicketLineItem struct {
	Id             int64      `json:"id"`
	TicketId       int64      `json:"ticket_id"`
	AwayTeam       string     `json:"away_team"`
	AwayScore      *int16     `json:"away_score"`
	HomeTeam       string     `json:"home_team"`
	HomeScore      *int16     `json:"home_score"`
	LineItemDate   *time.Time `json:"line_item_date"`
	LineItemSpread string     `json:"line_item_spread"`
}

func GetTicketLineItems(ticketIds []int64) map[int64][]TicketLineItem {
	var ticketLineItems []TicketLineItem
	ticketLineItemsMap := make(map[int64][]TicketLineItem)

	daoDb := db.GetDb()
	defer daoDb.Close()

	daoDb.Where("ticket_id in (?)", ticketIds).Find(&ticketLineItems)

	for i := range ticketLineItems {
		ticketLineItem := ticketLineItems[i]

		ticketLineItemArray := ticketLineItemsMap[ticketLineItem.TicketId]
		ticketLineItemsMap[ticketLineItem.TicketId] = append(ticketLineItemArray, ticketLineItem)
	}

	return ticketLineItemsMap
}
