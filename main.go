package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"time"
	"net/http"
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type User struct {
	Name 		string `json:"name"`
	Email 		string `json:"email"`
	Password 	string `json:"password"`
}

type Post struct {
	Postowner 		string 		`json:"postowner"`
	Caption 		string 		`json:"caption"`
	ImageURL 		string 		`json:"imageURL"`
	Timestamp 		time.Time 	`json:"timestamp"`
}

func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil{
		panic(err)
	}
    return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

func users(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		postUser(w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
		return
	}
}

func posts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		createPost(w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
		return
	}
}

func postUser(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if(err != nil) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	var user User
	err = json.Unmarshal(bodyBytes, &user)
	if(err != nil) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}

	hash, err := HashPassword(user.Password)
	if !CheckPasswordHash(user.Password, hash) || err != nil {
		panic(err)
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	// Ping the primary
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected and pinged.")

	coll := client.Database("Instagram").Collection("users")
	newUser := bson.D{{"name", user.Name}, {"email", user.Email}, {"password", hash}}
	
	result, err := coll.InsertOne(context.TODO(), newUser)
	fmt.Printf("Inserted document with _id: %v\n", result.InsertedID)
}

func createPost(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if(err != nil) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	var post Post
	err = json.Unmarshal(bodyBytes, &post)
	if(err != nil) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}
	post.Timestamp = time.Now()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	// Ping the primary
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected and pinged.")

	coll := client.Database("Instagram").Collection("posts")
	newPost := bson.D{{"postowner", post.Postowner}, {"caption", post.Caption}, {"imageurl", post.ImageURL}, {"timestamp", post.Timestamp}}
	
	result, err := coll.InsertOne(context.TODO(), newPost)
	fmt.Printf("Inserted document with _id: %v\n", result.InsertedID)
}

func getPost(w http.ResponseWriter, r *http.Request) {
	postID, err := r.URL.Query()["id"]
    
    if !err || len(postID[0]) < 1 {
        panic(err)
        return
    }

    id := string(postID[0])
	fmt.Fprintln(w, "Url Param 'id' is: " + id)

	client, err2 := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	
	if err2 != nil {
		panic(err2)
	}

	defer func() {
		if err2 = client.Disconnect(context.TODO()); err2 != nil {
			panic(err2)
		}
	}()
	// Ping the primary
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected and pinged.")

	objectIDS, _ := primitive.ObjectIDFromHex(id)
	coll := client.Database("Instagram").Collection("posts")
	filter := bson.M{"_id": objectIDS}
	query, err3 := coll.Find(context.TODO(), filter)
	if err3 != nil {
		fmt.Println("errror retrieving post postid : " + id)
	}

	var findResults []bson.M
	if err3 = query.All(context.TODO(), &findResults); err3 != nil {
		panic(err3)
	}
	for _, result := range findResults {
		fmt.Fprintln(w, result)
	}
}

func getAllPostsOfUser(w http.ResponseWriter, r *http.Request) {
	userID, err := r.URL.Query()["id"]
    
    if !err || len(userID[0]) < 1 {
        panic(err)
        return
    }

    id := string(userID[0])
	fmt.Fprintln(w, "Url Param 'id' is: " + id)

	client, err2 := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	
	if err2 != nil {
		panic(err2)
	}

	defer func() {
		if err2 = client.Disconnect(context.TODO()); err2 != nil {
			panic(err2)
		}
	}()
	// Ping the primary
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected and pinged.")

	coll := client.Database("Instagram").Collection("posts")
	filter := bson.M{"postowner": id}
	query, err3 := coll.Find(context.TODO(), filter)
	if err3 != nil {
		fmt.Println("errror retrieving user userid : " + id)
	}

	var findResults []bson.M
	if err3 = query.All(context.TODO(), &findResults); err3 != nil {
		panic(err3)
	}
	for _, result := range findResults {
		fmt.Fprintln(w, result)
	}
}

func get(w http.ResponseWriter, r *http.Request) {
	userID, err := r.URL.Query()["id"]
    
    if !err || len(userID[0]) < 1 {
        panic(err)
        return
    }

    id := string(userID[0])
	fmt.Fprintln(w, "Url Param 'id' is: " + id)

	client, err2 := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	
	if err2 != nil {
		panic(err2)
	}

	defer func() {
		if err2 = client.Disconnect(context.TODO()); err2 != nil {
			panic(err2)
		}
	}()
	// Ping the primary
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected and pinged.")

	objectIDS, _ := primitive.ObjectIDFromHex(id)
	coll := client.Database("Instagram").Collection("users")
	filter := bson.M{"_id": objectIDS}
	query, err3 := coll.Find(context.TODO(), filter)
	if err3 != nil {
		fmt.Println("errror retrieving user userid : " + id)
	}

	var findResults []bson.M
	if err3 = query.All(context.TODO(), &findResults); err3 != nil {
		panic(err3)
	}
	for _, result := range findResults {
		fmt.Fprintln(w, result)
	}
}

func getUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: Get User Endpoint")
	get(w, r)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Homepage endpoint hit")
}

func handleRequests() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/users", users)
	http.HandleFunc("/users/", getUser)
	http.HandleFunc("/posts", posts)
	http.HandleFunc("/posts/", getPost)
	http.HandleFunc("/posts/users/", getAllPostsOfUser)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

const uri = "mongodb+srv://shubh:shubhgupta@intagramapicluster.rtzze.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"

func main() {
	handleRequests()
}