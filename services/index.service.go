package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

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
	}

	col := database.Client.Database("Go-test").Collection("user")

	result, err := col.InsertOne(context.TODO(), user)
	if err != nil {
		utils.WriteError(err, w)
	}
	utils.WriteResponse(result, w)
}

func GetForm(w http.ResponseWriter, r *http.Request) {
	col := database.Client.Database("Go-test").Collection("form")
	cursor, err := col.Find(context.TODO(), bson.D{})
	if err != nil {
		utils.WriteError(err, w)
	}
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		utils.WriteError(err, w)
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

	col := database.Client.Database("Go-test").Collection("form")

	result, err := col.InsertOne(context.TODO(), form)
	if err != nil {
		utils.WriteError(err, w)
	}
	utils.WriteResponse(result, w)
}
