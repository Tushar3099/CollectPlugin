package main

import (
	"fmt"

	"github.com/Tushar3099/CollectPlugin/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Structure models.GoogleSheets

type NewForm struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserId     primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
	Title      string             `json:"title" bson:"title" `
	Questions  []models.Question  `json:"questions" bson:"questions"`
	Answers    []models.Answer    `json:"answers" bson:"answers"`
	ActionName string             `json:"action_name" bson:"action_name"`
	Action     Structure          `json:"action" bson:"action"`
}

func Constructor(b *[]byte) *models.Form {
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
	return fmt.Errorf("Invalid Answer")
}

func (g *Structure) Initialize(form *models.Form) error {
	if err := g.validate(); err != nil {
		return err
	}
	// var ok bool
	act, _ := form.Action.(*Structure)

	act.Link = "tujhse kya matlab"
	form.Action = act
	// fmt.Println("Hello Universe")
	return nil
}

func (g *Structure) validate() error {
	//validates if the action requires specific fields in forms
	return nil
}

// exported as symbol named "Action"
var Action Structure
