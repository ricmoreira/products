package controllers

import (
	"products/models/request"
	"products/services"
	"products/util/errors"
	"encoding/json"

	"github.com/gin-gonic/gin"
)

type (
	// ProductController represents the controller for operating on the products resource
	ProductController struct {
		ProductService services.ProductServiceContract
	}
)

// NewProductController is the constructor of ProductController
func NewProductController(ps *services.ProductService) *ProductController {
	return &ProductController{
		ProductService: ps,
	}
}

// CreateAction creates a new product
func (pc ProductController) CreateAction(c *gin.Context) {
	pReq := mrequest.ProductCreate{}
	json.NewDecoder(c.Request.Body).Decode(&pReq)

	e := errors.ValidateRequest(&pReq)
	if e != nil {
		c.JSON(e.HttpCode, e)
		return
	}

	pRes, err := pc.ProductService.CreateOne(&pReq)

	if err != nil {
		c.JSON(err.HttpCode, err)
		return
	}

	c.JSON(200, pRes)
}
