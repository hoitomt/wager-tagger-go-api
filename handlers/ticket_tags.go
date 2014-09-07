package handlers

import (
	"fmt"
	"strconv"
	"wager_tagger_go_api/dao"

	"github.com/ant0ine/go-json-rest/rest"
)

func CreateTicketTag(w rest.ResponseWriter, req *rest.Request) {
	fmt.Println("Request Body", req.Body)
	ticketIdStr := req.PathParam("ticket_id")

	ticketId, err := strconv.ParseInt(ticketIdStr, 10, 64)
	if err != nil {
		rest.Error(w, "invalid ticket id", 400)
		return
	}

	ticketTag := dao.TicketTag{}
	req.DecodeJsonPayload(&ticketTag)

	fmt.Printf("Amount %v\n", ticketTag.Amount)

	fmt.Println("Json Ticket Tag", ticketTag)

	ticketTag.TicketId = ticketId

	fmt.Println("Json Ticket Tag Form", ticketTag)
	ticketTag.Create()

	w.WriteJson(ticketTag)
}

func DeleteTicketTag(w rest.ResponseWriter, req *rest.Request) {
	ticketTagIdStr := req.PathParam("ticket_tag_id")

	ticketTagId, err := strconv.ParseInt(ticketTagIdStr, 10, 64)
	if err != nil {
		rest.Error(w, "invalid ticket id", 400)
		return
	}

	dao.DeleteTicketTag(ticketTagId)

	message := map[string]string{
		"message": "success",
	}
	w.WriteJson(message)
}
