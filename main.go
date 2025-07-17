package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func home(w http.ResponseWriter,r*http.Request){
	w.Write([]byte("hello from snippetbox"))
}

func snippetview(w http.ResponseWriter,r*http.Request){
	// w.Write([]byte("display the snippet.."))
	id,err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1{
		http.NotFound(w,r)
		return
	}
	msg := fmt.Sprintf("display a specific snippet with id %d..",id)
	w.Write([]byte(msg))
}

func snippetCreate(w http.ResponseWriter,r*http.Request){
	w.Write([]byte("display create a snipptex form"))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}",home)
	mux.HandleFunc("GET /snippet/view/{id}",snippetview)
	mux.HandleFunc("GET /snippet/create",snippetCreate)
	log.Print("starting server on 4000")
	

	err := http.ListenAndServe(":4000",mux)
	log.Fatal(err)
	
}