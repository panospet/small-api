package api

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPagination(t *testing.T) {
	r, _ := http.NewRequest("GET", "http://myapi.gr/v1/products?limit=100&perPage=20", nil)
	pag, err := getPaginationFromRequest(r)
	assert.Nil(t, err)

	assert.Equal(t, 100, pag.limit)
	assert.Equal(t, 20, pag.perPage)
	assert.Equal(t, 0, pag.offset)
	assert.Equal(t, 1, pag.page)
	assert.Equal(t, 0, pag.start)
	assert.Equal(t, 20, pag.end)
}

func TestPagination2(t *testing.T) {
	r, _ := http.NewRequest("GET", "http://myapi.gr/v1/products?limit=10&perPage=20", nil)
	pag, err := getPaginationFromRequest(r)
	assert.Nil(t, err)

	assert.Equal(t, 10, pag.limit)
	assert.Equal(t, 10, pag.perPage)
	assert.Equal(t, 0, pag.offset)
	assert.Equal(t, 1, pag.page)
	assert.Equal(t, 0, pag.start)
	assert.Equal(t, 10, pag.end)
}

func TestPagination3(t *testing.T) {
	r, _ := http.NewRequest("GET", "http://myapi.gr/v1/products?perPage=20&offset=100&page=4", nil)
	pag, err := getPaginationFromRequest(r)
	assert.Nil(t, err)

	assert.Equal(t, 0, pag.limit)
	assert.Equal(t, 20, pag.perPage)
	assert.Equal(t, 100, pag.offset)
	assert.Equal(t, 4, pag.page)
	assert.Equal(t, 160, pag.start)
	assert.Equal(t, 180, pag.end)
}

func TestPagination4(t *testing.T) {
	r, _ := http.NewRequest("GET", "http://myapi.gr/v1/products?perPage=20&offset=100&page=4&limit=1000", nil)
	pag, err := getPaginationFromRequest(r)
	assert.Nil(t, err)

	assert.Equal(t, 1000, pag.limit)
	assert.Equal(t, 20, pag.perPage)
	assert.Equal(t, 100, pag.offset)
	assert.Equal(t, 4, pag.page)
	assert.Equal(t, 160, pag.start)
	assert.Equal(t, 180, pag.end)
}

func TestPagination5(t *testing.T) {
	r, _ := http.NewRequest("GET", "http://myapi.gr/v1/products?offset=7&page=1&limit=100&perPage=1", nil)
	pag, err := getPaginationFromRequest(r)
	assert.Nil(t, err)

	assert.Equal(t, 100, pag.limit)
	assert.Equal(t, 1, pag.perPage)
	assert.Equal(t, 7, pag.offset)
	assert.Equal(t, 1, pag.page)
	assert.Equal(t, 7, pag.start)
	assert.Equal(t, 8, pag.end)
}

func TestCalculatePaginationHeaders(t *testing.T) {
	r, _ := http.NewRequest("GET", "http://myapi.gr/v1/products?perPage=20&offset=100&page=4&limit=1000", nil)
	pag, err := getPaginationFromRequest(r)
	assert.Nil(t, err)

	links := calculatePaginationHeaders(r, pag, 1000)
	assert.Len(t, links, 5)
	assert.Contains(t, links, `<http://myapi.gr/http:/myapi.gr/v1/products?limit=1000&offset=100&page=4&perPage=20>; rel="self"`)
	assert.Contains(t, links, `<http://myapi.gr/http:/myapi.gr/v1/products?limit=1000&offset=100&page=1&perPage=20>; rel="first"`)
	assert.Contains(t, links, `<http://myapi.gr/http:/myapi.gr/v1/products?limit=1000&offset=100&page=5&perPage=20>; rel="next"`)
	assert.Contains(t, links, `<http://myapi.gr/http:/myapi.gr/v1/products?limit=1000&offset=100&page=3&perPage=20>; rel="prev"`)
	assert.Contains(t, links, `<http://myapi.gr/http:/myapi.gr/v1/products?limit=1000&offset=100&page=50&perPage=20>; rel="last"`)
}

func TestCalculatePaginationHeadersNoNextPage(t *testing.T) {
	r, _ := http.NewRequest("GET", "http://myapi.gr/v1/products?perPage=20&page=5&limit=100", nil)
	pag, err := getPaginationFromRequest(r)
	assert.Nil(t, err)

	links := calculatePaginationHeaders(r, pag, 1000)
	assert.Len(t, links, 4)
	assert.Contains(t, links, `<http://myapi.gr/http:/myapi.gr/v1/products?limit=100&page=5&perPage=20>; rel="self"`)
	assert.Contains(t, links, `<http://myapi.gr/http:/myapi.gr/v1/products?limit=100&page=1&perPage=20>; rel="first"`)
	assert.Contains(t, links, `<http://myapi.gr/http:/myapi.gr/v1/products?limit=100&page=4&perPage=20>; rel="prev"`)
	assert.Contains(t, links, `<http://myapi.gr/http:/myapi.gr/v1/products?limit=100&page=5&perPage=20>; rel="last"`)
}

func TestCalculatePaginationHeadersNoPrevPage(t *testing.T) {
	r, _ := http.NewRequest("GET", "http://myapi.gr/v1/products?perPage=20&page=1&limit=100", nil)
	pag, err := getPaginationFromRequest(r)
	assert.Nil(t, err)

	links := calculatePaginationHeaders(r, pag, 1000)
	assert.Len(t, links, 4)
	assert.Contains(t, links, `<http://myapi.gr/http:/myapi.gr/v1/products?limit=100&page=1&perPage=20>; rel="self"`)
	assert.Contains(t, links, `<http://myapi.gr/http:/myapi.gr/v1/products?limit=100&page=1&perPage=20>; rel="first"`)
	assert.Contains(t, links, `<http://myapi.gr/http:/myapi.gr/v1/products?limit=100&page=2&perPage=20>; rel="next"`)
	assert.Contains(t, links, `<http://myapi.gr/http:/myapi.gr/v1/products?limit=100&page=5&perPage=20>; rel="last"`)
}