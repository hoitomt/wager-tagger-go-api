package dao

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
	"time"
	"wager_tagger_go_api/db"

	"github.com/jinzhu/gorm"
)

type Ticket struct {
	Id              int64            `json:"id"`
	SbBetId         int64            `json:"sb_bet_id"`
	WagerDate       time.Time        `json:"wager_date"`
	WagerType       string           `json:"wager_type"`
	AmountWagered   *float32         `json:"amount_wagered"`
	AmountToWin     *float32         `json:"amount_to_win"`
	Outcome         string           `json:"outcome"`
	TicketLineItems []TicketLineItem `json:"ticket_line_items"`
	TicketTags      []TicketTag      `json:"ticket_tags"`
}

func GetTicket(ticketId int) Ticket {
	var ticket Ticket
	var ticketLineItems []TicketLineItem
	// var ticketTags []TicketTag
	ticketTags := []TicketTag{}

	daoDb := db.GetDb()
	defer daoDb.Close()

	if err := daoDb.Where("id = ?", ticketId).First(&ticket).Error; err != nil {
		log.Printf("Invalid Query for Tickets: %v", err)
	}

	if err := daoDb.Where("ticket_id = ?", ticketId).Find(&ticketLineItems).Error; err != nil {
		log.Printf("Invalid Query for Ticket Line Items: %v", err)
	}

	daoDb.Table("ticket_tags").
		Select("tags.name, ticket_tags.*").
		Joins("left join tags on ticket_tags.tag_id = tags.id").
		Where("ticket_id = ?", ticketId).
		Scan(&ticketTags)

	ticket.TicketLineItems = ticketLineItems
	ticket.TicketTags = ticketTags

	return ticket
}

func GetTickets(queryParams url.Values) []Ticket {
	var tickets []Ticket
	var ticketIds []int64
	daoDb := db.GetDb()
	defer daoDb.Close()

	ticketQuery := daoDb.Order("wager_date desc")
	ticketQuery = setLimit(queryParams, ticketQuery)
	ticketQuery = setRange(queryParams, ticketQuery)

	ticketQuery.Find(&tickets)

	for i := range tickets {
		ticketIds = append(ticketIds, tickets[i].Id)
	}

	// Maps of line items and tags
	ticketLineItems := GetTicketLineItems(ticketIds)
	ticketTags := GetTicketTags(ticketIds)

	for j := range tickets {
		ticket := &tickets[j]

		// Assign empty arrays if no items
		if len(ticketLineItems[ticket.Id]) == 0 {
			ticket.TicketLineItems = []TicketLineItem{}
		} else {
			ticket.TicketLineItems = ticketLineItems[ticket.Id]
		}

		if len(ticketTags[ticket.Id]) == 0 {
			ticket.TicketTags = []TicketTag{}
		} else {
			ticket.TicketTags = ticketTags[ticket.Id]
		}
	}
	return tickets
}

func setLimit(queryParams url.Values, ticketQuery *gorm.DB) *gorm.DB {
	qLimit := queryParams["limit"]
	if len(qLimit) > 0 {
		limit, _ := strconv.Atoi(qLimit[0])

		var qPage []string

		qPage = queryParams["page"]
		if len(qPage) == 0 {
			qPage = []string{"1"}
		}
		page, _ := strconv.Atoi(qPage[0])
		offset := (page - 1) * limit

		fmt.Println("Page", page, "Offset", offset)

		ticketQuery = ticketQuery.Limit(limit).Offset(offset)
	}
	return ticketQuery
}

func setRange(queryParams url.Values, ticketQuery *gorm.DB) *gorm.DB {
	qStartDate := queryParams["start_date"]
	qStopDate := queryParams["stop_date"]
	if len(qStartDate) > 0 {
		startDate := qStartDate[0]
		if len(qStopDate) > 0 {
			stopDate := qStopDate[0]
			ticketQuery = ticketQuery.Where("wager_date BETWEEN ? AND ?", startDate, stopDate)
		} else {
			ticketQuery = ticketQuery.Where("wager_date >= ?", startDate)
		}
	}
	return ticketQuery
}
