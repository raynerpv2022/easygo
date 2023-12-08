package controller

import (
	"context"
	"encoding/json"
	"errors"
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

func updateOneData(idO string, podcast data.Podcast) (int64, error) {
	id, err := primitive.ObjectIDFromHex(idO)
	if err != nil {
		return 0, fmt.Errorf("ID is not valid format, try again")
	}
	filter := bson.M{"_id": id}

	//replaceOne method ,$set is not used, if  a field is not present, replace  will not include it.
	// podcastReplace := bson.M{"name": podcast.Name, "updateTime": podcast.UpdateTime}
	// updateResult, err := collection.ReplaceOne(context.Background(), filter, podcastReplace)

	// UpdateOne update data filter by filter variable, data by $set.
	podcastUpdateOne := bson.M{"$set": bson.M{"name": podcast.Name, "title": podcast.Title, "updateTime": podcast.UpdateTime}}
	updateResult, err := collection.UpdateOne(context.Background(), filter, podcastUpdateOne)

	if err != nil {
		return 0, fmt.Errorf("error, ufffff try again %v ", err.Error())
	}

	return updateResult.ModifiedCount, nil

}

func deleteAllData() (int64, error) {
	// delete all data bson.M{} as filter
	deleteResult, err := collection.DeleteMany(context.Background(), bson.M{})
	if err != nil {
		//  if somethings was wrong set error message,
		return deleteResult.DeletedCount, fmt.Errorf("%v", err.Error())
	}
	return deleteResult.DeletedCount, nil
}

func deleteOneData(idUser string) (bson.M, int64, error) { // criteria to delete id or name? , is just an example, is not the best way...

	// check for id primitive format
	id, err := primitive.ObjectIDFromHex(idUser)

	if err != nil {
		// if not valid? set error
		return nil, 0, fmt.Errorf("no valid ID, try again, %v ", err.Error())
	}
	// get data to delete because  i want to return it... check how to get delete data from collection.deleteOne
	toDeleteData, errIsfound := getOneDatabyId(idUser)
	if errIsfound != nil {
		// if no found data set error
		return nil, 0, fmt.Errorf("no data found: %v", errIsfound)
	}
	// delete data
	deleteResult, err := collection.DeleteOne(context.Background(), bson.M{"_id": id})

	if err != nil {
		// if something happend, set error
		return nil, 0, fmt.Errorf("error, %v ", err)
	}
	// if no one deleted set error.
	if deleteResult.DeletedCount == 0 {
		return nil, 0, fmt.Errorf("no Data deleted,  try again, ")
	}
	// return  data deleted, number of data deleted, an d nil for error
	return toDeleteData, deleteResult.DeletedCount, nil

}
func getOneDatabyId(idOb string) (bson.M, error) { //used only to delete one bacause i want to show delete data
	// check for id primitive format
	id, err := primitive.ObjectIDFromHex(idOb)
	if err != nil {
		// set error if no valid primitive id
		return nil, errors.New("no valid ID, try again")
	}
	// var to set data if found
	var oneID bson.M
	// look for data
	err = collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&oneID)
	if err != nil {
		// set error if something happend
		return nil, errors.New(err.Error())
	}
	// if ok return dada found
	return oneID, nil

}

func getOneDataByName(name string) ([]bson.M, error) {

	var searchResult []bson.M
	// err := collection.FindOne(context.Background(), bson.M{"name": name}).Decode(&singleResult) // get only one data with filter
	cur, err := collection.Find(context.Background(), bson.M{"name": name}) // get all data with filter

	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		var data bson.M
		if err := cur.Decode(&data); err != nil {
			return nil, err
		}

		searchResult = append(searchResult, data)
	}

	if len(searchResult) == 0 {
		err = mongo.ErrNoDocuments
	}
	return searchResult, err

}

func insertOne(data data.Podcast) *mongo.InsertOneResult {

	idInserted, err := collection.InsertOne(context.Background(), data)
	if err != nil {
		log.Fatal(err)
	}
	return idInserted

}

func getAllData() []bson.D {
	cFindAll, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		fmt.Println("Error", err)
		log.Fatal(err)
	}
	defer cFindAll.Close(context.Background())
	var allData []bson.D
	for cFindAll.Next(context.Background()) {
		var potcats bson.D
		if err := cFindAll.Decode(&potcats); err != nil {

			log.Fatal(err)
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
	updatePodcast.UpdateTime = time.Now()

	updateResult, err := updateOneData(id["id"], updatePodcast)
	json.NewEncoder(w).Encode(fmt.Sprintf("Items updated : %v, Error : %v", updateResult, err))

}

func DeleteAllData(w http.ResponseWriter, r *http.Request) {
	fmt.Println("DELETE ALL")
	deleteResult, err := deleteAllData()
	if err != nil {
		json.NewEncoder(w).Encode(fmt.Sprintln(err))
		return
	}
	json.NewEncoder(w).Encode(fmt.Sprintf("Delete All Data succesfully, %v items deleted", deleteResult))

}

func DeleteOne(w http.ResponseWriter, r *http.Request) {
	fmt.Println("DELETE ONE")
	idDelete := mux.Vars(r)
	// if idDelete["id"] == "" {
	// 	json.NewEncoder(w).Encode("Id or name is not valid, please try again...")
	// 	return
	// }
	dataDeleted, count, err := deleteOneData(idDelete["id"])

	if err != nil {

		json.NewEncoder(w).Encode(fmt.Sprint(err))
		return
	}
	json.NewEncoder(w).Encode("Deleted Succesfully .." + fmt.Sprintf("%v element deleted", count))
	json.NewEncoder(w).Encode(fmt.Sprintf("%v ", dataDeleted))

}

func GetOneData(w http.ResponseWriter, r *http.Request) {

	param := mux.Vars(r)

	oneData, err := getOneDataByName(param["name"])

	fmt.Println("Get One Data by name")

	if err == mongo.ErrNoDocuments {
		json.NewEncoder(w).Encode(err) //"No data found"
		fmt.Println(err)
		return
	}

	json.NewEncoder(w).Encode(oneData)

}

func InsertData(w http.ResponseWriter, r *http.Request) {
	fmt.Println("CREATE DATA")
	var podcast data.Podcast

	if err := json.NewDecoder(r.Body).Decode(&podcast); err == io.EOF {
		json.NewEncoder(w).Encode("Body Empty, Please try again ...")
		return

	}
	if podcast.IsNameEmpty() {
		json.NewEncoder(w).Encode("Data is Empty, Please try again ...")
		return

	}
	podcast.CreateTime = time.Now()
	podcast.UpdateTime = podcast.CreateTime
	id_inserted := insertOne(podcast)
	all := getAllData()
	json.NewEncoder(w).Encode(id_inserted.InsertedID)
	json.NewEncoder(w).Encode(" ")
	json.NewEncoder(w).Encode(all)

}

func GetAllData(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET all data")
	allData := getAllData()

	json.NewEncoder(w).Encode(allData)

}
