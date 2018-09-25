package controllers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"products/models/request"
	"products/models/response"
	"products/util/errors"
	"testing"

	"bytes"
	"encoding/json"

	"github.com/gin-gonic/gin"
)

// stub ProductService behaviour
type MockProductService struct{}

// mocked behaviour for CreateOne
func (ps *MockProductService) CreateOne(request *mrequest.ProductCreate) (*mresponse.ProductCreate, *mresponse.ErrorResponse) {
	// validate request
	err := errors.ValidateRequest(request)
	if err != nil {
		return nil, err
	}

	pRes := mresponse.ProductCreate{}
	pRes.ID = "some-unique-id"

	return &pRes, nil
}

func (ps *MockProductService) CreateMany(request *[]*mrequest.ProductCreate) (*[]*mresponse.ProductCreate, *mresponse.ErrorResponse) {
	// TODO: implement in the future
	return nil, nil
}

func (ps *MockProductService) List(req *mrequest.ListRequest) (*mresponse.ProductList, *mresponse.ErrorResponse) {

	// success case
	if req.Page == 1 && req.PerPage == 10 {
		err := errors.ValidateRequest(req)
		if err != nil {
			fmt.Printf("%v", err)
			return nil, err
		}

		res := mresponse.ProductList{}
		items := make([]*mresponse.ProductRead, 0)
		pRes1 := mresponse.ProductRead{}
		pRes1.ID = "some-id-1"
		pRes2 := mresponse.ProductRead{}
		pRes2.ID = "some-id-2"
		items = append(items, &pRes1)
		items = append(items, &pRes2)
		res.Items = &items

		return &res, nil
	}

		// error
		if req.Page == 1 && req.PerPage == 99 {
			err := mresponse.ErrorResponse{}
			err.HttpCode = 502
			err.Response = "error ocurred on service"
			err.Code = "SERVICE_ERROR"

			return nil, &err
		}

	return nil, nil
}
func TestCreateProductAction(t *testing.T) {

	// Mock the server

	// Switch to test mode in order to don't get such noisy output
	gin.SetMode(gin.TestMode)

	pps := &MockProductService{}

	pc := ProductController{
		ProductService: pps,
	}

	r := gin.Default()

	r.POST("/api/v1/product", pc.CreateAction)

	// TEST SUCCESS

	// Mock a request
	body := mrequest.ProductCreate{
		ProductType:        "P",
		ProductCode:        "some-product-code",
		ProductGroup:       "some-product-group",
		ProductDescription: "some-product-description",
		ProductNumberCode:  "some-product-number-code",
	}

	jsonValue, _ := json.Marshal(body)

	req, err := http.NewRequest(http.MethodPost, "/api/v1/product", bytes.NewBuffer(jsonValue))
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}

	// Create a response recorder in order to inspect the response
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)

	// Do asssertions
	if w.Code != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(w.Body)
		bodyString := string(bodyBytes)

		t.Fatalf("Expected to get status %d but instead got %d\nResponse body:\n%s", http.StatusOK, w.Code, bodyString)
	}
}
func TestListActionOnServiceError(t *testing.T) {

	// Mock the server

	// Switch to test mode in order to don't get such noisy output
	gin.SetMode(gin.TestMode)

	pps := &MockProductService{}

	pc := ProductController{
		ProductService: pps,
	}

	r := gin.Default()

	r.GET("/api/v1/product", pc.ListAction)

	// TEST SUCCESS

	req, err := http.NewRequest(http.MethodGet, "/api/v1/product?per_page=99&page=1&sort=id&order=normal", nil)
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}

	// Create a response recorder in order to inspect the response
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)

	// Do asssertions
	if w.Code != 502 {
		bodyBytes, _ := ioutil.ReadAll(w.Body)
		bodyString := string(bodyBytes)

		t.Fatalf("Expected to get status %d but instead got %d\nResponse body:\n%s", 502, w.Code, bodyString)
	}
}


func TestListActionSuccess(t *testing.T) {

	// Mock the server

	// Switch to test mode in order to don't get such noisy output
	gin.SetMode(gin.TestMode)

	pps := &MockProductService{}

	pc := ProductController{
		ProductService: pps,
	}

	r := gin.Default()

	r.GET("/api/v1/product", pc.ListAction)

	// TEST SUCCESS

	req, err := http.NewRequest(http.MethodGet, "/api/v1/product?per_page=10&page=1&sort=id&order=normal", nil)
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}

	// Create a response recorder in order to inspect the response
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)

	// Do asssertions
	if w.Code != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(w.Body)
		bodyString := string(bodyBytes)

		t.Fatalf("Expected to get status %d but instead got %d\nResponse body:\n%s", http.StatusOK, w.Code, bodyString)
	}
}
