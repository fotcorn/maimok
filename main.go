package main

import (
	"fmt"
	"net/http"

	"github.com/maimok/backend/maimok"
)

func main() {
	fmt.Println("Starting server on port 3000")
	http.ListenAndServe(":3000", maimok.GetRouter())
}
