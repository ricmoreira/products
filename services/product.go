package services

import (
	"context"
	"fmt"
	"products/models/request"
	"products/models/response"
	"products/repositories"
	"products/util/errors"

	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/mongodb/mongo-go-driver/mongo"
)

// ProductServiceContract is the abstraction for service layer on roles resource
type ProductServiceContract interface {
	CreateOne(*mrequest.ProductCreate) (*mresponse.ProductCreate, *mresponse.ErrorResponse)
	ReadOne(*mrequest.ProductRead) (*mresponse.Product, *mresponse.ErrorResponse)
	UpdateOne(*mrequest.ProductUpdate) (*mresponse.Product, *mresponse.ErrorResponse)
	DeleteOne(*mrequest.ProductDelete) (*mresponse.Product, *mresponse.ErrorResponse)
	CreateMany(*[]*mrequest.ProductCreate) (*[]*mresponse.ProductCreate, *mresponse.ErrorResponse)
	List(request *mrequest.ListRequest) (*mresponse.ProductList, *mresponse.ErrorResponse)
}

// ProductService is the layer between http client and repository for product resource
type ProductService struct {
	productRepository *repositories.ProductRepository
}

// NewProductService is the constructor of ProductService
func NewProductService(pr *repositories.ProductRepository) *ProductService {
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

// TODO: implement
func (this *ProductService) ReadOne(p *mrequest.ProductRead) (*mresponse.Product, *mresponse.ErrorResponse) {
	return nil, nil
}

// TODO: implement
func (this *ProductService) UpdateOne(p *mrequest.ProductUpdate) (*mresponse.Product, *mresponse.ErrorResponse) {
	return nil, nil
}

// TODO: implement
func (this *ProductService) DeleteOne(p *mrequest.ProductDelete) (*mresponse.Product, *mresponse.ErrorResponse) {
	return nil, nil
}

// CreateMany saves many products in one bulk operation
func (this *ProductService) CreateMany(request *[]*mrequest.ProductCreate) (*[]*mresponse.ProductCreate, *mresponse.ErrorResponse) {

	res, err := this.productRepository.InsertMany(request)

	if err != nil {
		mngBulkError := err.(mongo.BulkWriteError)
		writeErrors := mngBulkError.WriteErrors
		for _, err := range writeErrors {
			fmt.Println(err)
		}
	}

	result := make([]*mresponse.ProductCreate, len(res.InsertedIDs))
	for i, insertedID := range res.InsertedIDs {
		id := insertedID.(objectid.ObjectID)
		result[i] = &mresponse.ProductCreate{
			ID: id.Hex(),
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
		if err := cursor.Decode(&doc); err != nil {
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
		Products: &docs,
	}
	return &resp, nil
}
