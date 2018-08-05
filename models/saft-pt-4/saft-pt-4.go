package msaft

import (
	"encoding/xml"
)

type Product struct {
	XMLName            xml.Name        `xml:"Product"`
	ProductType        string          `bson:"ProductType" json:"ProductType" xml:"ProductType"`
	ProductCode        string          `bson:"ProductCode" json:"ProductCode" xml:"ProductCode"`
	ProductGroup       string          `bson:"ProductGroup" json:"ProductGroup" xml:"ProductGroup"`
	ProductDescription string          `bson:"ProductDescription" json:"ProductDescription" xml:"ProductDescription"`
	ProductNumberCode  string          `bson:"ProductNumberCode" json:"ProductNumberCode" xml:"ProductNumberCode"`
	CustomsDetails     *CustomsDetails `bson:"CustomsDetails" json:"CustomsDetails" xml:"CustomsDetails"`
}

type CustomsDetails struct {
	XMLName  xml.Name `xml:"CustomsDetails"`
	CNCode   string   `bson:"CNCode" json:"CNCode" xml:"CNCode"`
	UNNumber string   `bson:"UNNumber" json:"UNNumber" xml:"UNNumber"`
}

type AuditFile struct {
	XMLName  xml.Name  `xml:"AuditFile"`
	Products []*Product `json:"Products" xml:"MasterFiles>Product"`
}
