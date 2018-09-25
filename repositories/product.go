package repositories

import (
	"context"
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

type ProductRepositoryContract interface {
	CreateOne(request *mrequest.ProductCreate) (*mongo.InsertOneResult, error)
	ReadOne(p *mrequest.ProductRead) (*mresponse.Product, error)
	InsertMany(request *[]*mrequest.ProductCreate) (*mongo.InsertManyResult, error)
	List(req *mrequest.ListRequest) (int64, int64, int64, mongo.Cursor, error)
}

// NewProductRepository is the constructor for ProductRepository
func NewProductRepository(db *DBCollections) ProductRepositoryContract {
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

// List will return a mongo.Cursor along with pagination utility values
// total, perPage, page, cursor, error - these are the return values 
func (this *ProductRepository) List(req *mrequest.ListRequest) (int64, int64, int64, mongo.Cursor, error) {

	args := []*bson.Element{}

	for key, value := range req.Filters {
		if key != "_id" { // filter by text fields
			pattern := value.(string) 
			elem := bson.EC.Regex(key, pattern, "i")
			args = append(args, elem)
		} else { // filter by _id
			elem := bson.EC.String(key, value.(string))
			args = append(args, elem)
		}
	}

	total, e := this.products.Count(
		context.Background(),
		bson.NewDocument(args...),
	)

	sorting := map[string]int{}
	var sortingValue int
	if req.Order == "reverse" {
		sortingValue = -1
	} else {
		sortingValue = 1
	}
	sorting[req.Sort] = sortingValue

	perPage := int64(req.PerPage)
	page := int64(req.Page)
	cursor, e := this.products.Find(
		context.Background(),
		bson.NewDocument(args...),
		findopt.Sort(sorting),
		findopt.Skip(int64(req.PerPage*(req.Page-1))),
		findopt.Limit(perPage),
	)

	return total, perPage, page, cursor, e
}
