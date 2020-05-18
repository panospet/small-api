package api

import (
	"bytes"
	"encoding/json"
	"github.com/panospet/small-api/pkg/cache"
	"github.com/panospet/small-api/pkg/model"
	"github.com/panospet/small-api/pkg/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type Suite struct {
	suite.Suite
	api Api
}

func (s *Suite) SetupTest() {
	s.api = Api{
		Db:    services.NewMockDb(),
		Cache: cache.NewCacherMock(),
	}
}

func (s *Suite) TestHealthCheckHandler() {
	req, err := http.NewRequest("GET", "", nil)
	assert.Nil(s.T(), err)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.api.health)
	handler.ServeHTTP(rr, req)

	assert.Equal(s.T(), http.StatusTeapot, rr.Code)
	expected := `{"Message":"health good!"}`
	assert.Equal(s.T(), expected, rr.Body.String())
}

func (s *Suite) TestGetProducts() {
	req, err := http.NewRequest("GET", "/v1/products", nil)
	assert.Nil(s.T(), err)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.api.getListProducts)
	handler.ServeHTTP(rr, req)

	assert.Equal(s.T(), http.StatusOK, rr.Code)
	bodyBytes := []byte(rr.Body.String())
	var prods []model.Product
	err = json.Unmarshal(bodyBytes, &prods)
	assert.Nil(s.T(), err)
	assert.Len(s.T(), prods, 10)
}

func (s *Suite) TestGetProductsPerPage() {
	req, err := http.NewRequest("GET", "/v1/products?perPage=40", nil)
	assert.Nil(s.T(), err)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.api.getListProducts)
	handler.ServeHTTP(rr, req)

	assert.Equal(s.T(), http.StatusOK, rr.Code)
	bodyBytes := []byte(rr.Body.String())
	var prods []model.Product
	err = json.Unmarshal(bodyBytes, &prods)
	assert.Nil(s.T(), err)
	assert.Len(s.T(), prods, 40)
}

func (s *Suite) TestGetProductsPerPageWithLimit() {
	req, err := http.NewRequest("GET", "/v1/products?perPage=40&limit=20", nil)
	assert.Nil(s.T(), err)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.api.getListProducts)
	handler.ServeHTTP(rr, req)

	assert.Equal(s.T(), http.StatusOK, rr.Code)
	bodyBytes := []byte(rr.Body.String())
	var prods []model.Product
	err = json.Unmarshal(bodyBytes, &prods)
	assert.Nil(s.T(), err)
	assert.Len(s.T(), prods, 20)
}

func (s *Suite) TestGetProduct() {
	req, err := http.NewRequest("GET", "/v1/products/123", nil)
	assert.Nil(s.T(), err)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.api.getProduct)
	handler.ServeHTTP(rr, req)
	assert.Equal(s.T(), http.StatusOK, rr.Code)
	bodyBytes := []byte(rr.Body.String())
	var prod model.Product
	err = json.Unmarshal(bodyBytes, &prod)
	assert.Equal(s.T(), "product0", prod.Title)
}

func (s *Suite) TestPatchProduct() {
	reqBody, err := json.Marshal(map[string]interface{}{
		"title":       "updated",
		"image_url":   "http://www.bestprice.gr/updated.png",
		"price":       123,
		"description": "updated",
	})
	assert.Nil(s.T(), err)
	req, err := http.NewRequest("POST", "/v1/products/123", bytes.NewBuffer(reqBody))
	assert.Nil(s.T(), err)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.api.getProduct)
	handler.ServeHTTP(rr, req)
	assert.Equal(s.T(), http.StatusOK, rr.Code)
	bodyBytes := []byte(rr.Body.String())
	type Res struct {
		Message string
	}
	var res Res
	err = json.Unmarshal(bodyBytes, &res)
	assert.Contains(s.T(), "was updated", res.Message)
}

func (s *Suite) TestGetCategories() {
	req, err := http.NewRequest("GET", "/v1/categories", nil)
	assert.Nil(s.T(), err)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.api.getListCategories)
	handler.ServeHTTP(rr, req)

	assert.Equal(s.T(), http.StatusOK, rr.Code)
	bodyBytes := []byte(rr.Body.String())
	var categories []model.Category
	err = json.Unmarshal(bodyBytes, &categories)
	assert.Nil(s.T(), err)
	assert.Len(s.T(), categories, 10)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
