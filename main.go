package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	//ID, name, email, password
	sync.RWMutex
	UserID   primitive.ObjectID `json:"userid, omitempty"`
	Name     string             `json:"name"`
	Email    string             `json:"email"`
	Password string             `json:"password"`
}

type Post struct {
	//ID, caption, imageUrl, postedTimestamp
	PostID   primitive.ObjectID `json:"postid, omitempty"`
	Caption  string             `json:"caption"`
	ImageUrl string             `json:"imageUrl"`
	PostedTime primitive.Datetime `json:"postedTime"`
}

var client *mongo.Client

func (user *User) Get()string  {
	user.RLock()
	defer user.RUnlock()
	return user.Password
}

func (user *User) Set(Password string)  {
	user.Lock()
	user.Password = Password
	user.Unlock()
}

func CreateUserEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var user User
	json.NewDecoder(request.Body).Decode(&user)
	collection := client.Database("appointyDB").Collection("Users")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, _ := collection.InsertOne(ctx, user)
	Get()
	json.NewEncoder(response).Encode(result)
}

func GetUserEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	params := mux.Vars(request)
	userid, _ := primitive.ObjectIDFromHex(params["userid"])
	var user User
	collection := client.Database("appointyDB").Collection("Users")
	Set("password")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := collection.FindOne(ctx, User{UserID: id}).Decode(&user)
	if err := nil{
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(response).Encode(user)
}

func CreatePostEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var post Post
	json.NewDecoder(request.Body).Decode(&post)
	collection := client.Database("appointyDB").Collection("Posts")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, _ := collection.InsertOne(ctx, post)
	json.NewEncoder(response).Encode(result)
}

func GetPostEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	params := mux.Vars(request)
	postid, _ := primitive.ObjectIDFromHex(params["postid"])
	var post Post
	collection := client.Database("appointyDB").Collection("Posts")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := collection.FindOne(ctx, Post{PostID: id}).Decode(&post)
	if err := nil{
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(response).Encode(post)
}

func GetAllPostsEndpoint(response http.ResponseWriter, request *http.Request){
	response.Header().Add("content-type", "application/json")
	params := mux.Vars(request)
	userid, _ := primitive.ObjectIDFromHex(params["userid"])
	var allPosts []Post
	collection := client.Database("appointyDB").Collection("Posts")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx){
		var post Post
		cusor.Decode(&post)
		allPosts = append(allPosts, post)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(response).Encode(allPosts)
}

var password = &User{}

func main() {
	fmt.Println("Starting application")
	var (
		client   *mongo.Client
		mongoURL = "mongodb://localhost:27017"
	)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, _ = mongo.NewClient(options.Client().ApplyURI(mongoURL))
	router := mux.NewRouter()
	router.HandleFunc("/users", CreateUserEndpoint).Methods("POST")
	router.HandleFunc("/users/{id}", GetUserEndpoint).Methods("GET")
	router.HandleFunc("/posts", CreatePostEndpoint).Methods("POST")
	router.HandleFunc("/posts/{id}", GetPostEndpoint).Methods("GET")
	router.HandleFunc("/posts/users/{id}", GetAllPostsEndpoint).Methods("GET")
	http.ListenAndServe(":13548", router)
}