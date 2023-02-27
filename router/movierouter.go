package router

import (
	"github.com/gorilla/mux"
	"github.com/subhammahanty235/netflix-api-golang/controller"
)

func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/movies", controller.GetAllMovies).Methods("GET")
	router.HandleFunc("/api/movie/{id}", controller.GetOneMovie).Methods("GET")
	router.HandleFunc("/api/movie", controller.CreateMovie).Methods("POST")
	router.HandleFunc("/api/updatemovie/{id}", controller.MarkAsWatched).Methods("PUT")
	router.HandleFunc("/api/deletemovie/{id}", controller.DeleteMovie).Methods("DELETE")

	return router

}
