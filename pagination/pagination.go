package pagination

import (
	"fmt"
	"math"
	"net/url"

	"github.com/gin-gonic/gin"
)

// Request represents common pagination request parameters
type Request struct {
	Page     int32 `form:"page" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=100"`
}

// Response is a generic pagination response structure
type Response[T any] struct {
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Count    int64   `json:"count"`
	PageSize int32   `json:"page_size"`
	Results  []T     `json:"results"`
}

// Params represents database query parameters
type Params struct {
	Limit  int32
	Offset int32
}

// GetParams converts a Request into database parameters
func (r *Request) GetParams() Params {
	return Params{
		Limit:  r.PageSize,
		Offset: (r.Page - 1) * r.PageSize,
	}
}

// NewResponse creates a paginated response
func NewResponse[T any](
	ctx *gin.Context,
	req Request,
	results []T,
	totalCount int64,
) Response[T] {
	totalPages := int32(math.Ceil(float64(totalCount) / float64(req.PageSize)))

	baseURL := ctx.Request.URL
	var nextLink, prevLink *string

	if req.Page < totalPages {
		next := createPaginationURL(baseURL, req.Page+1, req)
		nextLink = &next
	}

	if req.Page > 1 {
		prev := createPaginationURL(baseURL, req.Page-1, req)
		prevLink = &prev
	}

	return Response[T]{
		Next:     nextLink,
		Previous: prevLink,
		Count:    totalCount,
		PageSize: req.PageSize,
		Results:  results,
	}
}

func createPaginationURL(baseURL *url.URL, page int32, req Request) string {
	q := baseURL.Query()
	q.Set("page", fmt.Sprintf("%d", page))
	q.Set("page_size", fmt.Sprintf("%d", req.PageSize))
	baseURL.RawQuery = q.Encode()
	return baseURL.String()
}
