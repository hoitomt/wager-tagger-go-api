package dao

import "wager_tagger_go_api/db"

type Finance struct {
	TagId   int64   `json:"tag_id"`
	Name    string  `json:"name"`
	Won     float32 `json:"won"`     // outcome is Won
	Lost    float32 `json:"lost"`    // outcome is Lost
	Pending float32 `json:"pending"` // outcome is Pending
	Total   float32 `json:"total"`   // all outcomes
}

type dbResponse struct {
	Id            int64
	TicketId      int64
	TagId         int64
	Amount        *float32
	AmountWagered *float32
	Name          string
	AmountToWin   *float32
	Outcome       string
}

func GetFinances(tagId int64, startDate string, stopDate string) Finance {
	var responses []dbResponse
	daoDb := db.GetDb()
	defer daoDb.Close()

	daoDb.Table("ticket_tags").
		Select("ticket_tags.*, tickets.amount_wagered, tags.name, tickets.amount_to_win, tickets.outcome").
		Joins("left join tickets on ticket_tags.ticket_id = tickets.id left join tags on ticket_tags.tag_id = tags.id").
		Where("ticket_tags.tag_id = ? and tickets.wager_date > ? and tickets.wager_date < ?", tagId, startDate, stopDate).
		Scan(&responses)

	finance := Finance{TagId: tagId}
	for _, response := range responses {
		amtWageredPerTag := *response.Amount
		amtWageredPerTicket := *response.AmountWagered
		pctStake := amtWageredPerTag / amtWageredPerTicket
		finance.Name = response.Name

		switch response.Outcome {
		case "Won":
			finance.Won += (*response.AmountToWin * pctStake) + *response.Amount
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
