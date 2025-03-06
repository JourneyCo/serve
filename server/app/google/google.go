package google

import (
	"github.com/kelvins/geocoder"
	"serve/helpers"
)

const (
// testGoogleSheetID = "1ZVkmOMblClYaiUsbEzfur-fsHxSus3l-bxwq-4qeKTU"
// readRange         = "Serve Day!A4:S"
// journeySheet      = "1QvHP4eax2ve4UIMcVuaTCGuGcOighVhgHm8RNE0qzis"
)

func SetKey() {
	key := helpers.GetEnvVar("GOOGLE_KEY")
	geocoder.ApiKey = key
}

// // getTokenFromWeb Requests a token from the web, then returns the retrieved token.
// func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
// 	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
// 	fmt.Printf(
// 		"Go to the following link in your browser then type the "+
// 			"authorization code: \n%v\n", authURL,
// 	)
//
// 	var authCode string
// 	if _, err := fmt.Scan(&authCode); err != nil {
// 		log.Printf("Unable to read authorization code: %v", err)
// 		return nil
// 	}
//
// 	tok, err := config.Exchange(context.TODO(), authCode)
// 	if err != nil {
// 		log.Printf("Unable to retrieve token from web: %v", err)
// 		return nil
// 	}
// 	return tok
// }

// // tokenFromFile Retrieves a token from a local file.
// func tokenFromFile(file string) (*oauth2.Token, error) {
// 	f, err := os.Open(file)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer f.Close()
// 	tok := &oauth2.Token{}
// 	err = json.NewDecoder(f).Decode(tok)
// 	return tok, err
// }
//
// // saveToken Saves a token to a file path.
// func saveToken(path string, token *oauth2.Token) {
// 	fmt.Printf("Saving credential file to: %s\n", path)
// 	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
// 	if err != nil {
// 		log.Printf("Unable to cache oauth token: %v", err)
// 		return
// 	}
// 	defer f.Close()
// 	json.NewEncoder(f).Encode(token)
// }

// // FetchProjects returns the projects that have been entered into a Google Sheet.
// func FetchProjects() error {
// 	ctx := context.Background()
// 	b, err := os.ReadFile("NOT_USED_GOOGLESHEET_oauth_google_creds.json")
// 	if err != nil {
// 		log.Printf("Unable to read client secret file: %v", err)
// 		return err
// 	}
//
// 	// If modifying these scopes, delete your previously saved NOTUSEDgoogle_sheet_token.json.
// 	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets.readonly")
// 	if err != nil {
// 		log.Printf("Unable to parse client secret file to config: %v", err)
// 		return err
// 	}
//
// 	client := getClient(config)
//
// 	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
// 	if err != nil {
// 		log.Printf("Unable to retrieve Sheets client: %v", err)
// 		return err
// 	}
//
// 	resp, err := srv.Spreadsheets.Values.Get(journeySheet, readRange).Do()
// 	if err != nil {
// 		log.Printf("Unable to retrieve data from sheet: %v", err)
// 		return err
// 	}
//
// 	if len(resp.Values) == 0 {
// 		fmt.Println("No projects found.")
// 		err = errors.New("no projects found")
// 		return err
// 	}
//
// 	for _, row := range resp.Values {
// 		for _, cell := range row {
// 			log.Printf("%v ", cell)
// 		}
// 		print("\n")
// 	}
//
// 	return nil
// }

// getClient Retrieve a token, saves the token, then returns the generated client.
// func getClient(config *oauth2.Config) *http.Client {
// 	// The file NOTUSEDgoogle_sheet_token.json.json stores the user's access and refresh tokens, and is
// 	// created automatically when the authorization flow completes for the first
// 	// time.
// 	tokFile := "NOTUSEDgoogle_sheet_token.json"
// 	tok, err := tokenFromFile(tokFile)
// 	if err != nil {
// 		tok = getTokenFromWeb(config)
// 		saveToken(tokFile, tok)
// 	}
// 	return config.Client(context.Background(), tok)
// }
