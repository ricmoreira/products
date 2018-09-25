package services

import (
	"context"
	"log"
	"products/models/request"
	"products/models/response"
	"products/repositories"
	"products/util/errors"

	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/mongodb/mongo-go-driver/mongo"
)

// ProductServiceContract is the abstraction for service layer on products resource
type ProductServiceContract interface {
	CreateOne(request *mrequest.ProductCreate) (*mresponse.ProductCreate, *mresponse.ErrorResponse)
	CreateMany(request *[]*mrequest.ProductCreate) (*[]*mresponse.ProductCreate, *mresponse.ErrorResponse)
	List(request *mrequest.ListRequest) (*mresponse.ProductList, *mresponse.ErrorResponse)
}

// ProductService is the layer between http client and repository for product resource
type ProductService struct {
	productRepository repositories.ProductRepositoryContract
}

// NewProductService is the constructor of ProductService
func NewProductService(pr repositories.ProductRepositoryContract) ProductServiceContract {
	return &ProductService{
		productRepository: pr,
	}
}

// CreateOne saves provided model instance to database
func (this *ProductService) CreateOne(request *mrequest.ProductCreate) (*mresponse.ProductCreate, *mresponse.ErrorResponse) {

	// validate request
	e := errors.ValidateRequest(request)
	if e != nil {
		return nil, e
	}

	res, err := this.productRepository.CreateOne(request)

	if err != nil {
		errR := errors.HandleErrorResponse(errors.SERVICE_UNAVAILABLE, nil, err.Error())
		return nil, errR
	}

	id := res.InsertedID.(objectid.ObjectID)
	p := mresponse.ProductCreate{
		ID: id.Hex(),
	}

	return &p, nil
}

// CreateMany saves many products in one bulk operation
func (this *ProductService) CreateMany(request *[]*mrequest.ProductCreate) (*[]*mresponse.ProductCreate, *mresponse.ErrorResponse) {

	res, err := this.productRepository.InsertMany(request)

	if err != nil {
		mngBulkError := err.(mongo.BulkWriteError)
		writeErrors := mngBulkError.WriteErrors
		for _, err := range writeErrors {
			log.Println(err) // for now only print errors
		}
	}

	var length int
	result := make([]*mresponse.ProductCreate, length)
	if res != nil {
		length = len(res.InsertedIDs)
		for _, insertedID := range res.InsertedIDs {
			id := insertedID.(objectid.ObjectID)
			result = append(result, &mresponse.ProductCreate{
				ID: id.Hex(),
			})
		}
	}
	
	return &result, nil
}

// List returns a list of products with pagination and filtering options
func (this *ProductService) List(request *mrequest.ListRequest) (*mresponse.ProductList, *mresponse.ErrorResponse) {

	total, perPage, page, cursor, err := this.productRepository.List(request)

	if err != nil {
		e := errors.HandleErrorResponse(errors.SERVICE_UNAVAILABLE, nil, err.Error())
		return nil, e
	}

	docs := []*mresponse.ProductRead{}

	for cursor.Next(context.Background()) {
		doc := mresponse.ProductRead{}
		err := cursor.Decode(&doc)
		if err != nil {
			errR := errors.HandleErrorResponse(errors.SERVICE_UNAVAILABLE, nil, err.Error())
			return nil, errR
		}

		doc.ID = doc.IDdb.Hex()

		docs = append(docs, &doc)
	}
	
	resp := mresponse.ProductList{
		Total: total,
		PerPage: perPage,
		Page: page,
		Items: &docs,
	}
	return &resp, nil
}
