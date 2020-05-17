package api

import (
	"encoding/json"
	"fmt"
	"github.com/panospet/small-api/pkg/model"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/panospet/small-api/pkg/cache"
	"github.com/panospet/small-api/pkg/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
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
	fmt.Println(rr.Body.String())
	bodyBytes := []byte(rr.Body.String())
	var prods []model.Product
	err = json.Unmarshal(bodyBytes, &prods)
	assert.Nil(s.T(), err)
	assert.Len(s.T(), prods, 10)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}