# Wager Tagger API
An api to serve ticket information. Hosted on [Heroku](https://wager-tagger-go-api.herokuapp.com/)

## Wager Tagger Application
The entire wager-tagger application consists of three parts:

1. UI - [source](https://github.com/hoitomt/wager-tagger-ui)
2. API - [source](https://github.com/hoitomt/wager-tagger-go-api)
3. Scraper - [source](https://github.com/hoitomt/wager-tagger-scraper)

##Usage
An API key is required

##Endpoints
The following endpoints are available

### GET /tickets[limit=][offset=][start_date=][stop_date=]
Return all of the tickets for the given query params

  - limit: the number of tickets that will be returned
  - offset: the starting point at which to query for tickets
  - start_date: the date at which every ticket (wager_date) will be after
  - stop_date: the date at which every ticket (wager_date) will be before

### GET /tickets/:ticket_id
Get the ticket based on the specified :ticket_id

### POST /tickets/:ticket_id/ticket_tags
Add a new tag to an existing ticket

### DELETE /ticket/:ticket_id/ticket_tags/:ticket_tag_id
Delete a ticket_tag based on the specified :ticket_tag_id

### GET /finances
Returns the balances for each user

### GET /tags
Returns a list of tags with which the tickets can be tagged

