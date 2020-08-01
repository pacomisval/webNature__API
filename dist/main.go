package main

import (
	"database/sql"
	"encoding/json"
	"fmt"	
	"log"
	"net/http"
	//"reflect"
	"sync"
	"strconv"
	"time"
	_"github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"gopkg.in/guregu/null.v3"
	"github.com/dgrijalva/jwt-go"

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

type Value struct {
	Token  string `json:"token"`
}

type Claims struct {
	Name    string
	Secret  string
	*jwt.StandardClaims
}

////////////////////////////// SINGLETON /////////////////////////////////////////
// Creamos una variable única para toda la aplicación, sirve de identificación. //
type Singleton struct {
	Tiempo int64
}

var instancia *Singleton
var once sync.Once
/////////////////////////////////////////////////

var db *sql.DB
var err error

var v string


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

	router.HandleFunc("/init", getCookieT).Methods("GET")

	log.Print("Server started on localhost:8000 .........")

	log.Fatal(http.ListenAndServe(":8000", handlers.CORS(
		handlers.AllowCredentials(),
		handlers.AllowedHeaders([]string{"X-Requested-Widt", "Content-Type", "Authorization", "Accept", "Accept-Language"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "DELETE", "OPTIONS"}),
		handlers.AllowedOrigins([]string{"http://localhost:8088"}))(router)))

}

func GetInstancia() *Singleton {
	once.Do(func() {
		instancia = &Singleton {
			time.Now().Unix(),
		}
	})

	return instancia
}

//////////////////// ENDPONITS CAROUSEL ///////////////////////
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
///////////////////////////// FIN CAROUSEL ////////////////////////
/////////////////////// ENDPOINTS CATEGORIAS /////////////////////////
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
///////////////////// FIN CATEGORIAS /////////////////////////////
///////////////////// ENDPOINTS PRODUCTOS //////////////////////////
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

///////////////////////////// FIN PRODUCTOS /////////////////////////

////////////////////// ENDPOINT COOKIET ///////////////////////

/**
* Crea identificador único para la plicación.
* Comprobamos que la cookie existe.
* Si no existe, la creamos.
* Si existe, la refrescamos.
*/
func getCookieT(w http.ResponseWriter, r *http.Request){     
	// GetInstancia nos devuelve una instancia unica global para toda la aplicación.
	// convertimos tipo int64 a string
	vAux64 := GetInstancia().Tiempo     // vAux64 tipo int64
	v = strconv.FormatInt(vAux64, 10)   //  tipo string

	// GUARDAR EN LA BASE DE DATOS EL TOKEN //


	// RECUPERAR EL TOKEN DE LA BASE DE DATOS //

	cookie := verificarCookies(r)

	var value Value            // value tipo struct Value
	value = Value{v}  // conversion de v tipo string a value tipo struct Value

	if cookie != 0 {  // si no existe la cookieT, la crea.
		crearCookie(w, r, value)
	} else {         // si existe la cookieT, la refresca o actualiza.
		eliminarCookie(w, r)
		crearCookie(w, r, value) // EL VALUE TIENE QUE SER DE LA BASE DE DATOS: "valueDB"
		fmt.Println("Se ha actualizado la cookie")
	}	
}

func crearCookie(w http.ResponseWriter, r *http.Request, value Value) {
	T := value.Token

	expiration := time.Now().Add(time.Minute * 10)

	cookie1 := &http.Cookie {
		Name:     "natureT",
		Value:    T,
		Path:     "/",
		Expires:  expiration,
		HttpOnly: false, 
		Secure:   false,
	}
	http.SetCookie(w, cookie1)
	r.AddCookie(cookie1)

	fmt.Println("valor de cookie1: " , cookie1)
	json.NewEncoder(w).Encode(cookie1)
}

func eliminarCookie(w http.ResponseWriter, r *http.Request) {
	expiration := time.Now().Add(time.Minute - 1)

	cookie1 := &http.Cookie {
		Name:    "natureT",
		Value:   "",
		Path:    "/",
		Expires: expiration,
	}
	http.SetCookie(w, cookie1)
	r.AddCookie(cookie1)

	fmt.Println("valor de cookie1 despues de ser eliminada: ", cookie1)
}

/**
* Comprueba si existe la cookieT
* Si devuelve 1 es que no existe
* Si devuelve 0 si existe.
*/
func verificarCookies(r *http.Request) int {
	t, err := r.Cookie("natureT")
	if err != nil {
		fmt.Println("ERROR cookie natureT: ", err)
		if err == http.ErrNoCookie {
			//w.WriteHeader(http.StatusUnauthorized)
			fmt.Println("Error de StatusUnauthorized:", err)
			return 1
		}
		return 1
	}
	tknStr := t.Value  // valor de la cookie natureT
	fmt.Println("Valor de tknStr: ", tknStr)
	return 0
}

////////////////////// FIN ENDPOINT COOKIET ////////////////////


/* func crearToken(name, secret string) string {
	expiresAt := time.Now().Add(time.Minute * 130).Unix()

	claims := Claims {
		Name: name,
		Secret: secret, 
		StandardClaims: &jwt.StandardClaims {
			ExpiresAt: expiresAt,
		},
	}

	Secret := []byte(secret)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(Secret)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("valor de tokenString: ", tokenString)
	return tokenString
}
 */

