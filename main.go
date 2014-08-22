package main

import (
  // "net"
  "net/http"
  "log"
  "fmt"
  "regexp"
  "strings"
  // "encoding/json"
  "github.com/jinzhu/gorm"
  _ "github.com/lib/pq"
  "github.com/ant0ine/go-json-rest/rest"
  // "database/sql"
  "time"
)

type Ticket struct {
  Id int64
  SbBetId int64
  WagerDate time.Time
  WagerType string
  AmountWagered *float32
  AmountToWin *float32
  Outcome string
  TicketLineItems []TicketLineItem
}

type TicketLineItem struct {
  Id int64
  TicketId int64
  AwayTeam string
  AwayScore *int16
  HomeTeam string
  HomeScore *int16
  LineItemDate time.Time
  LineItemSpread string
}

type Message struct {
  Body string
}

type AuthenticationMiddleware struct{}

func main() {
  handler := rest.ResourceHandler{
    PreRoutingMiddlewares: []rest.Middleware{
      &AuthenticationMiddleware{},
    },
  }
  err := handler.SetRoutes(
    &rest.Route{"GET", "/", rootHandler},
    &rest.Route{"GET", "/tickets", ticketHandler},
  )
  log.Println("listening...")
  if err != nil {
      log.Fatal(err)
  }
  log.Fatal(http.ListenAndServe(":4000", &handler))
}

func (mw *AuthenticationMiddleware) MiddlewareFunc(handler rest.HandlerFunc) rest.HandlerFunc {
  return func(writer rest.ResponseWriter, request *rest.Request) {
    fmt.Println("Authenticating...")

    if authenticatedRequest(request) {
      handler(writer, request)
    } else {
      mw.unauthorized(writer)
    }
    return
  }
}

func (mw *AuthenticationMiddleware) unauthorized(writer rest.ResponseWriter) {
  rest.Error(writer, "Not Authorized", http.StatusUnauthorized)
}

func authenticatedRequest(request *rest.Request) bool {
  authHeader := request.Header.Get("Authorization")
  if authHeader == "" {
    return false
  }
  re := regexp.MustCompile(`"(.*?)"`)
  match := re.FindString(authHeader)
  access_token := strings.Trim(match, "\"")

  if match == "" {
    return false
  }

  db := getDb()
  defer db.Close()

  var count int8
  db.Table("api_keys").Select("access_token").Where("access_token = ?", access_token).Count(&count)
  return count > 0
}

func rootHandler(w rest.ResponseWriter, req *rest.Request) {
  responseMap := map[string]string{"result": "Welcome to the jungle"}
  w.WriteJson(responseMap)
}

func ticketHandler(w rest.ResponseWriter, req *rest.Request) {
  db := getDb()
  defer db.Close()

  tickets := ticketsWithTicketLineItems(db)
  w.WriteJson(tickets)
}

func ticketsWithTicketLineItems(db gorm.DB) []Ticket {
  ticketLineItemsMap := getTicketLineItems(db)
  tickets := getTickets(db)
  for i := range tickets {
    tickets[i].TicketLineItems = ticketLineItemsMap[tickets[i].Id]
  }
  return tickets
}

func getTickets(db gorm.DB) []Ticket {
  var tickets []Ticket

  db.Find(&tickets)
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
