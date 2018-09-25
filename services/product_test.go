package services

import (
	"errors"
	"log"
	"products/models/request"
	"products/models/response"
	"products/repositories"
	"testing"

	"context"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/mongodb/mongo-go-driver/mongo"
	"go.uber.org/dig"
)

// Mock ProductRepository behaviour
type ProductRepositoryMock struct{}

func NewProductRepositoryMock() repositories.ProductRepositoryContract {
	pr := ProductRepositoryMock{}
	return &pr
}

func (prm *ProductRepositoryMock) CreateOne(request *mrequest.ProductCreate) (*mongo.InsertOneResult, error) {
	if request.ProductCode == "product-code-that-cause-repository-error" {
		return nil, errors.New("error ocurred on repository")
	}

	if request.ProductCode == "product-code-for-success" {
		res := mongo.InsertOneResult{}
		id, _ := objectid.FromHex("507f191e810c19729de860ea")
		res.InsertedID = id

		return &res, nil
	}

	return nil, nil
}

func (prm *ProductRepositoryMock) ReadOne(p *mrequest.ProductRead) (*mresponse.Product, error) {
	return nil, nil
}

func (prm *ProductRepositoryMock) InsertMany(request *[]*mrequest.ProductCreate) (*mongo.InsertManyResult, error) {

	res := mongo.InsertManyResult{}
	res.InsertedIDs = make([]interface{}, 0)
	for i, productCreate := range *request {
		if productCreate.ProductCode == "product-code-for-error" {
			e := mongo.BulkWriteError{}
			return nil, e
		}

		var id objectid.ObjectID
		if i == 0 {
			id, _ = objectid.FromHex("507f191e810c19729de860ea")
		}

		if i == 1 {
			id, _ = objectid.FromHex("507f191e810c19729de860eb")
		}

		res.InsertedIDs = append(res.InsertedIDs, id)
	}

	return &res, nil
}

// return values: total, perPage, page, cursor, error - these are the return values 
func (prm *ProductRepositoryMock) List(req *mrequest.ListRequest) (int64, int64, int64, mongo.Cursor, error) {
	// send an error
	if req.Order != "normal" && req.Order != "reverse" {
		return 0, 0, 0, nil, errors.New("invalid order type")
	}

	// send a cursor that will cause an error
	if req.Page == 3 {
		cursor := MongoCursorMock{
			Size:     6,
			Position: 0,
		}
		return 3, 10, 3, &cursor, nil
	}

	// successful cursor
	if req.Page == 1 && req.PerPage == 10 {
		cursor := MongoCursorMock{
			Size:     2,
			Position: 0,
		}
		return 2, 10, 1, &cursor, nil
	}

	return 0, 0, 0, nil, nil
}

// Mock Mongo cursor behaviour
type MongoCursorMock struct {
	Size     int
	Position int
}

func (mc *MongoCursorMock) ID() int64 {
	return 0
}

func (mc *MongoCursorMock) Next(context.Context) bool {
	if mc.Position >= mc.Size {
		return false
	}

	mc.Position++

	return true
}

func (mc *MongoCursorMock) Decode(obj interface{}) error {
	if mc.Position == 5 { // position that will cause an error
		return errors.New("error decoding")
	}

	return nil
}

func (mc *MongoCursorMock) DecodeBytes() (bson.Reader, error) {
	return nil, nil
}

func (mc *MongoCursorMock) Err() error {
	return nil
}

func (mc MongoCursorMock) Close(context.Context) error {
	return nil
}

// dependency injection provided
func buildTestProductContainer() *dig.Container {

	container := dig.New()

	// product repository
	err := container.Provide(NewProductRepositoryMock)
	if err != nil {
		panic(err)
	}

	// product service
	err = container.Provide(NewProductService)
	if err != nil {
		panic(err)
	}

	return container
}

func TestCreateOneErrorOnRequestValidation(t *testing.T) {
	container := buildTestProductContainer()

	err := container.Invoke(func(ps ProductServiceContract) {
		pc := mrequest.ProductCreate{
			// missing required field ProductType to cause an error on validation
			ProductCode:        "some-product-code",
			ProductGroup:       "some-product-group",
			ProductDescription: "some-product-description",
			ProductNumberCode:  "some-product-number-code",
		}

		resp, err := ps.CreateOne(&pc)

		if resp != nil {
			t.Fail()
		}

		if err.Code != "INVALID_REQUEST" {
			t.Fail()
		}
	})

	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}
}

