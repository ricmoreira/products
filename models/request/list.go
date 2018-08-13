package mrequest

import (
	"net/url"
	"strconv"
)

type ListRequest struct {
	PerPage int                    `json:"per_page" valid:"required"`
	Page    int                    `json:"page" valid:"required"`
	Sort    map[string]interface{} `json:"sort" valid:"required in(id)"`
	Filters map[string]interface{} `json:"filters" valid:""`
}

// NewListRequest creates a ListRequest from params sent in URL query string
// url example: http://products?per_page=10&page=1&sort=id&order=normal
func NewListRequest(params url.Values, allowedSorts map[string]string, allowedFilters map[string]string) *ListRequest {
	allowedOrders := make(map[string]string)
	allowedOrders["normal"] = "normal"
	allowedOrders["reverse"] = "reverse"

	var req ListRequest

	
	if ok := params.Get("per_page"); ok != "" {
		req.PerPage, _ = strconv.Atoi(params.Get("per_page"))
	}
	if ok := params.Get("page"); ok != "" {
		req.Page, _ = strconv.Atoi(params.Get("page"))
	}

	req.Filters = make(map[string]interface{})
	for _, filter := range allowedFilters {
		if val, ok := params[filter]; ok {
			req.Filters[filter] = val[0]
		}
	}

	// set default values
	if req.Page <= 0 {
		req.Page = 1
	}

	if req.PerPage <= 0 {
		req.PerPage = 20
	}

	req.Sort = make(map[string]interface{})
	if _, ok := allowedSorts[params.Get("sort")]; ok {
		req.Sort[allowedSorts[params.Get("sort")]] = 1
		if _, ok := allowedOrders[params.Get("order")]; ok && params.Get("order") == "reverse" {
			req.Sort[allowedSorts[params.Get("sort")]] = -1
		}
	}

	return &req
}
