/*
Go. Homework 2 task 2
Zaur Malakhov, dated Sep 18, 2019
Данные посылал через Advanced REST client

*/

package main

import (
	"net/http"
	"log"
	"encoding/json"
)


func main(){
	router := http.NewServeMux()
	router.HandleFunc( "/get" , getHandler)
	router.HandleFunc( "/read" , readHandler)
	
	log.Fatal(http.ListenAndServe( ":8080" , router))
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET"{
		cookie    :=    http.Cookie{Name: "name", Value: "John"}
		http.SetCookie(w, &cookie)
	}
}

func readHandler(w http.ResponseWriter, r *http.Request) {
	cookie, _ :=r.Cookie( "name" )
	foundCookie, err := json.Marshal(cookie)
	if err != nil {
		log.Panic(err)
	}
	w.Write([] byte (foundCookie))
}
