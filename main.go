package main

import (
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
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
}

type TicketLineItem struct {
	Id             int64     `json:"id"`
	TicketId       int64     `json:"ticket_id"`
	AwayTeam       string    `json:"away_team"`
	AwayScore      *int16    `json:"away_score"`
	HomeTeam       string    `json:"home_team"`
	HomeScore      *int16    `json:"home_score"`
	LineItemDate   time.Time `json:"line_item_date"`
	LineItemSpread string    `json:"line_item_spread"`
}

type Message struct {
	Body string
}

type MyAuthenticationMiddleware struct{}
type MyCorsMiddleware struct{}

func main() {
	handler := rest.ResourceHandler{
		PreRoutingMiddlewares: []rest.Middleware{
			&MyCorsMiddleware{},
			&MyAuthenticationMiddleware{},
		},
	}
	err := handler.SetRoutes(
		&rest.Route{"GET", "/", rootHandler},
		&rest.Route{"GET", "/tickets", GetTickets},
		&rest.Route{"GET", "/tickets/:id", GetTicketById},
	)
	log.Println("listening...")
	if err != nil {
		log.Fatal(err)
	}

	port := "4001"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	log.Fatal(http.ListenAndServe(":"+port, &handler))
}

func (mw *MyCorsMiddleware) MiddlewareFunc(handler rest.HandlerFunc) rest.HandlerFunc {
	return func(writer rest.ResponseWriter, request *rest.Request) {

		corsInfo := request.GetCorsInfo()

		if !corsInfo.IsCors {
			handler(writer, request)
			return
		}

		if corsInfo.IsPreflight {
			// check the request methods
			allowedMethods := map[string]bool{
				"GET":  true,
				"POST": true,
				"PUT":  true,
			}
			if !allowedMethods[corsInfo.AccessControlRequestMethod] {
				rest.Error(writer, "Invalid Preflight Request", http.StatusForbidden)
				return
			}
			// check the request headers
			allowedHeaders := map[string]bool{
				"Accept":          true,
				"Content-Type":    true,
				"X-Custom-Header": true,
				"Authorization":   true,
			}
			for _, requestedHeader := range corsInfo.AccessControlRequestHeaders {
				if !allowedHeaders[requestedHeader] {
					rest.Error(writer, "Invalid Preflight Request", http.StatusForbidden)
					return
				}
			}

			for allowedMethod, _ := range allowedMethods {
				writer.Header().Add("Access-Control-Allow-Methods", allowedMethod)
			}
			for allowedHeader, _ := range allowedHeaders {
				writer.Header().Add("Access-Control-Allow-Headers", allowedHeader)
			}
			writer.Header().Set("Access-Control-Allow-Origin", corsInfo.Origin)
			writer.Header().Set("Access-Control-Allow-Credentials", "true")
			writer.Header().Set("Access-Control-Max-Age", "3600")
			writer.WriteHeader(http.StatusOK)
			return
		} else {
			writer.Header().Set("Access-Control-Expose-Headers", "X-Powered-By")
			writer.Header().Set("Access-Control-Allow-Origin", corsInfo.Origin)
			writer.Header().Set("Access-Control-Allow-Credentials", "true")
			handler(writer, request)
			return
		}
	}
}

func (mw *MyAuthenticationMiddleware) MiddlewareFunc(handler rest.HandlerFunc) rest.HandlerFunc {
	return func(writer rest.ResponseWriter, request *rest.Request) {
		fmt.Println("Authenticating...")
		if authenticatedRequest(request) || request.URL.Path == "/" {
			handler(writer, request)
		} else {
			mw.unauthorized(writer)
		}
		return
	}
}

func (mw *MyAuthenticationMiddleware) unauthorized(writer rest.ResponseWriter) {
	rest.Error(writer, "Not Authorized", http.StatusUnauthorized)
}

