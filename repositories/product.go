package repositories

import (
	"context"
	"fmt"
	"products/models/request"
	"products/models/response"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/findopt"
	"github.com/mongodb/mongo-go-driver/mongo/insertopt"
)

// ProductRepository performs CRUD operations on users resource
type ProductRepository struct {
	products MongoCollection
}

// NewProductRepository is the constructor for ProductRepository
func NewProductRepository(db *DBCollections) *ProductRepository {
	return &ProductRepository{products: db.Product}
}

// CreateOne saves provided model instance to database
func (this *ProductRepository) CreateOne(request *mrequest.ProductCreate) (*mongo.InsertOneResult, error) {

	return this.products.InsertOne(context.Background(), request)
}

// ReadOne returns a product based on ProductCode sent in request
// TODO: implement better query based on full request and not only the ProducCode
func (this *ProductRepository) ReadOne(p *mrequest.ProductRead) (*mresponse.Product, error) {
	result := this.products.FindOne(
		context.Background(),
		bson.NewDocument(bson.EC.String("ProductCode", p.ProductCode)),
	)

	res := mresponse.Product{}
	err := result.Decode(p)

	if err != nil {
		return nil, err
	}

	return &res, nil
}

// TODO: implement
func (this *ProductRepository) UpdateOne(p *mrequest.ProductUpdate) (*mresponse.Product, error) {
	return nil, nil
}

// TODO: implement
func (this *ProductRepository) DeleteOne(p *mrequest.ProductDelete) (*mresponse.Product, error) {
	return nil, nil
}

func (this *ProductRepository) InsertMany(request *[]*mrequest.ProductCreate) (*mongo.InsertManyResult, error) {
	// transform to []interface{} (https://golang.org/doc/faq#convert_slice_of_interface)
	s := make([]interface{}, len(*request))
	for i, v := range *request {
		s[i] = v
	}

	// { ordered: false } ordered is false in order to don't stop execution because an error ocurred on one of the inserts
	opt := insertopt.Ordered(false)
	return this.products.InsertMany(context.Background(), s, opt)
}

func (this *ProductRepository) List(req *mrequest.ListRequest) (int64, int64, int64, mongo.Cursor, error) {

	args := []*bson.Element{}

	for i, v := range req.Filters {
		args = append(args, bson.EC.String(i, fmt.Sprintf("%v", v)))
	}

	total, e := this.products.Count(
		context.Background(),
		bson.NewDocument(args...),
	)

	perPage := int64(req.PerPage)
	page := int64(req.Page)
	cursor, e := this.products.Find(
		context.Background(),
		bson.NewDocument(args...),
		findopt.Sort(req.Sort),
		findopt.Skip(int64(req.PerPage*(req.Page-1))),
		findopt.Limit(perPage),
	)

	return total, perPage, page, cursor, e
}
