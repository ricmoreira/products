package mrequest

import (
	"net/url"
	"strconv"
)

type ListRequest struct {
	PerPage int                    `json:"per_page" valid:"required"`
	Page    int                    `json:"page" valid:"required"`
	Sort    string                 `json:"sort" valid:"required in(id|_id)"`
	Order   string                 `json:"order" valid:"required in(normal|reverse)"`
	Filters map[string]interface{} `json:"filters" valid:""`
}

// NewListRequest creates a ListRequest from params sent in URL query string
// url example: http://products?per_page=10&page=1&sort=id&order=normal
func NewListRequest(params url.Values, allowedSorts map[string]string, allowedFilters map[string]string) *ListRequest {
	allowedOrders := make(map[string]string)
	allowedOrders["normal"] = "normal"
	allowedOrders["reverse"] = "reverse"

	var req ListRequest

	// set sort
	if sort := params.Get("sort"); sort != "" {
		req.Sort = sort
	} else {
		req.Sort = "id"
	}

	// set order
	if order := params.Get("order"); order != "" {
		req.Order = order
	} else {
		req.Order = "normal"
	}

	// set per_page
	if ok := params.Get("per_page"); ok != "" {
		req.PerPage, _ = strconv.Atoi(params.Get("per_page"))
	}
	if req.PerPage <= 0 {
		req.PerPage = 20
	}

	// set page
	if ok := params.Get("page"); ok != "" {
		req.Page, _ = strconv.Atoi(params.Get("page"))
	}
	if req.Page <= 0 {
		req.Page = 1
	}

	// set filter
	req.Filters = make(map[string]interface{})
	for _, filter := range allowedFilters {
		if val, ok := params[filter]; ok {
			req.Filters[filter] = val[0]
		}
	}

	return &req
}
