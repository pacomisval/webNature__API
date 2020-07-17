package main

import (
	"database/sql"
	"encoding/json"
	"fmt"	
	"log"
	"net/http"
	//"reflect"
	_"github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"gopkg.in/guregu/null.v3"

)

type Carousel struct {
	Id           string `json:"id"`
	Nombre       string `json:"nombre"`
	Descripcion  string `json:"descripcion"`
	Imagen       string `json:"imagen"`
}

type Categoria struct {
	Id             string `json:"id"`
	Nombre         string `json:"nombre"`
	Descripcion    string `json:"descripcion"`
	FechaCreacion  string `json:"fechaCreacion"`
}

type Producto struct {
	Id			  string `json:"id"`
	Nombre		  string `json:"nombre"`
	Descripcion	  string `json:"descripcion"`
	Precio		  float64 `json:"precio"`
	Stock		  bool `json:"stock"`
	Oferta		  null.String `json:"oferta"`  // permite valores nulos
	IdCategoria	  string `json:"idCategoria"`
	Foto		  string `json:"foto"`
}

var db *sql.DB
var err error

func main() {
	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/nature")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	router := mux.NewRouter()

	router.HandleFunc("/carousel", getImagesCarousel).Methods("GET")

	router.HandleFunc("/categorias", getCategorias).Methods("GET")

	router.HandleFunc("/productos", getAllProductos).Methods("GET")
	router.HandleFunc("/productos/categoria/{id}", getProductosByIdCategoria).Methods("GET")

	log.Print("Server started on localhost:8000 .........")

	log.Fatal(http.ListenAndServe(":8000", handlers.CORS(
		handlers.AllowCredentials(),
		handlers.AllowedHeaders([]string{"X-Requested-Widt", "Content-Type", "Authorization", "Accept", "Accept-Language"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "DELETE", "OPTIONS"}),
		handlers.AllowedOrigins([]string{"*"}))(router)))

}

func getImagesCarousel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var images []Carousel

	result, err := db.Query("SELECT * FROM carousel")
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

	for result.Next() {
		var image Carousel
		err := result.Scan(&image.Id, &image.Nombre, &image.Descripcion, &image.Imagen)
		if err != nil {
			panic(err.Error())
		}
		images = append(images, image)
	}

	json.NewEncoder(w).Encode(images)

	fmt.Println("RESPONSE EN GET IMAGES CAROUSEL")
}

func getCategorias(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var categorias []Categoria

	result, err := db.Query("SELECT * FROM categorias")
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

	for result.Next() {
		var categoria Categoria
		err := result.Scan(&categoria.Id, &categoria.Nombre, &categoria.Descripcion, &categoria.FechaCreacion)
		if err != nil {
			panic(err.Error())
		}
		categorias = append(categorias, categoria)
	}

	json.NewEncoder(w).Encode(categorias)

	fmt.Println("RESPONSE EN GET CATEGORIAS")	
}

func getAllProductos(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	var productos []Producto

	result, err := db.Query("SELECT * FROM productos")
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

	for result.Next() {
		var producto Producto
		err := result.Scan(&producto.Id, &producto.Nombre, &producto.Descripcion, &producto.Precio, &producto.Stock, &producto.Oferta, &producto.IdCategoria, &producto.Foto)
		if err != nil {
			panic(err.Error())
		}
		
		productos = append(productos, producto)
	}

	json.NewEncoder(w).Encode(productos)

	fmt.Println("RESPONSE EN GET ALL PRODUCTOS")
}

func getProductosByIdCategoria(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	var productos []Producto

	params := mux.Vars(r)

	result, err := db.Query("SELECT * FROM productos WHERE idCategoria = ? ORDER BY nombre", params["id"])
	if err != nil {
		panic(err.Error())
	}

	defer result.Close()

	for result.Next() {
		var producto Producto
		err := result.Scan(&producto.Id, &producto.Nombre, &producto.Descripcion, &producto.Precio, &producto.Stock, &producto.Oferta, &producto.IdCategoria, &producto.Foto)
		if err != nil {
			panic(err.Error())
		}

		productos = append(productos, producto)
	}

	json.NewEncoder(w).Encode(productos)

	fmt.Println("RESPONSE EN GET PRODUCTO POR IDCATEGORIA")
}