func authenticatedRequest(request *rest.Request) bool {
	authHeader := request.Header.Get("Authorization")
	if authHeader == "" {
		return false
	}
	re := regexp.MustCompile(`"(.*?)"`)
	match := re.FindString(authHeader)
	accessToken := strings.Trim(match, "\"")
	fmt.Println("Access Token", accessToken)
	if match == "" {
		return false
	}

	db := getDb()
	defer db.Close()

	var count int8
	db.Table("api_keys").Select("access_token").Where("access_token = ?", accessToken).Count(&count)
	return count > 0
}

func rootHandler(w rest.ResponseWriter, req *rest.Request) {
	responseMap := map[string]string{"result": "Welcome to the jungle"}
	w.WriteJson(responseMap)
}

func GetTickets(w rest.ResponseWriter, req *rest.Request) {
	db := getDb()
	defer db.Close()

	queryParams := req.URL.Query()
	fmt.Println("Query Params: ", queryParams)

	tickets := ticketsWithTicketLineItems(db, queryParams)
	w.WriteJson(tickets)
}

func GetTicketById(w rest.ResponseWriter, req *rest.Request) {
	ticketId, err := strconv.Atoi(req.PathParam("id"))
	if err != nil {
		rest.Error(w, "invalid ticket id", 400)
	}

	db := getDb()
	defer db.Close()

	ticket := getCompleteTicket(db, ticketId)
	w.WriteJson(ticket)
}

func getCompleteTicket(db gorm.DB, ticketId int) Ticket {
	var ticket Ticket
	db.First(&ticket, ticketId)

	ticketLineItems := getTicketLineItemsByTicketId(db, ticketId)

	ticket.TicketLineItems = ticketLineItems

	return ticket
}

func ticketsWithTicketLineItems(db gorm.DB, queryParams url.Values) []Ticket {
	tickets := getTickets(db, queryParams)
	ticketLineItemsMap := getTicketLineItems(db)
	for i := range tickets {
		tickets[i].TicketLineItems = ticketLineItemsMap[tickets[i].Id]
	}
	return tickets
}

func getTickets(db gorm.DB, queryParams url.Values) []Ticket {
	var tickets []Ticket

	ticketQuery := db.Order("wager_date desc")

	qLimit := queryParams["limit"]
	qStartDate := queryParams["start_date"]
	qStopDate := queryParams["stop_date"]

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

	if len(qStartDate) > 0 {
		startDate := qStartDate[0]
		if len(qStopDate) > 0 {
			stopDate := qStopDate[0]
			ticketQuery = ticketQuery.Where("wager_date BETWEEN ? AND ?", startDate, stopDate)
		} else {
			ticketQuery = ticketQuery.Where("wager_date >= ?", startDate)
		}
	}

	ticketQuery.Find(&tickets)
	return tickets
}

func getTicketLineItems(db gorm.DB) map[int64][]TicketLineItem {
	var ticketLineItems []TicketLineItem
	tleMap := make(map[int64][]TicketLineItem)

	db.Find(&ticketLineItems)
	fmt.Println("Number of records: ", len(ticketLineItems))

	for i := range ticketLineItems {
		tle := ticketLineItems[i]
		s := append(tleMap[tle.TicketId], tle)
		tleMap[tle.TicketId] = s
	}
	return tleMap
}

func queryTickets(db gorm.DB, queryParams map[string]interface{}) []Ticket {
	var tickets []Ticket

	db.Order("wager_date desc").Find(&tickets)
	return tickets
}

func getTicketLineItemsByTicketId(db gorm.DB, ticketId int) []TicketLineItem {
	var ticketLineItems []TicketLineItem
	db.Where("ticket_id = ?", ticketId).Find(&ticketLineItems)
	return ticketLineItems
}

// Run `heroku pg:credentials DATABASE` from the ruby sportsbook_api for the connection string
func getDb() gorm.DB {
	db, err := gorm.Open("postgres",
		"dbname=d8ihclom31aprt host=ec2-54-197-246-197.compute-1.amazonaws.com port=5432 user=mffrovlynusepg password=uL6PkJjdATMqsSnvhwFzCIlh3X sslmode=require")
	if err != nil {
		log.Fatal(err)
	}
	db.LogMode(true)
	return db
}
