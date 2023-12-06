package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/raynerpv2022/mongodb/router"
)

func main() {
	r := router.Router()
	fmt.Println("Server UP and listening on port 4000 ... ")
	log.Fatal(http.ListenAndServe(":4000", r)) 
	

}
