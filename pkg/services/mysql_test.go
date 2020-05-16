package services

import (
	"database/sql"
	"io/ioutil"
	"log"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/panospet/small-api/pkg/model"
)

type Suite struct {
	suite.Suite
	appDb  AppDb
	dbMock sqlmock.Sqlmock
	db     *sql.DB
}

func (s *Suite) SetupTest() {
	// no need to log stuff
	log.SetFlags(0)
	log.SetOutput(ioutil.Discard)

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal("could not initialize dbMock test db class")
	}
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	db := AppDb{
		Conn: sqlxDB,
	}

	s.appDb = db
	s.dbMock = mock
	s.db = mockDB
}

func (s *Suite) TestGetAllProducts() {
	rows := sqlmock.NewRows([]string{"id", "category_id", "title", "image_url", "price", "description", "created_at", "updated_at"}).AddRow(
		uuid.New().String(), 2, "test title", "http://www.bestprice.gr/test.png", 100, "test description", time.Now(), time.Now()).AddRow(
		uuid.New().String(), 5, "test title 2", "http://www.bestprice.gr/test222.png", 200, "test description 2", time.Now(), time.Now())
	s.dbMock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM product")).WillReturnRows(rows)
	res, err := s.appDb.GetAllProducts(0, 0, "id", true)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 2, len(res))
	assert.Equal(s.T(), "test description", res[0].Description)
	assert.Equal(s.T(), "test description 2", res[1].Description)
	assert.Equal(s.T(), float32(100), res[0].Price)
	assert.Equal(s.T(), float32(200), res[1].Price)
}

func (s *Suite) TestGetProduct() {
	rows := sqlmock.NewRows([]string{"id", "category_id", "title", "image_url", "price", "description", "created_at", "updated_at"}).AddRow(
		uuid.New().String(), 2, "test title", "http://www.bestprice.gr/test.png", 100, "test description", time.Now(), time.Now())
	s.dbMock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM product WHERE id=?")).WithArgs("asdf").WillReturnRows(rows)
	res, err := s.appDb.GetProduct("asdf")
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), "test description", res.Description)
}

func (s *Suite) TestAddProduct() {
	product := model.Product{
		CategoryId:  3,
		Title:       "test",
		ImageUrl:    "http://www.bestprice.gr/test.png",
		Price:       12.12,
		Description: "test description",
	}
	q := `INSERT INTO product (id, category_id, title, image_url, price, description) VALUES (?,?,?,?,?,?);`
	s.dbMock.ExpectExec(regexp.QuoteMeta(q)).WithArgs(sqlmock.AnyArg(), product.CategoryId,
		product.Title, product.ImageUrl, product.Price, product.Description).WillReturnResult(
		sqlmock.NewResult(1, 1))
	id, err := s.appDb.AddProduct(product)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 36, len(id))
}

func (s *Suite) TestUpdateProduct() {
	id := uuid.New().String()
	product := model.Product{
		Id:          id,
		CategoryId:  3,
		Title:       "test",
		ImageUrl:    "http://www.bestprice.gr/test.png",
		Price:       12.12,
		Description: "test description",
	}
	q := `UPDATE product SET category_id=?, title=?, image_url=?, price=?, description=?, updated_at=NOW() WHERE id=?`
	s.dbMock.ExpectExec(regexp.QuoteMeta(q)).WithArgs(product.CategoryId,
		product.Title, product.ImageUrl, product.Price, product.Description, product.Id).WillReturnResult(
		sqlmock.NewResult(1, 1))
	err := s.appDb.UpdateProduct(product)
	assert.Nil(s.T(), err)
}

func (s *Suite) TestGetAllCategories() {
	rows := sqlmock.NewRows([]string{"id", "title", "pos", "image_url", "created_at", "updated_at"}).AddRow(
		1, "cat1", 2, "http://www.bestprice.gr/cat1.png", time.Now(), time.Now()).AddRow(
		2, "cat2", 6, "http://www.bestprice.gr/cat2.png", time.Now(), time.Now())
	s.dbMock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM category")).WillReturnRows(rows)
	res, err := s.appDb.GetAllCategories(0, 0, "id", true)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 2, len(res))
	assert.Equal(s.T(), "cat1", res[0].Title)
	assert.Equal(s.T(), "cat2", res[1].Title)
	assert.Equal(s.T(), 2, res[0].Position)
	assert.Equal(s.T(), 6, res[1].Position)
}

func (s *Suite) TestGetCategory() {
	rows := sqlmock.NewRows([]string{"id", "title", "pos", "image_url", "created_at", "updated_at"}).AddRow(
		1, "cat1", 2, "http://www.bestprice.gr/cat1.png", time.Now(), time.Now())
	s.dbMock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM category WHERE id=?")).WithArgs(10).WillReturnRows(rows)
	res, err := s.appDb.GetCategory(10)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), "cat1", res.Title)
}

func (s *Suite) TestAddCategory() {
	category := model.Category{
		Title:    "cat2",
		Position: 2,
		ImageUrl: "http://www.bestprice.gr/cat2.png",
	}
	q := `INSERT INTO category (title, pos, image_url) VALUES (?,?,?);`
	s.dbMock.ExpectExec(regexp.QuoteMeta(q)).WithArgs(category.Title, category.Position, category.ImageUrl).WillReturnResult(
		sqlmock.NewResult(1, 1))
	err := s.appDb.AddCategory(category)
	assert.Nil(s.T(), err)
}

func (s *Suite) TestUpdateCategory() {
	category := model.Category{
		Id:       44,
		Title:    "cat2",
		Position: 2,
		ImageUrl: "http://www.bestprice.gr/cat2.png",
	}
	q := `UPDATE category SET title=?, pos=?, image_url=?, updated_at=NOW() WHERE id=?`
	s.dbMock.ExpectExec(regexp.QuoteMeta(q)).WithArgs(category.Title, category.Position, category.ImageUrl, category.Id).WillReturnResult(
		sqlmock.NewResult(1, 1))
	err := s.appDb.UpdateCategory(category)
	assert.Nil(s.T(), err)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
