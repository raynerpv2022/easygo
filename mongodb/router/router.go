package router

import (
	"github.com/gorilla/mux"
	"github.com/raynerpv2022/mongodb/controller"
)

func Router() *mux.Router {
	m := mux.NewRouter()
	m.HandleFunc("/data", controller.GetAllData).Methods("GET")
	m.HandleFunc("/data/{name}", controller.GetOneData).Methods("GET")
	m.HandleFunc("/add", controller.InsertData).Methods("POST")
	m.HandleFunc("/delete/{id}", controller.DeleteOne).Methods("DELETE")
	m.HandleFunc("/deleteall", controller.DeleteAllData).Methods("DELETE")
	m.HandleFunc("/update/{id}", controller.UpdateOneData).Methods("PUT")
	return m
}
