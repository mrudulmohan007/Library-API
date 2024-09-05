package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// Model for course
type Book struct {
	BookId    string  `json:"bookid"`
	BookName  string  `json:"bookname"`
	BookPrice int     `json:"price"`
	Author    *Author `json:"author"`
}
type Author struct {
	Fullname     string `json:"fullname"`
	Author_Place string `json:"authorplace"`
}

// fake db
var books []Book

// middleware , helper-file
func (b *Book) IsEmpty() bool {
	//return b.BookId == "" && b.BookName == ""
	return b.BookName == ""
} //return true if bookid and bookname is empty

func main() {
	fmt.Println("API IN GOLANG")
	r := mux.NewRouter()

	//seeding of the data
	books = append(books, Book{BookId: "2", BookName: "Alchemist", BookPrice: 299, Author: &Author{Fullname: "Paulo Coelho", Author_Place: "Brazil"}})
	books = append(books, Book{BookId: "3", BookName: "Two states", BookPrice: 499, Author: &Author{Fullname: "Chetan Bhagat", Author_Place: "India"}})

	//routing
	r.HandleFunc("/", serveHome).Methods("GET")
	r.HandleFunc("/books", getAllBooks).Methods("GET")
	r.HandleFunc("/book/{id}", getOneBook).Methods("GET")
	r.HandleFunc("/book", createOneBook).Methods("POST")
	r.HandleFunc("/book/{id}", updateOneBook).Methods("PUT")
	r.HandleFunc("/book/{id}", deleteOneBook).Methods("DELETE")

	//listen to a port
	log.Fatal(http.ListenAndServe(":8080", r))
}

//controllers

// serve home route
func serveHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Hello we are trying to build an API</h1>"))
}
func getAllBooks(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Getting all Books")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func getOneBook(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get one book ")
	w.Header().Set("Content-Type", "application/json")
	//here we provide the id of the book to be searched so we should grab id from the request

	params := mux.Vars(r)

	//loop through books,find matching id and return the response
	for _, book := range books {
		if book.BookId == params["id"] {
			json.NewEncoder(w).Encode(book)
			return
		}
	}
	json.NewEncoder(w).Encode("No Book found with given id")
	return
}

func createOneBook(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create one book")
	w.Header().Set("Content-Type", "application/json")

	//what if the body is empty
	if r.Body == nil {
		json.NewEncoder(w).Encode("Please send some data")
	}
	//what about the data send like this: {}
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	if book.IsEmpty() {
		json.NewEncoder(w).Encode("NO data inside the json you are senting !")
		return
	}
	// Loop through courses to check for duplicate CourseName
	for _, existingBook := range books {
		if existingBook.BookName == book.BookName {
			json.NewEncoder(w).Encode("Book name already exists")
			return
		}
	}
	//generate unique id,string
	//append book into books
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)
	book.BookId = strconv.Itoa(random.Intn(100))
	books = append(books, book)
	json.NewEncoder(w).Encode(book)
}

func updateOneBook(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Updating one book")
	w.Header().Set("Content-Type", "application/json")

	//first grab id from the req

	params := mux.Vars(r)
	//loop through the value, once we get the id,remove the book and add with the id we are passing

	for index, book := range books {
		if book.BookId == params["id"] {
			books = append(books[:index], books[index+1:]...)
			var book Book
			_ = json.NewDecoder(r.Body).Decode(&book)
			book.BookId = params["id"]
			books = append(books, book)
			json.NewEncoder(w).Encode(book)
			return
		}

	}

}

func deleteOneBook(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Deleting one course")
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	//loop,find id,remove
	for index, book := range books {
		if book.BookId == params["id"] {
			books = append(books[:index], books[index+1:]...)
			break
		}
	}

}
