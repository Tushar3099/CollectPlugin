package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Tushar3099/CollectPlugin/models"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type Structure struct {
	Link    string `json:"link" bson:"link"`
	SheetId string `json:"form_id" bson:"form_id"`
}

type NewForm struct {
	Action Structure `json:"action" bson:"action"`
}

func (g *Structure) Constructor(b *[]byte) *models.Form {
	var newform NewForm
	var form models.Form
	bson.Unmarshal(*b, &form)
	bson.Unmarshal(*b, &newform)
	form.Action = &newform.Action
	return &form
}

func (g *Structure) Execute(form *models.Form, ans *models.Answer) error {
	fmt.Println(g.Link)
	// fmt.Println("Hello Universe")
	return nil
}

func (g *Structure) Initialize(form *models.Form) error {

	err := g.validate()
	if err != nil {
		return err
	}

	ctx := context.Background()
	path, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("unable to get dir path in plugin: %v", err)
	}
	path = filepath.Join(path, "/actions/"+form.ActionName+"/credentials.json")

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		return fmt.Errorf("unable to parse client secret file to config: %v", err)
	}
	client, err := getClient(config, form.ActionName, "sheet")
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	// Retrieving Sheets Client
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return fmt.Errorf("unable to retrieve Sheets client: %v", err)
	}

	var sheet sheets.Spreadsheet
	sheet.Properties = &sheets.SpreadsheetProperties{
		Title: form.Title,
	}

	// Creating new spreadsheet
	newSpreadSheet, err := srv.Spreadsheets.Create(&sheet).Context(context.TODO()).Do()

	if err != nil {
		return fmt.Errorf("unable to create a new spreadsheet: %v", err)
	}

	g.Link = newSpreadSheet.SpreadsheetUrl
	g.SheetId = newSpreadSheet.SpreadsheetId

	spreadsheetId := newSpreadSheet.SpreadsheetId

	fmt.Println(newSpreadSheet.SpreadsheetUrl)

	writeRange := "A1"
	var vr sheets.ValueRange
	var myval []interface{}
	myval = append(myval, "User_Id")

	for i := range form.Questions {
		myval = append(myval, "Question - "+fmt.Sprintf("%v", i+1))
	}

	vr.Values = append(vr.Values, myval)
	// Updating the spread sheet
	_, err = srv.Spreadsheets.Values.Update(spreadsheetId, writeRange, &vr).ValueInputOption("RAW").Do()

	if err != nil {
		return fmt.Errorf("unable to retrieve data from sheet. %v", err)
	}

	writeRange = "A2"
	vr.Values = nil
	myval = nil
	myval = append(myval, "----")

	for _, q := range form.Questions {
		myval = append(myval, q.Title)
	}

	vr.Values = append(vr.Values, myval)
	_, err = srv.Spreadsheets.Values.Update(spreadsheetId, writeRange, &vr).ValueInputOption("RAW").Do()

	if err != nil {
		return fmt.Errorf("unable to retrieve data from sheet. %v", err)
	}

	// The created spreadsheet will be private
	// To make the spreadsheet public we have to
	// used google drive api. Retrieve drive api Client
	// and the create a permission for the spreadsheet
	// with the permission "anyone"
	gconfig, err := google.ConfigFromJSON(b, drive.DriveScope)
	if err != nil {
		return fmt.Errorf("unable to parse client secret file to config: %v", err)
	}

	gclient, _ := getClient(gconfig, form.ActionName, "drive")
	gsrv, err := drive.NewService(ctx, option.WithHTTPClient(gclient))
	if err != nil {
		return fmt.Errorf("unable to retrieve Drive client: %v", err)
	}

	p := &drive.Permission{
		Type: "anyone",
		Role: "reader",
	}
	_, err = gsrv.Permissions.Create(newSpreadSheet.SpreadsheetId, p).Do()
	if err != nil {

		fmt.Printf("unable to make sheet public : %v", err.Error())
	}

	return nil
}

func (g *Structure) validate() error {
	//validates if the action requires specific fields in forms
	return nil
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config, action_name string, api string) (*http.Client, error) {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	path, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("unable to get dir path in plugin: %v", err)
	}
	if api == "sheet" {
		path = filepath.Join(path, "/actions/"+action_name+"/token.json")
	} else if api == "drive" {
		path = filepath.Join(path, "/actions/"+action_name+"/gtoken.json")
	}

	tok, err := tokenFromFile(path)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(path, tok)
	}
	return config.Client(context.Background(), tok), nil
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

// exported as symbol named "Action"
var Action Structure
