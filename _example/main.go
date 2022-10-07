package main

import (
	"github.com/Ja7ad/forker"
	"log"
	"net/http"
)

func main() {
	srv := &http.Server{
		Handler: GreetingHandler(),
	}

	f := forker.New(srv)

	log.Fatalln(f.ListenAndServe(":8080"))

}

func GreetingHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("greeting!!!"))
	}
}
