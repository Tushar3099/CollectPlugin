package models

import "go.mongodb.org/mongo-driver/bson/primitive"

//Create Struct
type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name,omitempty" bson:"name,omitempty"`
	Email    string             `json:"email,omitempty" bson:"email,omitempty"`
	Password string             `json:"password,omitempty" bson:"password,omitempty"`
	Phone    string             `json:"phone,omitempty" bson:"phone,omitempty"`
}

type Action interface {
	Constructor(b *[]byte) *Form
	Execute(form *Form, ans *Answer) error
	Initialize(form *Form) error
}

type Form struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserId     primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
	Title      string             `json:"title" bson:"title" `
	Questions  []Question         `json:"questions" bson:"questions"`
	Answers    []Answer           `json:"answers" bson:"answers"`
	ActionName string             `json:"action_name" bson:"action_name"`
	Action     Action             `json:"action" bson:"action"`
}

type Question struct {
	Title     string     `json:"title" bson:"title"`
	Type      string     `json:"type" bson:"type"`
	Responses []Response `json:"responses" bson:"responses" `
}

type Response struct {
	Data string `json:"data" bson:"data"`
}

type Answer struct {
	UserId primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
	List   []string           `json:"list,omitempty" bson:"list,omitempty"`
}
