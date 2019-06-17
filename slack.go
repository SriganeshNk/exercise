package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"os"
	"slackProject/Api"
)


func main() {
	Api.ReadDeploys()
	r := mux.NewRouter()
	r.HandleFunc("/engineers", Api.GetEngineers).Methods("GET")
	r.HandleFunc("/actions/engineer/{engineer}", Api.GetActionsOfEngineer).Methods("GET")
	r.HandleFunc("/events/from/{from}/to/{to}", Api.GetEvents).Methods("GET")
	r.HandleFunc("/eventStats", Api.GetEventStats).Methods("GET")
	r.HandleFunc("/eventStats/{date}", Api.GetEventStats).Methods("GET")
	http.Handle("/", r)
	logHandler := handlers.LoggingHandler(os.Stdout, r)
	log.Fatal(http.ListenAndServe(":8080", logHandler))
}
