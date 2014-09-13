package main

import (
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"wager_tagger_go_api/dao"
	"wager_tagger_go_api/handlers"
	"wager_tagger_go_api/wflags"

	"github.com/ant0ine/go-json-rest/rest"
)

type Message struct {
	Body string
}

type MyAuthenticationMiddleware struct{}
type MyCorsMiddleware struct{}

func main() {
	wflags.ProcessFlags()

	handler := rest.ResourceHandler{
		PreRoutingMiddlewares: []rest.Middleware{
			&MyCorsMiddleware{},
			&MyAuthenticationMiddleware{},
		},
	}
	err := handler.SetRoutes(
		&rest.Route{"GET", "/", rootHandler},
		&rest.Route{"GET", "/tickets", handlers.GetTickets},
		&rest.Route{"GET", "/tickets/:ticket_id", handlers.GetTicket},
		&rest.Route{"POST", "/tickets/:ticket_id/ticket_tags", handlers.CreateTicketTag},
		&rest.Route{"DELETE", "/tickets/:ticket_id/ticket_tags/:ticket_tag_id", handlers.DeleteTicketTag},
		&rest.Route{"GET", "/tags", handlers.GetTags},
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

		log.Println("Cors Info IsCors", corsInfo.IsCors)
		log.Println("Cors Info IsPreflight", corsInfo.IsPreflight)
		log.Println("Cors Info AccessControlRequestMethod", corsInfo.AccessControlRequestMethod)
		log.Println("Cors Info AccessControlRequestHeaders", corsInfo.AccessControlRequestHeaders)
		log.Println("Cors Info Origin", corsInfo.Origin)
		log.Println("Cors Info OriginUrl", corsInfo.OriginUrl)

		if !corsInfo.IsCors {
			handler(writer, request)
			return
		}

		if corsInfo.IsPreflight {
			// check the request methods
			allowedMethods := map[string]bool{
				"GET":    true,
				"POST":   true,
				"PUT":    true,
				"DELETE": true,
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
				"Origin":          true,
			}
			for _, requestedHeader := range corsInfo.AccessControlRequestHeaders {
				if !allowedHeaders[requestedHeader] {
					log.Println("Invalid Preflight Request")
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
		log.Println("Authenticating...")
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
		log.Println("Missing the Authorization Request Header")
		return false
	}
	re := regexp.MustCompile(`"(.*?)"`)
	match := re.FindString(authHeader)
	accessToken := strings.Trim(match, "\"")
	log.Println("Access Token", accessToken)
	if match == "" {
		log.Println("Missing Access Token")
		return false
	}
	return dao.ValidAccessToken(accessToken)
}

func rootHandler(w rest.ResponseWriter, req *rest.Request) {
	responseMap := map[string]string{"result": "Welcome to the jungle"}
	w.WriteJson(responseMap)
}
