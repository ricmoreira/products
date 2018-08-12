package controllers

import (
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
func (ps *MockProductService) CreateOne(pReq *mrequest.ProductCreate) (*mresponse.ProductCreate, *mresponse.ErrorResponse) {
	// validate request
	err := errors.ValidateRequest(pReq)
	if err != nil {
		return nil, err
	}

	pRes := mresponse.ProductCreate{}
	pRes.ID = "some-unique-id"

	return &pRes, nil
}

// mocked behaviour for ReadOne
func (ps *MockProductService) ReadOne(p *mrequest.ProductRead) (*mresponse.Product, *mresponse.ErrorResponse) {
	// TODO: implement in the future
	return nil, nil
}

// mocked behaviour for UpdateOne
func (ps *MockProductService) UpdateOne(p *mrequest.ProductUpdate) (*mresponse.Product, *mresponse.ErrorResponse) {
	// TODO: implement in the future
	return nil, nil
}

// mocked behaviour for DeleteOne
func (ps *MockProductService) DeleteOne(p *mrequest.ProductDelete) (*mresponse.Product, *mresponse.ErrorResponse) {
	// TODO: implement in the future
	return nil, nil
}

func (ps *MockProductService) CreateMany(*[]*mrequest.ProductCreate) (*[]*mresponse.ProductCreate, *mresponse.ErrorResponse) {
	// TODO: implement in the future
	return nil, nil
}

func (ps *MockProductService) List(*mrequest.ListRequest) (*mresponse.ProductList, *mresponse.ErrorResponse) {
	// TODO: implement in the future
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
