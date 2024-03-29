package main

import (
	"fmt"
	"net/http"

	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/shadowshot-x/micro-product-go/authservice"
)

func main() {
	mainRouter := mux.NewRouter()

	authRouter := mainRouter.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/signup", authservice.SignupHandler).Methods("POST")

	authRouter.HandleFunc("/signin", authservice.SigninHandler).Methods("POST")

	ch := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"http://localhost:8080"}))

	server := &http.Server{
		Addr:    "127.0.0.1:8000",
		Handler: ch(mainRouter),
	}
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Error Booting the Server")
	}
}
