package main

import (
	"fmt"

	"github.com/Tushar3099/CollectPlugin/models"
)

type Structure struct {
}

func (g *Structure) Constructor(b *[]byte) *models.Form {
	return &models.Form{}
}

func (g *Structure) Execute(form *models.Form, ans *models.Answer) error {
	fmt.Println("Executing SMS action")
	return nil
}

func (g *Structure) Initialize(form *models.Form) error {
	g.validate()
	fmt.Println("Initalizing SMS Action")
	return nil
}

func (g *Structure) validate() error {
	//validates if the action requires specific fields in forms
	return nil
}

// exported as symbol named "Action"
var Action Structure
