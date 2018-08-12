package mresponse

import (
	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

type Product struct {
	ID                 objectid.ObjectID `json:"id,omitempty" bson:"_id"`
	ProductType        string            `json:"ProductType,omitempty" bson:"ProductType"`
	ProductCode        string            `json:"ProductCode,omitempty" bson:"ProductCode"`
	ProductGroup       string            `json:"ProductGroup,omitempty" bson:"ProductGroup,omitempty"`
	ProductDescription string            `json:"ProductDescription,omitempty" bson:"ProductDescription"`
	ProductNumberCode  string            `json:"ProductNumberCode,omitempty" bson:"ProductNumberCode"`
	CustomsDetails     *CustomsDetails   `json:"CustomsDetails,omitempty" bson:"CustomsDetails,omitempty"`
}

type CustomsDetails struct {
	CNCode   []string `json:"CNCode" bson:"CNCode"`
	UNNumber []string `json:"UNNumber" bson:"UNNumber"`
}

type ProductCreate struct {
	ID string `json:"id,omitempty"`
}

type ProductRead struct {
	ID                 string            `json:"id,omitempty"`
	IDdb               objectid.ObjectID `json:"-" bson:"_id"`
	ProductType        string            `json:"ProductType,omitempty" bson:"ProductType"`
	ProductCode        string            `json:"ProductCode,omitempty" bson:"ProductCode"`
	ProductGroup       string            `json:"ProductGroup,omitempty" bson:"ProductGroup,omitempty"`
	ProductDescription string            `json:"ProductDescription,omitempty" bson:"ProductDescription"`
	ProductNumberCode  string            `json:"ProductNumberCode,omitempty" bson:"ProductNumberCode"`
	CustomsDetails     *CustomsDetails   `json:"CustomsDetails,omitempty" bson:"CustomsDetails,omitempty"`
}

type ProductList struct {
	Total    int64           `json:"total"`
	PerPage  int64           `json:"per_page"`
	Page     int64           `json:"page"`
	Products *[]*ProductRead `json:"products"`
}
