package handlers

import (
	"fmt"
	"strconv"
	"wager_tagger_go_api/dao"

	"github.com/ant0ine/go-json-rest/rest"
)

func GetTickets(w rest.ResponseWriter, req *rest.Request) {
	queryParams := req.URL.Query()
	fmt.Println("Query Params: ", queryParams)

	tickets := dao.GetTickets(queryParams)

	w.WriteJson(tickets)
}

func GetTicket(w rest.ResponseWriter, req *rest.Request) {
	ticketId, err := strconv.Atoi(req.PathParam("ticket_id"))
	if err != nil {
		rest.Error(w, "invalid ticket id", 400)
	}

	ticket := dao.GetTicket(ticketId)
	w.WriteJson(ticket)
}
