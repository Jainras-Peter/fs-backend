package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Forwarder struct {
	ID                   primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ForwarderID          string             `bson:"forwarderId" json:"forwarderId"`
	ForwarderCompanyName string             `bson:"companyName" json:"companyName"`
	ContactPhone         string             `bson:"phone" json:"phone"`
	Email                string             `bson:"email" json:"email"`
	FullAddress          string             `bson:"address" json:"address"`
	Logo                 string             `bson:"logo" json:"logo"`
	TermsAndConditions   string             `bson:"termsAndConditions" json:"termsAndConditions"`
	DefaultLanguage      string             `bson:"defaultLanguage" json:"defaultLanguage"`
	Username             string             `bson:"username" json:"username"`
	Password             string             `bson:"password" json:"password,omitempty"` // omitempty helps keep it out of some JSON responses if cleared
}
