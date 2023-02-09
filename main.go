package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"os"

	"example.com/m/helper"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

var collection = helper.ConnectDB()

type Movie struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Isbn     string             `json:"isbn,omitempty" bson:"isbn,omitempty"`
	Title    string             `json:"title" bson:"title,omitempty"`
	Director *Director          `json:"director" bson:"director,omitempty"`
}
type Director struct {
	FirstName string `json:"firstname,omitempty" bson:"firstname,omitempty"`
	LastName  string `json:"lastname,omitempty" bson:"lastname,omitempty"`
}

var movies []Movie

// func getMovies(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("content-type", "application/json")
// 	json.NewEncoder(w).Encode(movies)
// }

// func deleteMovie(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("content-type", "application/json")
// 	params := mux.Vars(r)
// 	for index, item := range movies {
// 		if item.ID == params["id"] {
// 			movies = append(movies[:index], movies[index+1:]...)
// 			break
// 		}
// 	}
// 	json.NewEncoder(w).Encode(movies)
// }

// func creaetMovie(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	var newMovie Movie
// 	_ = json.NewDecoder(r.Body).Decode(&newMovie)
// 	newMovie.ID = strconv.Itoa(rand.Intn(100000000))
// 	movies = append(movies, newMovie)
// 	json.NewEncoder(w).Encode(newMovie)
// }
// func updateMovie(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	params := mux.Vars(r)

// 	for index, item := range movies {
// 		if item.ID == params["id"] {
// 			movies = append(movies[:index], movies[index+1:]...)
// 			var newMovie Movie
// 			_ = json.NewDecoder(r.Body).Decode(&newMovie)
// 			newMovie.ID = params["id"]
// 			movies = append(movies, newMovie)
// 			json.NewEncoder(w).Encode(newMovie)
// 			return
// 		}
// 	}
// }

func main() {
	err := godotenv.Load("local.env")
	if err != nil {
		log.Fatal("Error loading .env file %s", err)
	}
	r := mux.NewRouter()
	fileserver := http.FileServer(http.Dir("./static"))

	// movies = append(movies, Movie{ID: "1", Isbn: "438227", Title: "Movie 1", Director: &Director{Firstname: "John", Lastname: "Doe"}})
	// movies = append(movies, Movie{ID: "2", Isbn: "458227", Title: "Movie 2", Director: &Director{Firstname: "Ishan", Lastname: "Joe"}})
	// movies = append(movies, Movie{ID: "3", Isbn: "458245", Title: "Movie 3", Director: &Director{Firstname: "Steve", Lastname: "Smith"}})
	r.Handle("/", fileserver)
	r.HandleFunc("/movies", getMovies).Methods("GET")
	// r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movie/add", creaetMovie).Methods("POST")
	r.HandleFunc("/movie/update", updateMovie).Methods("POST")
	r.HandleFunc("/movie/delete", deleteMovie).Methods("POST")

	PORT := os.Getenv("PORT")
	fmt.Printf("Starting server at port %s\n", PORT)
	log.Fatal(http.ListenAndServe(PORT, r))
}

func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var movies []Movie
	cur, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		helper.GetError(err, w)
		return
	}
	defer cur.Close(context.TODO())
	for cur.Next(context.TODO()) {
		var movie Movie
		err := cur.Decode(&movie)
		if err != nil {
			log.Fatal(err)
		}
		movies = append(movies, movie)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(movies)

}

func creaetMovie(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "application/json")
	if err := r.ParseForm(); err != nil {
		log.Fatal(err)
		return
	}

	var movie Movie
	// _ = json.NewDecoder(r.Body).Decode(&movie)
	movie.Isbn = r.FormValue("isbn")
	movie.Title = r.FormValue("title")
	movie.Director = &Director{FirstName: r.FormValue("directorfname"), LastName: r.FormValue("directorlname")}

	_, err := collection.InsertOne(context.TODO(), movie)
	if err != nil {
		helper.GetError(err, w)
		return
	}
	// json.NewEncoder(w).Encode(result)
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Fatal(err)
		return
	}
	id, _ := primitive.ObjectIDFromHex(r.FormValue("id"))
	filter := bson.M{"_id": id}

	result, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		helper.GetError(err, w)
		return
	}
	json.NewEncoder(w).Encode(result)
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "application/json")
	// params := mux.Vars(r)
	if err := r.ParseForm(); err != nil {
		log.Fatal(err)
		return
	}
	id, _ := primitive.ObjectIDFromHex(r.FormValue("id"))
	// id, _ := primitive.ObjectIDFromHex(params["id"])
	// fmt.Println(params["id"])
	filter := bson.M{"_id": id}
	// fmt.Println(filter)
	var movie Movie
	// _ = json.NewDecoder(r.Body).Decode(&movie)
	movie.Isbn = r.FormValue("isbn")
	movie.Title = r.FormValue("title")
	movie.Director = &Director{FirstName: r.FormValue("directorfname"), LastName: r.FormValue("directorlname")}
	update := bson.M{"$set": movie}
	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		helper.GetError(err, w)
		return
	}
	json.NewEncoder(w).Encode(result)
}
