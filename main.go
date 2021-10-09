package main

import (
	"context"
	"fmt"
	"time"
	"net/http"
	"encoding/json"
	
	
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	//"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct{
	//ID, name, email, password
	UserID primitive.ObjectID `json:"_userid, omitempty"`
	Name string `json:"name"`
	Email string `json:"email"`
	Password string `json:"password"`
}

type Post struct{
	//ID, caption, imageUrl, postedTimestamp
	PostID primitive.ObjectID `json:"_postid"`
	Caption string `json:"caption"`
	ImageUrl string `json:"imageUrl"`
}

var client *mongo.Client

func CreateUserEndpoint(response http.ResponseWriter, request *http.Request){
	response.Header().Add("content-type", "application/json")
	var user User
	json.NewDecoder (request.Body).Decode(&user)
	collection := client.Database("appointyDB").Collection("Users")
	ctx,_:= context.WithTimeout(context.Background(), 10*time.Second)
	result,_ := collection.InsertOne(ctx, user)
	json.NewEncoder(response).Encode(result)
}

func main(){
	fmt.Println("Starting application") 
	var (
		client     *mongo.Client
		mongoURL = "mongodb://localhost:27017"
	)
	ctx,_:= context.WithTimeout(context.Background(), 10*time.Second)
	client,_ = mongo.NewClient(options.Client().ApplyURI(mongoURL))
	router := mux.NewRouter()
	router.HandleFunc("/users", CreateUserEndpoint).Methods("POST")
	// router.HandleFunc("/users/{id}", GetUserEndpoint).Methods("GET")
	// router.HandleFunc("/posts", CreatePostEndpoint).Methods("POST")
	// router.HandleFunc("/posts/{id}", GetPostEndpoint).Methods("GET")
	// router.HandleFunc("/posts/users/{id}", GetAllPostsEndpoint).Methods("GET")
	http.ListenAndServe(":12345", router)
}