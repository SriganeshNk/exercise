package Api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"slackProject/Response"
	"slackProject/Service"
	"time"
)

var ( dep *Service.DeployService )

func ReadDeploys() {
	/**
	Initializes the deployService to be used by handlers of the API layer
	 */
	log.Print("hello world")
	f, err := os.Open("../deploys.sqlite")
	if err != nil {
		log.Fatalf("Not able to initialize SQLite DB connectivity \n %v", err)
	}
	dep, err = Service.PopulateDeploys(*f)
	if err != nil {
		log.Fatalf("Not able to initialize the deploy tables, %v", err)
	}
}

func responseFormatter(resp interface{}, err error) (int, []byte) {
	/**
	Response formatter formats the response with three values, the status, timestamp and result
	errors are reported separately with write Headers in the handler function
	 */
	if err != nil {
		message, _ := json.Marshal(Response.BaseResponse{Status:"FAILURE", Timestamp:time.Now().Unix(), Result:err.Error()})
		return http.StatusInternalServerError, message
	} else {
		b, err := json.Marshal(Response.BaseResponse{Status:"SUCCESS", Timestamp:time.Now().Unix(), Result: resp})
		if err != nil {
			message, _ := json.Marshal(Response.BaseResponse{Status:"FAILURE", Timestamp:time.Now().Unix(), Result:err.Error()})
			return http.StatusInternalServerError, message
		} else {
			return http.StatusOK, b
		}
	}
}

func GetEngineers(w http.ResponseWriter, r *http.Request) {
	/**
	Handler function to return the distinct engineers from the database
	 */
	engineers, err := dep.GetEngineers()
	status, message := responseFormatter(engineers, err)
	w.WriteHeader(status)
	n, err := w.Write(message)
	if err != nil {
		log.Printf("Not able to send the response back %v", err)
	}
	go log.Printf("response sent %d", n)
}

func GetActionsOfEngineer(w http.ResponseWriter, r *http.Request) {
	/**
	Handler function to return the actions of an engineer from the database
	*/
	vars := mux.Vars(r)
	log.Print(vars)
	actions, err := dep.GetActionsOfEng(vars["engineer"])
	status, message := responseFormatter(actions, err)
	w.WriteHeader(status)
	n, err := w.Write(message)
	if err != nil {
		log.Printf("Not able to send the response back %v", err)
	}
	go log.Printf("response sent %d", n)
}


func GetEvents(w http.ResponseWriter, r *http.Request) {
	/**
	Handler function to return the list of events within a time range
	Checks for the validity of the time range and also only accepts unix 10 digit timestamps as valid time input
	 */
	vars := mux.Vars(r)
	events, err := dep.GetEventsDuring(vars["from"], vars["to"])
	status, message := responseFormatter(events, err)
	w.WriteHeader(status)
	n, err := w.Write(message)
	if err != nil {
		log.Printf("Not able to send the response back %v", err)
	}
	go log.Printf("response sent %d", n)
}


func GetEventStats(w http.ResponseWriter, r *http.Request) {
	/**
	Handles two types of requests, either get event stats from all the dates in the DB
	else returns the stats for a particular day
	 */
	vars := mux.Vars(r)
	events, err := dep.GetEventStats(vars["date"])
	status, message := responseFormatter(events, err)
	w.WriteHeader(status)
	n, err := w.Write(message)
	if err != nil {
		log.Printf("Not able to send the response back %v", err)
	}
	go log.Printf("response sent %d", n)
}