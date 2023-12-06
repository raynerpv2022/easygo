package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/raynerpv2022/apihttp/controller"
)

func main() {
	fmt.Println(" API HTTP Example")
	m := mux.NewRouter()

	fmt.Println("Server up...")
	m.HandleFunc("/data", controller.GetAllData).Methods("GET")
	m.HandleFunc("/data/{id}", controller.GetOne).Methods("GET")
	m.HandleFunc("/add", controller.AddOne).Methods("POST")
	m.HandleFunc("/update/{id}", controller.UpdateOne).Methods("PUT")
	m.HandleFunc("/delete/{id}", controller.DeleteOne).Methods("DELETE")
	m.HandleFunc("/deleteall", controller.DeleteAll).Methods("DELETE")
	fmt.Println("Listen at port 2000 ...")
	log.Fatal(http.ListenAndServe(":2000", m))

}