func TestCreateOneErrorOnProductRepository(t *testing.T) {
	container := buildTestProductContainer()

	err := container.Invoke(func(ps ProductServiceContract) {
		pc := mrequest.ProductCreate{
			ProductType:        "P",
			ProductCode:        "product-code-that-cause-repository-error",
			ProductGroup:       "some-product-group",
			ProductDescription: "some-product-description",
			ProductNumberCode:  "some-product-number-code",
		}

		resp, err := ps.CreateOne(&pc)

		if resp != nil {
			t.Fail()
		}

		if err.Response != "error ocurred on repository" {
			t.Fail()
		}

	})

	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}
}

func TestCreateOneSuccess(t *testing.T) {
	container := buildTestProductContainer()

	err := container.Invoke(func(ps ProductServiceContract) {
		pc := mrequest.ProductCreate{
			ProductType:        "P",
			ProductCode:        "product-code-for-success",
			ProductGroup:       "some-product-group",
			ProductDescription: "some-product-description",
			ProductNumberCode:  "some-product-number-code",
		}

		resp, err := ps.CreateOne(&pc)

		if err != nil {
			t.Fail()
		}

		if resp == nil {
			t.Fail()
		}

		if resp.ID != "507f191e810c19729de860ea" {
			t.Fail()
		}
	})

	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}
}

func TestCreateManyErrorOnRepository(t *testing.T) {
	container := buildTestProductContainer()

	err := container.Invoke(func(ps ProductServiceContract) {
		pc1 := mrequest.ProductCreate{
			ProductType:        "P",
			ProductCode:        "product-code-for-error",
			ProductGroup:       "some-product-group",
			ProductDescription: "some-product-description",
			ProductNumberCode:  "some-product-number-code",
		}

		pc2 := mrequest.ProductCreate{
			ProductType:        "P",
			ProductCode:        "product-code-two",
			ProductGroup:       "some-product-group",
			ProductDescription: "some-product-description",
			ProductNumberCode:  "some-product-number-code",
		}

		req := make([]*mrequest.ProductCreate, 0)
		req = append(req, &pc1)
		req = append(req, &pc2)

		res, err := ps.CreateMany(&req)

		if err != nil {
			t.Fail()
		}

		if len(*res) != 0 {
			t.Fail()
		}
	})

	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}
}

func TestCreateManySuccess(t *testing.T) {
	container := buildTestProductContainer()

	err := container.Invoke(func(ps ProductServiceContract) {
		pc1 := mrequest.ProductCreate{
			ProductType:        "P",
			ProductCode:        "product-code-one",
			ProductGroup:       "some-product-group",
			ProductDescription: "some-product-description",
			ProductNumberCode:  "some-product-number-code",
		}

		pc2 := mrequest.ProductCreate{
			ProductType:        "P",
			ProductCode:        "product-code-two",
			ProductGroup:       "some-product-group",
			ProductDescription: "some-product-description",
			ProductNumberCode:  "some-product-number-code",
		}

		req := make([]*mrequest.ProductCreate, 0)
		req = append(req, &pc1)
		req = append(req, &pc2)

		res, err := ps.CreateMany(&req)

		if err != nil {
			t.Fail()
		}

		if len(*res) != 2 {
			t.Fail()
		}
	})

	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}
}

func TestCreateListErrorOnProductRepository(t *testing.T) {
	container := buildTestProductContainer()

	err := container.Invoke(func(ps ProductServiceContract) {

		req := mrequest.ListRequest{
			Order: "order that will cause error",
		}

		_, err := ps.List(&req)

		if err == nil {
			t.Fail()
		}

		if err.Response != "invalid order type" {
			t.Fail()
		}
	})

	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}
}

func TestCreateListErrorOnProductDecode(t *testing.T) {
	container := buildTestProductContainer()

	err := container.Invoke(func(ps ProductServiceContract) {

		req := mrequest.ListRequest{
			Order: "normal",
			Page:  3, // page that will cause an error on decode
		}

		_, err := ps.List(&req)

		if err == nil {
			t.Fail()
		}

		if err.Response != "error decoding" {
			t.Fail()
		}
	})

	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}
}

func TestCreateListSuccess(t *testing.T) {
	container := buildTestProductContainer()

	err := container.Invoke(func(ps ProductServiceContract) {

		req := mrequest.ListRequest{
			Order: "normal",
			Page:  1,
			PerPage: 10,
			Sort: "id",
		}

		succ, err := ps.List(&req)

		if err != nil {
			t.Fail()
		}

		if succ.Total != 2 || succ.PerPage != 10 || len(*succ.Items) != 2 || succ.Page != 1 {
			t.Fail()
		}
	})

	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}
}
