package mrequest

import (
	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

type ProductCreate struct {
	ProductType        string          `bson:"ProductType" json:"ProductType,omitempty" valid:"required~Field token cannot be empty or is missing,in(P|S|O)~Must be P|S|O"`
	ProductCode        string          `bson:"ProductCode" json:"ProductCode,omitempty" valid:"required~Field token cannot be empty or is missing"`
	ProductGroup       string          `bson:"ProductGroup" json:"ProductGroup,omitempty" valid:"runelength(1|50)~Must be between 1 and 50 characters"`
	ProductDescription string          `bson:"ProductDescription" json:"ProductDescription,omitempty" valid:"required~Field token cannot be empty or is missing,runelength(2|200)~Must be between 2 and 200 characters"`
	ProductNumberCode  string          `bson:"ProductNumberCode" json:"ProductNumberCode,omitempty" valid:"required~Field token cannot be empty or is missing,runelength(1|60)~Must be between 1 and 60 characters"`
	CustomsDetails     *CustomsDetails `bson:"CustomsDetails" json:"CustomsDetails,omitempty"`
}

type ProductRead struct {
	ID                 objectid.ObjectID `json:"id,omitempty" bson:"_id"`
	ProductType        string            `json:"ProductType,omitempty" bson:"ProductType"`
	ProductCode        string            `json:"ProductCode,omitempty" bson:"ProductCode"`
	ProductGroup       string            `json:"ProductGroup,omitempty" bson:"ProductGroup"`
	ProductDescription string            `json:"ProductDescription,omitempty" bson:"ProductDescription"`
	ProductNumberCode  string            `json:"ProductNumberCode,omitempty" bson:"ProductNumberCode"`
	CustomsDetails     *CustomsDetails   `json:"CustomsDetails,omitempty" bson:"CustomsDetails"`
}

type ProductUpdate struct {
	ProductType        string          `bson:"ProductType" json:"ProductType,omitempty" valid:"required~Field token cannot be empty or is missing,in(P|S|O)~Must be P|S|O"`
	ProductCode        string          `bson:"ProductCode" json:"ProductCode,omitempty" valid:"required~Field token cannot be empty or is missing"`
	ProductGroup       string          `bson:"ProductGroup" json:"ProductGroup,omitempty" valid:"runelength(1|50)~Must be between 1 and 50 characters"`
	ProductDescription string          `bson:"ProductDescription" json:"ProductDescription,omitempty" valid:"required~Field token cannot be empty or is missing,runelength(2|200)~Must be between 2 and 200 characters"`
	ProductNumberCode  string          `bson:"ProductNumberCode" json:"ProductNumberCode,omitempty" valid:"required~Field token cannot be empty or is missing,runelength(1|60)~Must be between 1 and 60 characters"`
	CustomsDetails     *CustomsDetails `bson:"CustomsDetails" json:"CustomsDetails,omitempty"`
}

type ProductDelete struct {
	ID                 objectid.ObjectID `bson:"_id" json:"id,omitempty" valid:"required~Cannot be empty" bson:"_id"`
	ProductType        string            `bson:"ProductType" json:"ProductType,omitempty" bson:"ProductType"`
	ProductCode        string            `bson:"ProductCode" json:"ProductCode,omitempty" bson:"ProductCode"`
	ProductGroup       string            `bson:"ProductGroup" json:"ProductGroup,omitempty" bson:"ProductGroup"`
	ProductDescription string            `bson:"ProductDescription" json:"ProductDescription,omitempty" bson:"ProductDescription"`
	ProductNumberCode  string            `bson:"ProductNumberCode" json:"ProductNumberCode,omitempty" bson:"ProductNumberCode"`
	CustomsDetails     *CustomsDetails   `bson:"CustomsDetails" json:"CustomsDetails,omitempty" bson:"CustomsDetails"`
}

type CustomsDetails struct {
	CNCode   []string `json:"CNCode" bson:"CNCode"`
	UNNumber []string `json:"UNNumber" bson:"UNNumber"`
}
