package main

import (
	"Week3/controllers"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	// testing linter 9

	// endpoint login
	router.HandleFunc("/login", controllers.Login).Methods("POST")

	// endpoint get
	router.HandleFunc("/users", controllers.GetAllUsers).Methods("GET")
	router.HandleFunc("/products", controllers.GetAllProduct).Methods("GET")
	router.HandleFunc("/transactions", controllers.GetAllTransactions).Methods("GET")
	router.HandleFunc("/detailTransactions", controllers.GetDetailUserTransaction).Methods(("GET"))

	// endpoint insert
	router.HandleFunc("/users", controllers.InsertUser).Methods("POST")
	router.HandleFunc("/products", controllers.InsertProduct).Methods("POST")
	router.HandleFunc("/transactions", controllers.InsertTransaction).Methods("POST")

	// endpoint update
	router.HandleFunc("/users", controllers.UpdateUser).Methods("PUT")
	router.HandleFunc("/products", controllers.UpdateProduct).Methods("PUT")
	router.HandleFunc("/transactions", controllers.UpdateTransaction).Methods("PUT")

	// endpoint delete
	router.HandleFunc("/users", controllers.DeleteUser).Methods("DELETE")
	router.HandleFunc("/products", controllers.DeleteProduct).Methods("DELETE")
	router.HandleFunc("/transactions", controllers.DeleteTransaction).Methods("DELETE")

	// endpoint gorm
	router.HandleFunc("/v2/users", controllers.GetAllUsersGorm).Methods("GET")
	router.HandleFunc("/v2/users", controllers.InsertUserGorm).Methods("POST")
	router.HandleFunc("/v2/users", controllers.UpdateUserGorm).Methods("PUT")
	router.HandleFunc("/v2/users", controllers.DeleteUserGorm).Methods("DELETE")
	router.HandleFunc("/v2/detailTransactions", controllers.GetDetailUserTransactionGorm).Methods("GET")

	http.Handle("/", router)
	fmt.Println("Connected to port 8888")
	log.Println("Connected to port 8888")
	log.Fatal(http.ListenAndServe(":8888", router))
}
