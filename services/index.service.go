package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"plugin"

	"github.com/Tushar3099/CollectPlugin/database"
	"github.com/Tushar3099/CollectPlugin/models"
	"github.com/Tushar3099/CollectPlugin/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func PostUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var user models.User
	err := decoder.Decode(&user)
	if err != nil {
		utils.WriteError(fmt.Errorf("invalid fields"), w)
		return
	}

	col := database.Client.Database("Go-test").Collection("user")

	result, err := col.InsertOne(context.TODO(), user)
	if err != nil {
		utils.WriteError(err, w)
		return
	}
	utils.WriteResponse(result, w)
}

func GetForm(w http.ResponseWriter, r *http.Request) {
	col := database.Client.Database("Go-test").Collection("form")
	cursor, err := col.Find(context.TODO(), bson.D{})
	if err != nil {
		utils.WriteError(err, w)
		return
	}
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		utils.WriteError(err, w)
		return
	}
	if results == nil {
		results = []bson.M{}
	}
	utils.WriteResponse(results, w)
}

func PutForm(w http.ResponseWriter, r *http.Request) {

}

func PostForm(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var form models.Form
	err := decoder.Decode(&form)
	if err != nil {
		utils.WriteError(fmt.Errorf("invalid fields"), w)
	}

	// Checks if action is present in form
	if form.ActionName != "" {
		// Checks if the action given is supported
		if !utils.ActionExists(form.ActionName) {
			utils.WriteError(fmt.Errorf("action does not exists"), w)
			return
		}

		// Gets the path of plugin with the plugin-name
		path, err := os.Getwd()
		if err != nil {
			log.Println(err)
			return
		}
		path = filepath.Join(path, "/actions/"+form.ActionName+"/action.so")

		// Loads the plugin
		plug, err := plugin.Open(path)
		if err != nil {
			log.Fatalf("Plugin Error : %s", err.Error())
			os.Exit(1)
		}

		// Looks up for symbol (in this case variable)
		actionPlugin, err := plug.Lookup("Action")
		if err != nil {
			fmt.Printf("Plugin Error : %s", err.Error())
			os.Exit(1)
		}
		var ok bool
		form.Action, ok = actionPlugin.(models.Action)

		if !ok {
			fmt.Println("unexpected type from module symbol")
			os.Exit(1)
		}

		// Initialize the plugin and populate the action-meta data in form
		err = form.Action.Initialize(&form)

		if err != nil {
			utils.WriteError(err, w)
			return
		}
	}

	col := database.Client.Database("Go-test").Collection("form")

	result, err := col.InsertOne(context.TODO(), form)
	if err != nil {
		utils.WriteError(err, w)
	}
	utils.WriteResponse(result, w)
}
