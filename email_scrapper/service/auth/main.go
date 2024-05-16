package auth

import (
	"fmt"
	"io"
	"log"
	"os"

	"email_scrapper/pkg/imap_auth"
	"html/template"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

const AUTH_PORT int = 8080

type Server struct {
	sessions []Session
}

type Session struct {
	// db connection
	// context stored in memory: object Id
}

func (server *Server) handleAuthRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Println("this is auth root path")
	fmt.Println(r.URL.Path)

	t, err := template.ParseFiles("root.html")
	if err != nil {
		log.Println("template not parsed %s", err)
	}

	// data written to Response!
	t.Execute(w, nil)
}

var mail_mapping = map[string]int{
	"gmail":   1,
	"outlook": 2,
}

func (server *Server) handleUserInfo(w http.ResponseWriter, r *http.Request) {
	// activate the db
	var client *mongo.Client
	var coll *mongo.Collection
	db, err := Connect(client)
	defer Disconnect(client)
	if err != nil {
		log.Fatal(err)
	}
	coll = useDoc(db)

	r.ParseForm()
	var user_email string = r.FormValue("email")
	var user_provider string = r.FormValue("provider")

	_, err = db_saveUser(coll, user_email, mail_mapping[user_provider])

	jumpToOauthPage()
}

func jumpToOauthPage() {
	// redirect website to google oauth2.0 page

}

func (server *Server) handleAuthRedirect(w http.ResponseWriter, r *http.Request) {
	// receive a request from mail service provider
	fmt.Println("either success or fail after wait")
	t, err := template.ParseFiles("gmail_auth_done.html")
	if err != nil {
		log.Println("template not parsed %s", err)
	}
	t.Execute(w, nil)
	// for gmail, parse token from url
	vcode := getGmailVerificationCode()
	fmt.Println(vcode)

	// send vcode to gmail server

	// save token to local nosql db
}

func getGmailVerificationCode() string {
	// parse request and get the verification code
	return ""
}

func saveToDB() error {

}

func main() {
	server := Server{
		sessions: make([]Session, 0),
	}

	http.HandleFunc("/auth", server.handleAuthRoot)
	http.HandleFunc("/auth/userinfo", server.handleUserInfo)
	http.HandleFunc("/auth/redirect", server.handleAuthRedirect)

	err := http.ListenAndServe(fmt.Sprintf(":%v", AUTH_PORT), nil)
	log.Fatalf("server error %s", err)
}
