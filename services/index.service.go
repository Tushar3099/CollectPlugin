package services

import (
	"net/http"
)

func PostUser(w http.ResponseWriter, r *http.Request) {

}

func GetForm(w http.ResponseWriter, r *http.Request) {

}

func PutForm(w http.ResponseWriter, r *http.Request) {

}

func PostForm(w http.ResponseWriter, r *http.Request) {
	// col := database.Client.Database("Go-test").Collection("forms");
	// questions := []models.Question{
	// 	{
	// 		Title: "What is you gender",
	// 		Type: "select",
	// 		Responses: []models.Response{
	// 			{
	// 				Data: "Male",
	// 			},
	// 			{
	// 				Data: "Female",
	// 			},
	// 			{
	// 				Data: "Others",
	// 			},
	// 		},
	// 	},
	// }
	// answers := []models.Answer{
	// 	{
	// 		Data: "Male",
	// 	},
	// },
	// form := models.Form{
	// 	Questions: questions,
	// 	Action: "Google_Sheets",
	// }
	// col.InsertOne(context.TODO(),form)
	w.Write([]byte("<h1>Hello World!</h1>"))
}
