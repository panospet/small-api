package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	"github.com/panospet/small-api/pkg/cache"
	"github.com/panospet/small-api/pkg/model"
	"github.com/panospet/small-api/pkg/services"
)

type Api struct {
	Db    services.DbService
	Cache cache.Cacher
}

func NewApi(db services.DbService, cache cache.Cacher) *Api {
	return &Api{
		Db:    db,
		Cache: cache,
	}
}

func Authenticator(nextHandler http.HandlerFunc, app *Api) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authenticated := false
		if username, password, ok := r.BasicAuth(); ok {
			authenticated = app.Db.UserExists(username, password)
		}
		if !authenticated {
			respondWithError(w, http.StatusUnauthorized, "Authorization failed")
			return
		}
		nextHandler(w, r)
	}
}

func (a *Api) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/", Authenticator(a.health, a))

	// products
	router.HandleFunc("/v1/products", a.getAllProducts).Methods("GET")
	router.HandleFunc("/v1/products/{id}", a.getProduct).Methods("GET")
	router.HandleFunc("/v1/products", Authenticator(a.createProduct, a)).Methods("POST")
	router.HandleFunc("/v1/products/{id}", Authenticator(a.updateProduct, a)).Methods("PATCH")
	router.HandleFunc("/v1/products/{id}", Authenticator(a.deleteProduct, a)).Methods("DELETE")

	// categories
	router.HandleFunc("/v1/categories", a.getAllCategories).Methods("GET")
	router.HandleFunc("/v1/categories/{id}", a.getCategory).Methods("GET")
	router.HandleFunc("/v1/categories", Authenticator(a.createCategory, a)).Methods("POST")
	router.HandleFunc("/v1/categories/{id}", Authenticator(a.updateCategory, a)).Methods("PATCH")
	router.HandleFunc("/v1/categories/{id}", Authenticator(a.deleteCategory, a)).Methods("DELETE")

	log.Println("API is starting...")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		panic(fmt.Sprintf("API error: %s", err))
	}
}

func (a *Api) getAllProducts(w http.ResponseWriter, r *http.Request) {
	orderByValue := r.FormValue("orderBy")
	var asc bool
	var orderBy string
	if orderByValue != "" {
		parts := strings.Split(orderByValue, ":")
		// todo if len(parts) > 2 then throw error, if len(parts)==1 then it has to be only orderBy
		orderBy = parts[0]
		if len(parts) > 1 {
			if parts[1] == "asc" {
				asc = true
			} else if parts[1] == "desc" {
				asc = false
			} else {
				respondWithError(w, http.StatusBadRequest, "Bad order by value. Example \"orderBy=price:asc\"")
				return
			}
		}
	}
	p, err := getPaginationFromRequest(r)
	if err != nil {
		// todo specific error for p value
		respondWithError(w, http.StatusInternalServerError, "Error in pagination values")
		return
	}
	products, err := a.Db.GetAllProducts(p.offset, p.limit, orderBy, asc)
	total := len(products)
	if err != nil {
		log.Println("error while getting products", err)
		respondWithError(w, http.StatusInternalServerError, "Error while getting products")
		return
	}
	if p.start > total-1 {
		respondWithError(w, http.StatusBadRequest, "Page does not exist")
		return
	}
	end := p.end
	if p.end > total {
		end = total
	}
	setPaginationHeaders(w, r, p, total)
	respondWithJSON(w, http.StatusOK, products[p.start:end])
}

func (a *Api) getProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if cacheRes, err := a.Cache.GetProduct(id); err == nil && cacheRes != "" {
		fmt.Println("got product", id, "from cache")
		respondCachedWithJson(w, http.StatusOK, []byte(cacheRes))
		return
	} else if err != nil {
		log.Println("error getting from cache product with id", id, err)
	}
	product, err := a.Db.GetProduct(id)
	if err != nil {
		log.Println("error while getting product", err)
		respondWithError(w, http.StatusInternalServerError, "Error while getting product")
		return
	}
	respondWithJSON(w, http.StatusOK, product)
}

func (a *Api) createProduct(w http.ResponseWriter, r *http.Request) {
	var product model.Product
	raw, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid POST request")
		return
	}
	if err := json.Unmarshal(raw, &product); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid POST request")
		return
	}

	id, err := a.Db.AddProduct(product)
	if err != nil {
		log.Println("product could not be added", err)
		respondWithError(w, http.StatusInternalServerError, "Product could not be added")
		return
	}
	respondWithJSON(w, http.StatusCreated, Response{Message: fmt.Sprintf("Product with id %s was created", id)})
}

func (a *Api) updateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	product, err := a.Db.GetProduct(id)
	if err != nil {
		log.Println("error while getting product", err)
		respondWithError(w, http.StatusInternalServerError, "Error while getting product")
		return
	}
	raw, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid POST request")
		return
	}
	if err := json.Unmarshal(raw, &product); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid POST request")
		return
	}
	err = a.Db.UpdateProduct(product)
	if err != nil {
		log.Println("error while updating product", err)
		respondWithError(w, http.StatusInternalServerError, "Product could not be updated")
		return
	}
	respondWithJSON(w, http.StatusCreated, Response{Message: fmt.Sprintf("Product with id %s was updated", id)})
}

