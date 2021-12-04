package main

import (
	"fmt"

	"github.com/Tushar3099/CollectPlugin/models"
)

type googleSheets struct {
	Link string `json:"link" bson:"link"`
}

func (g *googleSheets) Execute(form *models.Form) error {
	fmt.Println(g.Link)
	// fmt.Println("Hello Universe")
	return nil
}

func (g *googleSheets) Initialize(form *models.Form) error {
	g.validate()
	g.Link = "tujhse kya matlab"
	// fmt.Println("Hello Universe")
	return nil
}

func (g *googleSheets) validate() error {
	//validates if the action requires specific fields in forms
	return nil
}

// exported as symbol named "Action"
var Action googleSheets
