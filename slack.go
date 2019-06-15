package main

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"os"
	"slackProject/Response"
	"slackProject/Service"
	"strconv"
	"time"
)

var ( dep *Service.DeployService )

func ReadDeploys() {
	log.Print("hello world")
	f, err := os.Open("./src/slackProject/deploys.sqlite")
	if err != nil {
		log.Fatalf("Not able to initialize SQLite DB connectivity \n %v", err)
	}
	dep, err = Service.PopulateDeploys(*f)
}

func ResponseFormatter(resp interface{}, err error) (int, []byte) {
	if err != nil {
		message, _ := json.Marshal(Response.BaseResponse{"FAILURE", time.Now(), err.Error()})
		return http.StatusInternalServerError, message
	} else {
		b, err := json.Marshal(Response.BaseResponse{"SUCCESS", time.Now(), resp})
		if err != nil {
			message, _ := json.Marshal(Response.BaseResponse{"FAILURE", time.Now(), err.Error()})
			return http.StatusInternalServerError, message
		} else {
			return http.StatusOK, b
		}
	}
}

func getEngineers(w http.ResponseWriter, r *http.Request) {
	engineers, err := dep.GetEngineers()
	status, message := ResponseFormatter(engineers, err)
	w.WriteHeader(status)
	w.Write(message)
}

func getActionsOfEngineer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	actions, err := dep.GetActionsOfEng(vars["engineer"])
	status, message := ResponseFormatter(actions, err)
	w.WriteHeader(status)
	w.Write(message)
}


func validateInpute(fromString string, toString string) (int64, int64, error) {
	message := "I only support microsecond format"
	if len(fromString) != 10 || len(toString) != 10 {
		return -1, -1, errors.New(message)
	} else {
		from, err := strconv.Atoi(fromString)
		if err != nil {
			return -1, -1, errors.New(message)
		}
		to, err := strconv.Atoi(toString)
		if err != nil {
			return -1, -1, errors.New(message)
		}
		return int64(from), int64(to), nil
	}
}


func getEvents(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	from, to, err := validateInpute(vars["from"], vars["to"])
	if err != nil {
		status, message := ResponseFormatter(nil, err)
		w.WriteHeader(status)
		w.Write(message)
	} else {
		events, err := dep.GetEventsDuring(from, to)
		status, message := ResponseFormatter(events, err)
		w.WriteHeader(status)
		w.Write(message)
	}
}


func getEventStats(w http.ResponseWriter, r *http.Request) {
	events, err := dep.GetEventStats()
	status, message := ResponseFormatter(events, err)
	w.WriteHeader(status)
	w.Write(message)
}


func main() {
	ReadDeploys()
	r := mux.NewRouter()
	r.HandleFunc("/engineers", getEngineers).Methods("GET")
	r.HandleFunc("/actions/engineer/{engineer}", getActionsOfEngineer).Methods("GET")
	r.HandleFunc("/events/from/{from}/to/{to}", getEvents).Methods("GET")
	r.HandleFunc("/eventStats", getEventStats).Methods("GET")
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