func (a *Api) deleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	err := a.Db.DeleteProduct(id)
	if err != nil {
		log.Println("error while deleting product", err)
		respondWithError(w, http.StatusInternalServerError, "Error while deleting product")
		return
	}
	respondWithJSON(w, http.StatusOK, Response{Message: fmt.Sprintf("Product with id %s was deleted", id)})
}

func (a *Api) getAllCategories(w http.ResponseWriter, r *http.Request) {
	orderByValue := r.FormValue("orderBy")
	var asc bool
	var orderBy string
	if orderByValue != "" {
		parts := strings.Split(orderByValue, ":")
		// todo if len(parts) > 2 error, if len(parts)==1 then it has to be only orderBy
		orderBy = parts[0]
		if len(parts) > 1 {
			if parts[1] == "asc" {
				asc = true
			} else if parts[1] == "desc" {
				asc = false
			} else {
				respondWithError(w, http.StatusBadRequest, "Bad order by value. Example \"orderBy=title:asc\"")
				return
			}
		}
	}
	p, err := getPaginationFromRequest(r)
	if err != nil {
		// todo specific error for p value
		respondWithError(w, http.StatusInternalServerError, "Error in pagination values")
		return
	}
	if orderBy == "position" {
		orderBy = "pos"
	}
	categories, err := a.Db.GetAllCategories(p.offset, p.limit, orderBy, asc)
	total := len(categories)
	if err != nil {
		if _, ok := err.(*services.SqlInjectionAttemptError); ok {
			respondWithError(w, http.StatusBadRequest, "Bad parameters given (I saw what you did there ;) )")
			return
		}
		log.Println("error while getting categories", err)
		respondWithError(w, http.StatusInternalServerError, "Error while getting categories")
		return
	}
	if p.start > total-1 {
		respondWithError(w, http.StatusBadRequest, "Page does not exist")
		return
	}
	end := p.end
	if p.end > total {
		end = total
	}
	setPaginationHeaders(w, r, p, total)
	respondWithJSON(w, http.StatusOK, categories[p.start:end])
}

func (a *Api) getCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println("error with category id", err)
		respondWithError(w, http.StatusBadRequest, "Bad category id")
		return
	}
	if cacheRes, err := a.Cache.GetCategory(vars["id"]); err == nil && cacheRes != "" {
		fmt.Println("got category", id, "from cache")
		respondCachedWithJson(w, http.StatusOK, []byte(cacheRes))
		return
	} else if err != nil {
		log.Println("error getting from cache category with id", id, err)
	}
	category, err := a.Db.GetCategory(id)
	if err != nil {
		log.Println("error while getting category", err)
		respondWithError(w, http.StatusInternalServerError, "Error while getting category")
		return
	}
	respondWithJSON(w, http.StatusOK, category)
}

func (a *Api) createCategory(w http.ResponseWriter, r *http.Request) {
	var category model.Category
	raw, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid POST request")
		return
	}
	if err := json.Unmarshal(raw, &category); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid POST request")
		return
	}
	err = a.Db.AddCategory(category)
	if err != nil {
		log.Println("product could not be added", err)
		respondWithError(w, http.StatusInternalServerError, "Category could not be added")
		return
	}
	respondWithJSON(w, http.StatusCreated, Response{Message: "Category was created successfully"})
}

func (a *Api) updateCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println("error with category id", err)
		respondWithError(w, http.StatusBadRequest, "Bad category id")
		return
	}
	category, err := a.Db.GetCategory(id)
	if err != nil {
		log.Println("error while getting category", err)
		respondWithError(w, http.StatusInternalServerError, "Error while getting category")
		return
	}
	raw, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid POST request")
		return
	}
	if err := json.Unmarshal(raw, &category); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid POST request")
		return
	}
	err = a.Db.UpdateCategory(category)
	if err != nil {
		log.Println("error while updating category", err)
		respondWithError(w, http.StatusInternalServerError, "Category could not be updated")
		return
	}
	respondWithJSON(w, http.StatusCreated, Response{Message: fmt.Sprintf("Category with id %d was updated", id)})
}

func (a *Api) deleteCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println("error with category id", err)
		respondWithError(w, http.StatusBadRequest, "Bad category id")
		return
	}
	err = a.Db.DeleteCategory(id)
	if err != nil {
		log.Println("error while deleting product", err)
		respondWithError(w, http.StatusInternalServerError, "Error while deleting category")
		return
	}
	respondWithJSON(w, http.StatusOK, Response{Message: fmt.Sprintf("Category with id %d was deleted", id)})
}

func (a *Api) health(w http.ResponseWriter, r *http.Request) {
	type Health struct {
		Message string
		Code    int
	}
	health := Health{
		Message: "health good!",
	}
	data, _ := json.Marshal(health)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusTeapot)
	w.Write(data)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondCachedWithJson(w http.ResponseWriter, code int, payload []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(payload)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		type ApiError struct {
			Code    int
			Message string
		}
		apiError := ErrorResponse{Code: 500, Message: "Internal server error",
			Err: err}
		response, _ = json.Marshal(apiError)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		w.Write(response)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

type Response struct {
	Message string `json:"message"`
}
