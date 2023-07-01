package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Person represents a log document in MongoDB
type LogModel struct {
	ID        primitive.ObjectID    `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string    `json:"name,omitempty" bson:"name,omitempty"`
	Description string `json:"description,omitempty" bson:"description,omitempty"`
	Image string `json:"image,omitempty" bson:"image,omitempty"`
	Attributes []AttributeModel `json:"attributes,omitempty" bson:"attributes,omitempty"`
	TokenId int `json:"tokenid,omitempty" bson:"tokenid,omitempty"`
}

type AttributeModel struct {
	TraitType string `json:"trait_type,omitempty" bson:"trait_type,omitempty"`
	Value int `json:"value,omitempty" bson:"value,omitempty"`
	MaxValue int `json:"max_value,omitempty" bson:"max_value,omitempty"`
}