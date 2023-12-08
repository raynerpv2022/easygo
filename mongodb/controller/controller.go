package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/raynerpv2022/mongodb/data"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const URI = "mongodb+srv://resba:resba.net.2024@cluster0.awuzw0u.mongodb.net/?retryWrites=true&w=majority"

var collection *mongo.Collection

func init() {
	contextTemp := context.TODO()
	clientOP := options.Client().ApplyURI(URI)
	client, err := mongo.Connect(contextTemp, clientOP)
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Ping(contextTemp, nil); err != nil {
		log.Fatal(err)
	}

	// defer client.Disconnect(contextTemp)   Where to close it....???
	fmt.Println("MongoDb conection was succesfully...")
	db := "nextflix"
	col := "watchlist"
	collection = client.Database(db).Collection(col)
	fmt.Println("Colecttion is Ready ...")

}

// Helper methods

func updateOneData(idO string, podcast data.Podcast) int64 {
	id, err := primitive.ObjectIDFromHex(idO)
	if err != nil {
		fmt.Println("ID is not valid format, try again")
		return 0
	}
	filter := bson.M{"_id": id}

	//replaceOne method ,$set is not used, if  a field is not present, replace  will not include it.
	// podcastReplace := bson.M{"name": podcast.Name, "updateTime": podcast.UpdateTime}
	// updateResult, err := collection.ReplaceOne(context.Background(), filter, podcastReplace)

	// UpdateOne update data filter by filter variable, data by $set.
	podcastUpdateOne := bson.M{"$set": bson.M{"name": podcast.Name, "title": podcast.Title, "updateTime": podcast.UpdateTime}}
	updateResult, err := collection.UpdateOne(context.Background(), filter, podcastUpdateOne)

	if err != nil {
		fmt.Println(err.Error())
		return 0
	}

	return updateResult.ModifiedCount

}

func deleteAllData() int64 {
	// delete all data bson.M{} as filter
	deleteResult, err := collection.DeleteMany(context.Background(), bson.M{})
	if err != nil {
		//  if somethings was wrong set error message,
		fmt.Printf("%v", err.Error())
		return 0
	}
	return deleteResult.DeletedCount
}

func deleteOneData(idUser string) int64 { // criteria to delete id or name? , is just an example, is not the best way...

	// check for id primitive format
	id, err := primitive.ObjectIDFromHex(idUser)

	if err != nil {
		// if not valid? set error
		fmt.Println("no valid ID, try again, ", err.Error())
		return 0
	}
	// get data to delete because  i want to return it... check how to get delete data from collection.deleteOne

	// delete data
	deleteResult, err := collection.DeleteOne(context.Background(), bson.M{"_id": id})

	if err != nil {
		// if something happend, set error
		fmt.Println(err)
		return 0
	}
	// if no one deleted set error.
	if deleteResult.DeletedCount == 0 {
		fmt.Println("no Data found in DB, try again ... ")
		return 0
	}
	// return  data deleted, number of data deleted, an d nil for error
	return deleteResult.DeletedCount

}
func getOneDatabyId(idOb string) bson.M { //used only to delete one bacause i want to show delete data
	// check for id primitive format
	id, err := primitive.ObjectIDFromHex(idOb)
	if err != nil {
		// set error if no valid primitive id

		fmt.Println("Error , no valid ID, try again")
		return nil
	}
	// var to set data if found
	var oneID bson.M
	// look for data
	err = collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&oneID)
	if err != nil {
		// set error if something happend
		fmt.Println("Error :", err)
		return nil
	}
	// if ok return dada found
	return oneID

}

func getOneDataByName(name string) []bson.M {

	var searchResult []bson.M
	// err := collection.FindOne(context.Background(), bson.M{"name": name}).Decode(&singleResult) // get only one data with filter
	cur, err := collection.Find(context.Background(), bson.M{"name": name}) // get all data with filter

	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		var data bson.M
		if err := cur.Decode(&data); err != nil {
			fmt.Println(err)
			return nil
		}

		searchResult = append(searchResult, data)
	}

	if len(searchResult) == 0 {
		fmt.Println("No data found in DB")
		return searchResult
	}

	return searchResult

}

func insertOne(data data.Podcast) *mongo.InsertOneResult {

	idInserted, err := collection.InsertOne(context.Background(), data)
	if err != nil {
		fmt.Println("Error", err)

	}
	return idInserted

}

func getAllData() []bson.D {
	cFindAll, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		fmt.Println("Error", err)
		return nil
	}
	defer cFindAll.Close(context.Background())
	var allData []bson.D
	for cFindAll.Next(context.Background()) {
		var potcats bson.D
		if err := cFindAll.Decode(&potcats); err != nil {

			fmt.Println("Error", err)
			return nil
		}
		allData = append(allData, potcats)
	}

	return allData

}

// controler Methods

func UpdateOneData(w http.ResponseWriter, r *http.Request) {
	fmt.Println("UPDATE DATA")
	var updatePodcast data.Podcast
	id := mux.Vars(r)

	if err := json.NewDecoder(r.Body).Decode(&updatePodcast); err == io.EOF {
		json.NewEncoder(w).Encode("Data is empty, try again")
		return
	}
	if updatePodcast.IsNameEmpty() {
		fmt.Println("Data is not correctly, Please try again ...")
		json.NewEncoder(w).Encode("Data is not correctly, Please try again ...")
		return

	}
	updatePodcast.UpdateTime = time.Now()

	updateResult := updateOneData(id["id"], updatePodcast)
	json.NewEncoder(w).Encode(fmt.Sprintf("Items updated : %v", updateResult))
	fmt.Println("Items updated : ", updateResult)

}

func DeleteAllData(w http.ResponseWriter, r *http.Request) {
	fmt.Println("DELETE ALL")
	deleteResult := deleteAllData()

	json.NewEncoder(w).Encode(fmt.Sprintf("%v items deleted", deleteResult))
	fmt.Printf("%v element deleted\n", deleteResult)

}

func DeleteOne(w http.ResponseWriter, r *http.Request) {
	fmt.Println("DELETE ONE")
	idDelete := mux.Vars(r)

	count := deleteOneData(idDelete["id"])

	json.NewEncoder(w).Encode(fmt.Sprintf("%v element deleted", count))
	fmt.Printf("%v element deleted\n", count)

}

// Fix two function in One
func GetOneDatabyName(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get One Data by name")
	param := mux.Vars(r)

	oneData := getOneDataByName(param["name"])
	json.NewEncoder(w).Encode(oneData)

}

func GetOneDatabyId(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get One Data by ID")
	param := mux.Vars(r)

	oneData := getOneDatabyId(param["id"])
	json.NewEncoder(w).Encode(oneData)

}
func InsertData(w http.ResponseWriter, r *http.Request) {
	fmt.Println("CREATE DATA")
	var podcast data.Podcast

	if err := json.NewDecoder(r.Body).Decode(&podcast); err == io.EOF {
		fmt.Println("Body Empty, Please try again ...")
		json.NewEncoder(w).Encode("Body Empty, Please try again ...")
		return

	}
	if podcast.IsNameEmpty() {
		fmt.Println("Data is not correctly, Please try again ...")
		json.NewEncoder(w).Encode("Data is not correctly, Please try again ...")
		return

	}
	podcast.CreateTime = time.Now()
	podcast.UpdateTime = podcast.CreateTime
	_ = insertOne(podcast)

	json.NewEncoder(w).Encode(podcast)

}

func GetAllData(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET all data")
	allData := getAllData()

	json.NewEncoder(w).Encode(allData)

}
