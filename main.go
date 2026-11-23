package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("./upload"))
	http.Handle("/files/", http.StripPrefix("/files/", fs))

	fmt.Println("Sharing files on http://localhost:8080/files/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
