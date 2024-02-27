package main

import (
	"log"
	"net/http"
	"todoPchi/database"
	"todoPchi/web"
)

func main() {
	url := "mongodb://localhost:27017"
	store, err := database.NewStore(url)
	if err != nil {
		log.Fatal(err)
	}

	h := web.NewHadler(store)
	err = http.ListenAndServe(":3000", h)
	if err != nil {
		log.Fatal(err)
	}

}
