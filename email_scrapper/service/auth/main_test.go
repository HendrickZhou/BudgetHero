package auth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

// Retrieves a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok) // todo save it to database
	}
	return config.Client(context.Background(), tok)
}

// Requests a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.Background(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
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

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	defer f.Close()
	if err != nil {
		log.Fatalf("Unable to cache OAuth token: %v", err)
	}
	json.NewEncoder(f).Encode(token)
}

// refresh token
func refreshToken(path string, token *oauth2.Token) {
	saveToken(path, token)
}

func main_test() {
	ctx := context.Background()
	b, err := os.ReadFile("/Users/zhouhang/Project/BudgetHero/email_scrapper/credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, gmail.MailGoogleComScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve GMail client: %v", err)
	}

	// Prints the title of the requested doc:
	// https://docs.google.com/document/d/195j9eDD3ccgjQRttHhJPymLJUCOUjs-jmwTrekvdjFE/edit
	// docId := "195j9eDD3ccgjQRttHhJPymLJUCOUjs-jmwTrekvdjFE"
	// doc, err := srv.Get(docId).Do()
	// if err != nil {
	// 	log.Fatalf("Unable to retrieve data from document: %v", err)
	// }
	// fmt.Printf("The title of the doc is: %s\n", doc.Title)

	user := "me"
	lis, err := srv.Users.Messages.List(user).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve labels: %v", err)
	}
	fmt.Println("Ids:")
	for i, l := range lis.Messages {
		if i > 3 {
			break
		}
		fmt.Printf("- %s\n", l.Id)
		msgResp, _ := srv.Users.Messages.Get(user, l.Id).Format("Raw").Do()
		if msgResp != nil {
			decodedData, err := base64.RawURLEncoding.DecodeString(msgResp.Raw)
			if err != nil {
				log.Println("error b64 decoding: ", err)
			}
			fmt.Printf("- %s\n", decodedData)
		}

	}

}
