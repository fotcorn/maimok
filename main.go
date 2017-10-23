package main

import (
	"net/http"

	"github.com/maimok/backend/maimok"
)

func main() {
	http.ListenAndServe(":3000", maimok.GetRouter())
}
