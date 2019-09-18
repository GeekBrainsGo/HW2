/*
Go. Homework 2 task 1
Zaur Malakhov, dated Sep 18, 2019
Данные посылал через Advanced REST client

{
	"search": "автор"
	"sites": [
      "https://book24.ru",
	  "https://avidreaders.ru/books/",
	  "https://www.iherb.com"
  ]
}


*/


package main

import (
    "fmt"
    "encoding/json"
    "io/ioutil"
	"net/http"
	"log"
	"strings"
	"os"
)

//структура описывающая запрос
type Find struct {
    Search    string `json:"search"`
    Sites     []string  `json:"sites"`
}

func main(){
	router := http.NewServeMux()
	router.HandleFunc( "/" , requestHandler)
	
	log.Fatal(http.ListenAndServe( ":8080" , router))
}

func requestHandler(w http.ResponseWriter, r *http.Request) {

    var resp Find

    //читаем тело запроса
    body, err := ioutil.ReadAll(r.Body)
    //проверяем на наличие ошибки
    if err != nil {
        fmt.Fprintf(w, "err %q\n",  err.Error())
    } else {
        //если все нормально - пишем по указателю в структуру
        err = json.Unmarshal(body, &resp)
        if err != nil {
            fmt.Println(w, "can't unmarshal: ", err.Error())
        }
	}
	

	fmt.Println(resp.Search)
	fmt.Println(resp.Sites)

	query := resp.Search
	linksInput := resp.Sites

	

	//w.Write([] byte ( body ))

	//links := []string{"https://book24.ru", "https://avidreaders.ru/books/", "https://www.iherb.com"}
	links := findQuery(query, linksInput)

	foundLinks, err := json.Marshal(links)
	if err != nil {
		log.Panic(err)
	}

	w.Write([] byte (foundLinks))
	

    
}



func findQuery(q string, m1 []string) (m2 []string){
	for x := range m1{
		//fmt.Println(m1[x])
		if checkContainsQuery(q,m1[x]){
			m2 = append(m2, m1[x])
		}
	}
	return
}

func checkContainsQuery(query, url string) bool {
	res := false

	body := getBodyFromUrl(url)
	if strings.Contains(body, query){
		res = true
	}

	return res
}

func getBodyFromUrl(url string) string{

	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: %v\n",err)
		os.Exit(1)
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: чтение %s: %v\n", url, err)
		os.Exit(1)
	}
	//fmt.Printf("%s", body)
	return string(body)

}




