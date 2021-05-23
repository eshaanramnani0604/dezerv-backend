package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"net/mail"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"dezerv-test/helper"
	"dezerv-test/models"
)

var client *mongo.Client




//register function

func Register(response http.ResponseWriter, request *http.Request) {
	
	//setting cors headers
	
	response.Header().Set("Access-Control-Allow-Origin", "*")
	
	response.Header().Set("Content-type", "application/json")
	response.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	response.Header().Set("Access-Control-Allow-Credentials", "true")
	response.Header().Set("Access-Control-Max-Age", "86400")
	response.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Origin,Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	
	var user models.User
	var res models.User
	var custErr models.CustomError
	custErr.IsPasswordError=false
	custErr.Success=false
	custErr.Errors=[]string{}
	var success models.Success
	success.Success=true
	json.NewDecoder(request.Body).Decode(&user)
	//fmt.Println("after decode body")
	collection := client.Database("dezerv-test").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	//check if username already exists
	noErr := collection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&res)
	//fmt.Println(noErr)
	if noErr != nil {
		//check if valid email
		_, err := mail.ParseAddress(user.Email)
		if(err!=nil){
			custErr.Message="not a valid email"
			
			response.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(response).Encode(custErr)
			return 
		}
		if user.Email==user.Password{
			custErr.Message="error in Password"
			custErr.IsPasswordError=true
			custErr.Errors=append(custErr.Errors,"password cannot be same as username")
			response.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(response).Encode(custErr)
			return
		}
		//check if password meets criteria
		isValid, errors:=helper.CheckPassword(user.Password)
		
		if !isValid{
			custErr.Message="error in Password"
			custErr.IsPasswordError=true
			custErr.Errors= errors
			response.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(response).Encode(custErr)
			return
		}
		//save encrypted password
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
		user.Password =string(hashedPassword)
		collection.InsertOne(ctx, user)
		success.Message="user created successfully"
		json.NewEncoder(response).Encode(success)
		return
	}else{
		custErr.Message="user already exists"
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(custErr)
		return
	}
	
}
//sign in function
func SignIn(response http.ResponseWriter, request *http.Request) {
	//setting cors headers
	response.Header().Set("Access-Control-Allow-Origin", "*")
	//fmt.Println("here")
	response.Header().Set("Content-type", "application/json")
	response.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	response.Header().Set("Access-Control-Allow-Credentials", "true")
	response.Header().Set("Access-Control-Max-Age", "86400")
	response.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Origin,Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	var user models.User
	var res models.User
	var custErr models.CustomError
	var success models.Success
	collection := client.Database("dezerv-test").Collection("users")
	
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	json.NewDecoder(request.Body).Decode(&user)
	
	err := collection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&res)
	if err != nil {
		custErr.Message="user not found"
		
		//response.WriteHeader(http.StatusBadRequest) commented because causing CORS error on frontend
		json.NewEncoder(response).Encode(custErr)
		return
	}
	//check if password matches the saved encrypted one
	pwdNoMatch :=bcrypt.CompareHashAndPassword([]byte(res.Password),[]byte(user.Password) )
	if pwdNoMatch!=nil{
		custErr.Message="incorrect password"
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(custErr)
		return
	}
	success.Message="password matched! login Successful"
	success.Success=true
	json.NewEncoder(response).Encode(success)
	
	
 }
 //cahnge password function
func ChangePassword(response http.ResponseWriter, request *http.Request) { 
	//set CORS headers
	response.Header().Set("Access-Control-Allow-Origin", "*")
	//fmt.Println("here")
	response.Header().Set("Content-type", "application/json")
	response.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	response.Header().Set("Access-Control-Allow-Credentials", "true")
	response.Header().Set("Access-Control-Max-Age", "86400")
	response.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Origin,Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	var user models.User
	var res models.User
	var custErr models.CustomError
	var success models.Success
	json.NewDecoder(request.Body).Decode(&user)
	
	collection := client.Database("dezerv-test").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err := collection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&res)
	if err != nil {
		custErr.Message="user not found"
		//response.WriteHeader(http.StatusBadRequest) commented because causing CORS error on frontend
		json.NewEncoder(response).Encode(custErr)
		return
	}
	if user.Email==user.Password{
		custErr.Message="error in Password"
		custErr.IsPasswordError=true
		custErr.Errors=append(custErr.Errors,"password cannot be same as username")
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(custErr)
		return
	}
	
	pwdNoMatch :=bcrypt.CompareHashAndPassword([]byte(res.Password),[]byte(user.Password) )
	if pwdNoMatch==nil{
		custErr.Message="error in Password"
		custErr.IsPasswordError=true
		custErr.Errors=append(custErr.Errors,"password cannot be same as old password")
		response.WriteHeader(http.StatusBadRequest)
		
		json.NewEncoder(response).Encode(custErr)
		return
	}
	isValid, errors:=helper.CheckPassword(user.Password)
	
	if !isValid{
		custErr.Message="error in Password"
		custErr.IsPasswordError=true
		custErr.Errors= errors
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(custErr)
		return
	}
	//update user with new encrypted password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
	user.Password =string(hashedPassword)
	update:=bson.D{{
		"$set",bson.D{
			{"password",user.Password},
		},
	}}
	collection.UpdateOne(ctx,bson.M{"email": user.Email}, update)
	success.Success=true
	success.Message="user password updated successfully"
	json.NewEncoder(response).Encode(success)
	return


}

func main() {
	fmt.Println("Starting the application...")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(ctx, clientOptions)
	router := mux.NewRouter()
	router.HandleFunc("/register", Register).Methods("POST", "OPTIONS")
	router.HandleFunc("/sign-in", SignIn).Methods("POST", "OPTIONS")
	router.HandleFunc("/change-password", ChangePassword).Methods("POST", "OPTIONS")
	http.ListenAndServe(":8000", router)
}