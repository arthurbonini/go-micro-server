package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/arthurbonini/micro/data"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l*log.Logger) *Products {
	return &Products{l}
}

func (p*Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	//GET
	if r.Method == http.MethodGet {
		p.getProducts(rw, r)
		return
	}

	//ADD
	if r.Method == http.MethodPost {
		p.addProduct(rw, r)
		return
	}

	//UPDATE
	if r.Method == http.MethodPut {
		reg := regexp.MustCompile(`/([0-9]+)`)
		g := reg.FindAllStringSubmatch(r.URL.Path, 1)

		if len(g) != 1 {
			http.Error(rw, "Invalid URI", http.StatusNotFound)
			return
		}

		if len(g[0]) != 2 {
			http.Error(rw, "Invalid URI", http.StatusNotFound)
			return
		}

		idString := g[0][1]
		id, err := strconv.Atoi(idString)
		if err != nil {
			http.Error(rw, "Invalid URI", http.StatusNotFound)
			return
		}

		p.updateProduct(id, rw, r)
		return
	}

	//catchall
	rw.WriteHeader(http.StatusMethodNotAllowed)
}


func (p*Products) getProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET Request")
	lp := data.GetProducts()
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (p*Products) addProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Request")

	prod := &data.Product{}
	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	}

	data.AddProduct(prod)
}

func (p*Products) updateProduct(id int, rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle PUT Request")

	prod := &data.Product{}
	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	}

	err = data.UpdateProduct(id, prod)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product Not Found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(rw, "Product Not Found", http.StatusInternalServerError)
		return
	}
}