package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/arthurbonini/micro/data"
	"github.com/gorilla/mux"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l*log.Logger) *Products {
	return &Products{l}
}


func (p*Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET Request")
	lp := data.GetProducts()
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (p*Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Request")

	prod := r.Context().Value(KeyProduct{}).(data.Product)
	data.AddProduct(&prod)
}

func (p*Products) UpdateProduct(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
	}

	p.l.Println("Handle PUT Request", id)
	prod := r.Context().Value(KeyProduct{}).(data.Product)

	
	err = data.UpdateProduct(id, &prod)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product Not Found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(rw, "Product Not Found", http.StatusInternalServerError)
		return
	}
}

type KeyProduct struct{}

func (p *Products) MiddlewareValidateProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prod := &data.Product{}

		err := prod.FromJSON(r.Body)
		if err != nil {
			http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		req := r.WithContext(ctx)
		next.ServeHTTP(rw, req)
	})
}