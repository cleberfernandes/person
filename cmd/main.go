package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/bson"
	person "person.com"
)

var mongoHandler *person.MongoHandler

func registerRoutes() http.Handler {
	r := chi.NewRouter()
	r.Route("/persons", func(r chi.Router) {
		r.Get("/", getAllPerson)                 //GET /persons
		r.Get("/{phonenumber}", getPerson)       //GET /persons/991158888
		r.Post("/", addPerson)                   //POST /persons
		r.Put("/{phonenumber}", updatePerson)    //PUT /persons/991158888
		r.Delete("/{phonenumber}", deletePerson) //DELETE /persons/991158888
	})
	return r
}

func main() {
	mongoDbConnection := "mongodb://localhost:27017"
	mongoHandler := person.NewHandler(mongoDbConnection)
	fmt.Printf("mongoHandler: %v\n", mongoHandler)

	handler := registerRoutes()
	log.Fatal(http.ListenAndServe(":3060", handler))
}

func getPerson(responseWriter http.ResponseWriter, request *http.Request) {
	phoneNumber := chi.URLParam(request, "phonenumber")
	if phoneNumber == "" {
		http.Error(responseWriter, http.StatusText(404), 404)
		return
	}

	person := &person.Person{}
	err := mongoHandler.GetOne(person, bson.M{"phoneNumber": phoneNumber})
	if err != nil {
		http.Error(responseWriter, fmt.Sprintf("Person with phonenumber %S not found!", phoneNumber), 404)
		return
	}

	json.NewEncoder(responseWriter).Encode(person)
}

func getAllPerson(responseWriter http.ResponseWriter, request *http.Request) {
	persons := mongoHandler.Get(bson.M{})
	json.NewEncoder(responseWriter).Encode(persons)
}

func addPerson(responseWriter http.ResponseWriter, request *http.Request) {
	existingPerson := &person.Person{}
	var person person.Person
	json.NewDecoder(request.Body).Decode(&person)
	person.CreatedOn = time.Now()
	err := mongoHandler.GetOne(existingPerson, bson.M{"phoneNumber": person.PhoneNumber})
	if err == nil {
		http.Error(responseWriter, fmt.Sprintf("Person with phonenumber %S already exist!", person.PhoneNumber), 404)
		return
	}
	_, err1 := mongoHandler.AddOne(&person)
	if err != nil {
		http.Error(responseWriter, fmt.Sprint(err1), 400)
	}
	responseWriter.Write([]byte("Person created successfully."))
	responseWriter.WriteHeader(201)
}

func deletePerson(responseWriter http.ResponseWriter, request *http.Request) {
	existingPerson := &person.Person{}
	phoneNumber := chi.URLParam(request, "phonenumber")
	if phoneNumber == "" {
		http.Error(responseWriter, http.StatusText(404), 404)
		return
	}
	err := mongoHandler.GetOne(existingPerson, bson.M{"phoneNumber": phoneNumber})
	if err != nil {
		http.Error(responseWriter, fmt.Sprintf("Person with phonenumber %S does not exist!", phoneNumber), 400)
		return
	}
	_, err = mongoHandler.RemoveOne(bson.M{"phoneNumber": phoneNumber})
	if err != nil {
		http.Error(responseWriter, fmt.Sprint(err), 400)
		return
	}
	responseWriter.Write([]byte("Person deleted."))
	responseWriter.WriteHeader(200)
}

func updatePerson(responseWriter http.ResponseWriter, request *http.Request) {
	phoneNumber := chi.URLParam(request, "phoneNumber")
	if phoneNumber == "" {
		http.Error(responseWriter, http.StatusText(404), 404)
		return
	}
	person := &person.Person{}
	json.NewDecoder(request.Body).Decode(person)
	_, err := mongoHandler.Update(bson.M{"phoneNumber": phoneNumber}, person)
	if err != nil {
		http.Error(responseWriter, fmt.Sprint(err), 400)
		return
	}
	responseWriter.Write([]byte("Person updated."))
	responseWriter.WriteHeader(200)
}
