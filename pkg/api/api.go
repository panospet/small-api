package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/panospet/small-api/pkg/model"
	"github.com/panospet/small-api/pkg/services"
)

type Api struct {
	Db *services.AppDb
}

func NewApi(db *services.AppDb) *Api {
	return &Api{
		Db: db,
	}
}

func (a *Api) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/", a.health)

	// products
	router.HandleFunc("/v1/products", a.getAllProducts).Methods("GET")
	router.HandleFunc("/v1/products/{id}", a.getProduct).Methods("GET")
	router.HandleFunc("/v1/products", a.createProduct).Methods("POST")
	router.HandleFunc("/v1/products/{id}", a.updateProduct).Methods("PATCH")
	router.HandleFunc("/v1/products/{id}", a.deleteProduct).Methods("DELETE")

	// categories
	router.HandleFunc("/v1/categories", a.getAllCategories).Methods("GET")
	router.HandleFunc("/v1/categories/{id}", a.getCategory).Methods("GET")
	router.HandleFunc("/v1/categories", a.createCategory).Methods("POST")
	router.HandleFunc("/v1/categories/{id}", a.updateCategory).Methods("PATCH")
	router.HandleFunc("/v1/categories/{id}", a.deleteCategory).Methods("DELETE")

	log.Println("API is starting...")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		panic(fmt.Sprintf("API error: %s", err))
	}
}

func (a *Api) getAllProducts(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, []model.Product{})
}

func (a *Api) getProduct(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, model.Product{})
}

func (a *Api) createProduct(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusCreated, Response{Message: "product created"})
}

func (a *Api) updateProduct(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusCreated, Response{Message: "product updated"})
}

func (a *Api) deleteProduct(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, Response{Message: "product deleted"})
}

func (a *Api) getAllCategories(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, []model.Category{})
}

func (a *Api) getCategory(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, model.Category{})
}

func (a *Api) createCategory(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusCreated, Response{Message: "Category was created successfully"})
}

func (a *Api) updateCategory(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusCreated, Response{Message: "category updated"})
}

func (a *Api) deleteCategory(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, Response{Message: "category deleted"})
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
