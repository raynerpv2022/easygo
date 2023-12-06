package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/raynerpv2022/apihttp/db"
)

var student []db.Student

// helper methods from controller
func isDBempty() bool {

	return len(student) == 0

}

func dbEmptyMessage(w http.ResponseWriter) bool {
	if isDBempty() {
		json.NewEncoder(w).Encode("DB is Empty, please enter data")
		return true
	}
	return false
}

func getBodyData(w http.ResponseWriter, r *http.Request) (db.Student, bool) {
	var local db.Student
	if errorEOF := json.NewDecoder(r.Body).Decode(&local); errorEOF == io.EOF {

		return local, true
		//if no data in body send error message
	}
	return local, false
}

// controller methods

func GetAllData(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get all data")
	if dbEmptyMessage(w) {
		return
	}
	json.NewEncoder(w).Encode(" ********  Data List *********")
	json.NewEncoder(w).Encode(student)

}

func GetOne(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Get one data")
	id := mux.Vars(r)

	var localStudent db.Student
	if dbEmptyMessage(w) {
		return
	}

	for _, v := range student {
		if v.Id == id["id"] {
			localStudent.Id = v.Id
			localStudent.Name = v.Name
			json.NewEncoder(w).Encode(localStudent)
			return
		}
	}
	json.NewEncoder(w).Encode("No data found in DB, try again...")

}

func AddOne(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Add one data")

	ramdonID := rand.Intn(100)

	//check if data is present in r.Body
	localStudent, err := getBodyData(w, r)
	if err {
		json.NewEncoder(w).Encode("NOT ADD allow, data is empty, please send again...")
		return
	}

	localStudent.Id = strconv.Itoa(ramdonID)
	student = append(student, localStudent)
	json.NewEncoder(w).Encode(" ADD is DONE succesfully...")
	json.NewEncoder(w).Encode(" ********  Data List ********* ")
	json.NewEncoder(w).Encode(student)

}

func UpdateOne(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update one data")
	if dbEmptyMessage(w) {
		return
	}
	id := mux.Vars(r)
	localStudent, err := getBodyData(w, r)
	if err {
		json.NewEncoder(w).Encode("NOT UPDATE allow, Data is empty, please send again...")
		return
	}

	for i, v := range student {
		if v.Id == id["id"] {
			json.NewEncoder(w).Encode(fmt.Sprintf("Data %v is ready to be updated ...", v))
			localStudent.Id = v.Id
			student = append(student[:i], student[i+1:]...)
			student = append(student, localStudent)
			json.NewEncoder(w).Encode("Update is DONE succesfully")
			json.NewEncoder(w).Encode(" ********  Data List ********* ")
			json.NewEncoder(w).Encode(student)
			return
		}
	}
	json.NewEncoder(w).Encode("No data found in DB to update, please try again...")

}

func DeleteAll(w http.ResponseWriter, r *http.Request) {
	if dbEmptyMessage(w) {
		return
	}
	fmt.Println("Delete all data {not working yet}")
	json.NewEncoder(w).Encode("All data will be deleted ...")
	student = []db.Student{}
	json.NewEncoder(w).Encode(student)
	json.NewEncoder(w).Encode("Delete All data DONE succesfully !!! ")

}

func DeleteOne(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete one data")
	if dbEmptyMessage(w) {
		return
	}
	id := mux.Vars(r)

	for i, v := range student {
		if v.Id == id["id"] {
			json.NewEncoder(w).Encode(fmt.Sprintf("Deleting %v ...", v))
			json.NewEncoder(w).Encode("Data deleted succesfully")
			student = append(student[:i], student[i+1:]...)
			json.NewEncoder(w).Encode(" ********  Data List *********")
			json.NewEncoder(w).Encode(student)
			return
		}
	}
	json.NewEncoder(w).Encode("No data found in DB to delete, please try again...")

}
