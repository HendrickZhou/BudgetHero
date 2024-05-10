package main

import (
	"fmt"
	"log"

	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
	"github.com/emersion/go-sasl"

	"encoding/json"
	"os"

	"golang.org/x/oauth2"
	// "golang.org/x/oauth2/google"
	// "google.golang.org/api/gmail/v1"
	// "google.golang.org/api/option"
)

const gmail_tls_server_address string = "imap.gmail.com"
const tls_port int = 993

// const username string = "zhouhangseu@gmail.com"
// const pswd string = "Harheihei@6@6"

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func main() {
	// a spinner that periodically send request to email server
	// save the data to no-sql db and send a message to analyzer worker

	// * find a way to get email service more reponsively and lightweight
	var client *imapclient.Client
	address := fmt.Sprintf("%s:%v", gmail_tls_server_address, tls_port)
	log.Println(address)

	// connect to server
	client, err := imapclient.DialTLS(address, nil)
	if err != nil {
		log.Print(err)
		log.Fatalf("fail to connect to imap server")
	}
	defer client.Close()

	// todo try auth
	tokenJson, err := tokenFromFile("/Users/zhouhang/Project/BudgetHero/email_scrapper/token.json")
	if err != nil {
		log.Fatalf("no token file found")
	}

	if !client.Caps().Has(imap.AuthCap(sasl.OAuthBearer)) {
		log.Fatal("OAUTHBEARER not supported by the server")
	}
	saslClient := sasl.NewOAuthBearerClient(&sasl.OAuthBearerOptions{
		Username: "hang_zhou@alumni.brown.edu",
		Token:    tokenJson.AccessToken,
	})
	if err := client.Authenticate(saslClient); err != nil {
		log.Fatalf("authentication failed: %v", err)
	}

	// login (blocking
	// err = client.Login(username, pswd).Wait()
	// if err != nil {
	// 	log.Print(err)
	// 	log.Fatalf("login fail, check usrname pswd")
	// }

	// select a mailbox(in/out/trash etc)
	mailboxes, err := client.List("", "", nil).Collect()
	if err != nil {
		log.Print(err)
		log.Println("fail to get mailbox list")
	}
	log.Printf("Found %v mailboxes", len(mailboxes))
	for _, mbox := range mailboxes {
		log.Print(err)
		log.Printf(" - %v", mbox.Mailbox)
	}

	selectedMbox, err := client.Select("INBOX", nil).Wait()
	if err != nil {
		log.Print(err)
		log.Fatalf("failed to select INBOX: %v", err)
	}
	log.Printf("INBOX contains %v messages", selectedMbox.NumMessages)

	// test run, fetch the first email in the box
	if selectedMbox.NumMessages > 0 {
		seqSet := imap.SeqSetNum(1)
		fetchOptions := &imap.FetchOptions{Envelope: true}
		messages, err := client.Fetch(seqSet, fetchOptions).Collect()
		if err != nil {
			log.Fatalf("failed to fetch first message in INBOX: %v", err)
		}
		log.Printf("subject of first message in INBOX: %v", messages[0].Envelope.Subject)
	}

	// idle loop to keep getting new email notificated
	// for {

	// }

	// clean up
	if err := client.Logout().Wait(); err != nil {
		log.Fatalf("failed to logout: %v", err)
	}

}
