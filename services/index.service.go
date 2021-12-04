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
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	decoder := json.NewDecoder(r.Body)
	var res struct {
		Form_id primitive.ObjectID `json:"form_id"`
		Ans     models.Answer      `json:"answer"`
	}

	err := decoder.Decode(&res)
	if err != nil {
		utils.WriteError(fmt.Errorf("invalid fields"), w)
		return
	}

	// Fetching form with id of form_id
	var form models.Form
	var result bson.D
	col := database.Client.Database("Go-test").Collection("form")
	err = col.FindOne(context.TODO(), bson.D{{"_id", res.Form_id}}).Decode(&result)

	if err != nil {
		utils.WriteError(err, w)
		return
	}
	bsonBytes, _ := bson.Marshal(result)
	bson.Unmarshal(bsonBytes, &form)

	// Validating the answers length
	if len(form.Questions) != len(res.Ans.List) {
		utils.WriteError(fmt.Errorf("answers length is not equal to questions length"), w)
	}

	// Checks if action is present in form
	// Actions are executed once the form gets updates
	if form.ActionName != "" {
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
		constructor, err := plug.Lookup("Constructor")
		if err != nil {
			fmt.Printf("Plugin Error : %s", err.Error())
			os.Exit(1)
		}
		f := constructor.(func(b *[]byte) *models.Form)(&bsonBytes)

		// Executes the Plugin action
		err = f.Action.Execute(f, &res.Ans)

		if err != nil {
			utils.WriteError(err, w)
			return
		}
	}

	// Pushes answer to the forms answers array
	// and updates form
	change := bson.M{"$push": bson.M{"answers": res.Ans}}
	opts := options.Update().SetUpsert(true)
	_, err = col.UpdateByID(context.TODO(), res.Form_id, change, opts)
	if err != nil {
		utils.WriteError(err, w)
	}
	utils.WriteResponse("Succefully Answered", w)
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
	form.Answers = []models.Answer{}
	result, err := col.InsertOne(context.TODO(), form)
	if err != nil {
		utils.WriteError(err, w)
	}
	utils.WriteResponse(result, w)
}
