package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	server "github.com/m-a-r-a-t/go-jwt-auth/internal"
	_ "github.com/m-a-r-a-t/go-jwt-auth/internal/controllers/auth"
)

func main() {
	var host string
	for _, v := range os.Args {
		if v == "--host" {
			fmt.Println("Введите host: ")
			fmt.Fscan(os.Stdin, &host)
		}
	}
	// в make file добавить добавление аргументов в командной строке
	if &host == nil {
		host = "localhost:8000"
	}
	log.Println("Server starting on port: 8000")
	err := http.ListenAndServe(host, &server.R.Mux)
	log.Fatal(err)
}
