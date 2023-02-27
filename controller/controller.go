package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/subhammahanty235/netflix-api-golang/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	// "go.mongodb.org/mongo-driver/x/mongo/driver/mongocrypt/options"
)

const dbname = "netflix"
const colname = "watchlist"

var collection *mongo.Collection

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	connectionURI := os.Getenv("MONGODB_URI")

	clientOption := options.Client().ApplyURI(connectionURI)

	client, err := mongo.Connect(context.TODO(), clientOption)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Mongodb Connection Success")
	collection = client.Database(dbname).Collection(colname)

	fmt.Println("Connection instance is ready")

}

// insert 1 recod

func insertOneMovie(movie model.Netflix) {
	inserted, err := collection.InsertOne(context.Background(), movie)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted one movie with the id", inserted.InsertedID)

}

func updateMovie(id string) {
	movieid, _ := primitive.ObjectIDFromHex(id)

	filter := bson.M{"_id": movieid}
	update := bson.M{"$set": bson.M{"watched": true}}

	result, err := collection.UpdateOne(context.Background(), filter, update)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Modified count", result.ModifiedCount)
}

func deleteMovie(id string) {
	movieid, _ := primitive.ObjectIDFromHex(id)

	filter := bson.M{"_id": movieid}

	result, err := collection.DeleteOne(context.Background(), filter)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Deleted movie ", result.DeletedCount)
}

func getAllMovies() []primitive.M {
	curser, err := collection.Find(context.Background(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}

	var movies []primitive.M
	for curser.Next(context.Background()) {
		var movie bson.M
		err := curser.Decode(&movie)
		if err != nil {
			log.Fatal(err)
		}
		movies = append(movies, movie)

	}
	defer curser.Close(context.Background())
	return movies

}

func getOneMovie(id string) model.Netflix {

	filter := bson.M{"_id": id}
	fmt.Println(filter)
	var movie model.Netflix
	fmt.Println("1")
	err := collection.FindOne(context.Background(), filter).Decode(&movie)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("2")
	fmt.Println(movie.Movie)
	return movie
}

// Actual Controllers

func GetAllMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	allMovies := getAllMovies()
	json.NewEncoder(w).Encode(allMovies)
}

func GetOneMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	movie := getOneMovie(params["id"])

	json.NewEncoder(w).Encode(movie)
}

func CreateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	var movie model.Netflix

	_ = json.NewDecoder(r.Body).Decode(&movie)
	insertOneMovie(movie)
	json.NewEncoder(w).Encode(movie)
}

func MarkAsWatched(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "PUT")

	params := mux.Vars(r)

	updateMovie(params["id"])
	json.NewEncoder(w).Encode(params)

}
func DeleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "DELETE")

	params := mux.Vars(r)
	deleteMovie(params["id"])

	json.NewEncoder(w).Encode(params)
}
