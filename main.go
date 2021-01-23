package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	model "nononsensecode/rest-api-tutorial/model"

	"github.com/gorilla/mux"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homepage")
}

func returnAllArticles(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllArticles")
	json.NewEncoder(w).Encode(model.Articles)
}

func returnArticleByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnArticleByID")
	vars := mux.Vars(r)
	id := vars["id"]

	for _, article := range model.Articles {
		if article.ID == id {
			json.NewEncoder(w).Encode(article)
		}
	}
}

func createArticle(w http.ResponseWriter, r *http.Request) {
	requestBody, _ := ioutil.ReadAll(r.Body)
	var article model.Article
	json.Unmarshal(requestBody, &article)
	model.Articles = append(model.Articles, article)
	json.NewEncoder(w).Encode(article)
}

func deleteArticleByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hit endpoint deleteArticleByID")
	vars := mux.Vars(r)
	id := vars["id"]
	for index, article := range model.Articles {
		if article.ID == id {
			model.Articles = append(model.Articles[:index], model.Articles[index+1:]...)
		}
	}
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/articles", returnAllArticles)
	myRouter.HandleFunc("/article", createArticle).Methods("POST")
	myRouter.HandleFunc("/article/{id}", returnArticleByID).Methods("GET")
	myRouter.HandleFunc("/article/{id}", deleteArticleByID).Methods("DELETE")
	//myRouter.HandleFunc("/article/{id}", returnArticleByID)
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func main() {
	model.Articles = []model.Article{
		{ID: "1", Title: "Hello", Desc: "Article Description", Content: "Article Content"},
		{ID: "2", Title: "Hello 2", Desc: "Article Description", Content: "Article Content"},
	}
	handleRequests()
}
