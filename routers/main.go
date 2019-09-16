/*
 * HomeWork-2: Routers and simple cookie control
 * Created on 15.09.2019 22:02
 * Copyright (c) 2019 - Eugene Klimov
 */

package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	servAddr    = "localhost:8080"
	cookieName  = "session_id"
	cookieValue = "Klim_GeekBrains_Go"
)

func mainPage(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session_id")
	loggedIn := err != http.ErrNoCookie

	if loggedIn {
		io.WriteString(w, `<a href="/logout">Logout,</a>&nbsp`+session.Value)
	} else {
		io.WriteString(w, `<a href="/login">Login</a>, please!`)
	}
}

func loginPage(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:    cookieName,
		Value:   cookieValue,
		Expires: time.Now().AddDate(0, 0, 1),
	})
	http.Redirect(w, r, "/", http.StatusFound)
}

func logoutPage(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   cookieName,
		MaxAge: -1,
	})
	http.Redirect(w, r, "/", http.StatusFound)
}

func main() {

	// safe shutdown
	shutdown := make(chan os.Signal)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM) // os.Kill wrong for linters

	go func() {
		<-shutdown
		// any work here
		fmt.Printf("\nShutdown server at: %s\n", servAddr)
		os.Exit(0)
	}()

	// prepare server, no need smart router for simple scenario
	http.HandleFunc("/", mainPage)
	http.HandleFunc("/login", loginPage)
	http.HandleFunc("/logout", logoutPage)

	fmt.Println("Starting server at:", servAddr)
	log.Fatalln(http.ListenAndServe(servAddr, nil))
}
