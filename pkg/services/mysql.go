package services

import (
	"fmt"
	"regexp"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/panospet/small-api/internal/passwd"
	"github.com/panospet/small-api/pkg/model"
)

type AppDb struct {
	Conn *sqlx.DB
}

func NewDb(mysqlPath string) (*AppDb, error) {
	db, err := sqlx.Connect("mysql", mysqlPath)
	if err != nil {
		return &AppDb{}, err
	}
	return &AppDb{Conn: db}, nil
}

var valid = regexp.MustCompile("^[A-Za-z0-9_]+$")

func (a *AppDb) GetAllProducts(offset int, limit int, orderBy string, asc bool) ([]model.Product, error) {
	var products []model.Product
	q := "SELECT * FROM product"
	if len(orderBy) > 0 {
		if !valid.MatchString(orderBy) {
			return products, &SqlInjectionAttemptError{}
		}
		sort := "desc"
		if asc == true {
			sort = "asc"
		}
		q += fmt.Sprintf(` ORDER BY %s %s`, orderBy, sort)
	}
	if limit != 0 {
		q += fmt.Sprintf(` LIMIT %d OFFSET %d `, limit, offset)
	}
	rows, err := a.Conn.Queryx(q)
	if err != nil {
		return products, err
	}
	for rows.Next() {
		var p model.Product
		err = rows.StructScan(&p)
		if err != nil {
			return products, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (a *AppDb) GetProduct(id string) (model.Product, error) {
	var product model.Product
	err := a.Conn.QueryRowx("SELECT * FROM product WHERE id=?", id).StructScan(&product)
	if err != nil {
		return model.Product{}, err
	}
	return product, nil
}

func (a *AppDb) AddProduct(product model.Product) (string, error) {
	id := uuid.New().String()
	product.Id = id
	q := `INSERT INTO product (id, category_id, title, image_url, price, description) VALUES (?,?,?,?,?,?);`
	_, err := a.Conn.Exec(q, product.Id, product.CategoryId, product.Title, product.ImageUrl, product.Price, product.Description)
	if err != nil {
		return "", err
	}
	return id, err
}

func (a *AppDb) UpdateProduct(product model.Product) error {
	q := `UPDATE product SET category_id=?, title=?, image_url=?, price=?, description=?, updated_at=NOW() WHERE id=?`
	_, err := a.Conn.Exec(q, product.CategoryId, product.Title, product.ImageUrl, product.Price, product.Description, product.Id)
	if err != nil {
		return err
	}
	return nil
}

func (a *AppDb) DeleteProduct(id string) error {
	q := `DELETE FROM product WHERE id=?`
	_, err := a.Conn.Exec(q, id)
	if err != nil {
		return err
	}
	return nil
}

func (a *AppDb) GetAllCategories(offset int, limit int, orderBy string, asc bool) ([]model.Category, error) {
	var categories []model.Category
	q := "SELECT * FROM category"
	if len(orderBy) > 0 {
		if !valid.MatchString(orderBy) {
			return categories, &SqlInjectionAttemptError{}
		}
		sort := "desc"
		if asc == true {
			sort = "asc"
		}
		q += fmt.Sprintf(` ORDER BY %s %s`, orderBy, sort)
	}
	if limit != 0 {
		q += fmt.Sprintf(` LIMIT %d OFFSET %d `, limit, offset)
	}
	rows, err := a.Conn.Queryx(q)
	if err != nil {
		return categories, err
	}
	for rows.Next() {
		var cat model.Category
		err = rows.StructScan(&cat)
		if err != nil {
			return categories, err
		}
		categories = append(categories, cat)
	}
	return categories, nil
}

func (a *AppDb) GetCategory(id int) (model.Category, error) {
	var category model.Category
	err := a.Conn.QueryRowx("SELECT * FROM category WHERE id=?", id).StructScan(&category)
	if err != nil {
		return model.Category{}, err
	}
	return category, nil
}

func (a *AppDb) AddCategory(category model.Category) error {
	q := `INSERT INTO category (title, pos, image_url) VALUES (?,?,?);`
	_, err := a.Conn.Exec(q, category.Title, category.Position, category.ImageUrl)
	if err != nil {
		return err
	}
	return nil
}

func (a *AppDb) UpdateCategory(category model.Category) error {
	q := `UPDATE category SET title=?, pos=?, image_url=?, updated_at=NOW() WHERE id=?`
	_, err := a.Conn.Exec(q, category.Title, category.Position, category.ImageUrl, category.Id)
	if err != nil {
		return err
	}
	return nil
}

func (a *AppDb) DeleteCategory(id int) error {
	q := `DELETE FROM category WHERE id=?`
	_, err := a.Conn.Exec(q, id)
	if err != nil {
		return err
	}
	return nil
}

func (a *AppDb) AddUser(user model.User) error {
	q := `INSERT INTO user (username, password) VALUES (?,?);`
	_, err := a.Conn.Exec(q, user.Username, passwd.Hash([]byte(user.Password)))
	if err != nil {
		return err
	}
	return nil
}

func (a *AppDb) UserExists(username string, password string) bool {
	var user model.User
	err := a.Conn.QueryRowx("SELECT * FROM user WHERE username=?", username).StructScan(&user)
	if err != nil {
		return false
	}
	return passwd.Authenticate(user.Password, []byte(password))
}

type SqlInjectionAttemptError struct{}

func (s *SqlInjectionAttemptError) Error() string {
	return "value contains malicious chars for sql injection"
}
