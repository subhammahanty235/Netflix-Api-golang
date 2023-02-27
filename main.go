package main

import (
	"fmt"
	"net/http"

	"github.com/subhammahanty235/netflix-api-golang/router"
)

func main() {

	fmt.Println("Hello")
	r := router.Router()
	http.ListenAndServe(":5000", r)
	fmt.Println("listening to port 5000")

}
