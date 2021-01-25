package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"nononsensecode/rest-api-tutorial/articleerror"
	model "nononsensecode/rest-api-tutorial/model"
	utils "nononsensecode/rest-api-tutorial/utils"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

func init() {
	var err error
	model.DB, err = sql.Open("sqlite3", "./article.db")
	if err != nil {
		log.Fatal("DB cannot be opened!", err)
	}

	err = model.CreateTable()
	if err != nil {
		log.Fatal("Table cannot be created", err)
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homepage")
}

func returnAllArticles(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllArticles")
	articles, err := model.FindAllArticles()
	if err != nil {
		log.Panic(err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Unknown error")
	}
	utils.RespondWithJSON(w, http.StatusOK, articles)
}

func returnArticleByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnArticleByID")
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid ID!")
	}
	
	article, err := model.FindArticleByID(id)
	if err != nil {
		articleError, ok := err.(*articleerror.ArticleError)
		if ok {
			switch{
			case articleError.IsArticleEmpty():
				utils.RespondWithError(w, http.StatusNotFound, articleError.Error())
			default:
				utils.RespondWithError(w, http.StatusInternalServerError, "Unknown error occurred")
			}
			log.Println(err)
			return
		}
	}
	utils.RespondWithJSON(w, http.StatusOK, article)
}

func createArticle(w http.ResponseWriter, r *http.Request) {
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Request is not readable")
		log.Println(err)
	}
	
	var article model.Article
	err = json.Unmarshal(requestBody, &article)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Request cannot be converted")
		log.Println(err)
	}

	articleID, err := model.CreateArticle(article)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Article cannot be created")
	}
	article.ID = articleID

	utils.RespondWithJSON(w, http.StatusCreated, article)
}

func deleteArticleByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hit endpoint deleteArticleByID")
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {		
		utils.RespondWithError(w, http.StatusBadRequest, "ID is invalid")
	}
	
	err = model.DeleteArticleByID(id)
	if err != nil {
		message := fmt.Sprintf("Article with ID %d cannot be deleted", id)
		utils.RespondWithError(w, http.StatusInternalServerError, message)
	}

	utils.RespondWithJSON(w, http.StatusNoContent, model.Article{})
}

func updateArticle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hit endpoint updateArticle")
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Request is not readable")
		log.Println(err)
	}

	var updatedArticle model.Article
	err = json.Unmarshal(requestBody, &updatedArticle)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Request cannot be converted")
		log.Println(err)
	}

	article, err := model.UpdateArticle(updatedArticle)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Article cannot be updated due to unknown error")
		log.Panic(err)
	}
	utils.RespondWithJSON(w, http.StatusOK, article)

}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/articles", returnAllArticles)
	myRouter.HandleFunc("/article", createArticle).Methods("POST")
	myRouter.HandleFunc("/article/{id}", returnArticleByID).Methods("GET")
	myRouter.HandleFunc("/article/{id}", deleteArticleByID).Methods("DELETE")
	myRouter.HandleFunc("/article", updateArticle).Methods("PUT")
	//In the below code, as you can see method is not specified. So if we put this
	//code above the delete method, delete will also give you get method result 
	//instead of delete method
	//myRouter.HandleFunc("/article/{id}", returnArticleByID)
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func main() {	
	handleRequests()
}
