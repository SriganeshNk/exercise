package Api

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"net/http/httptest"
	"slackProject/Response"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

var once = sync.Once{}
var (
	router *mux.Router
	)

func getRouter() (*mux.Router) {
	once.Do(
		func(){
			router = mux.NewRouter()
			router.HandleFunc("/engineers", GetEngineers).Methods("GET")
			router.HandleFunc("/actions/engineer/{engineer}", GetActionsOfEngineer).Methods("GET")
			router.HandleFunc("/events/from/{from}/to/{to}", GetEvents).Methods("GET")
			router.HandleFunc("/eventStats", GetEventStats).Methods("GET")
			router.HandleFunc("/eventStats/{date}", GetEventStats).Methods("GET")
		})
	return router
}

func getResponse(rr *httptest.ResponseRecorder, t *testing.T) (bytes.Buffer, int){
	resp := Response.BaseResponse{}
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	byteArray, err := json.Marshal(resp.Result)
	if err != nil {
		t.Fatal(err)
	}
	result := bytes.Buffer{}
	n, err := result.Write(byteArray)
	if err != nil {
		t.Fatal(err)
	}
	return result, n
}

func TestGetEngineersHandler(t *testing.T) {
	ReadDeploys()
	req, err := http.NewRequest("GET", "/engineers", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	getRouter().ServeHTTP(rr, req)

	if status:= rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	expected := 111
	result, n := getResponse(rr, t)
	x := strings.Split(string(result.Bytes()[:n]), ",")
	if expected != len(x) {
		t.Errorf("Expected %d, got %d", expected, len(x))
	}
}


func TestGetActionsOfEngineer(t *testing.T) {
	ReadDeploys()
	req, err := http.NewRequest("GET", "/actions/engineer/vincent", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	getRouter().ServeHTTP(rr, req)

	if status:= rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	expected := 36
	result, n := getResponse(rr, t)
	t.Log(result.String())
	x := strings.Split(string(result.Bytes()[:n]), ",")
	if expected != len(x) {
		t.Errorf("Expected %d, got %d", expected, len(x))
	}
}


func TestGetEvents(t *testing.T) {
	ReadDeploys()
	start := strconv.Itoa(int(time.Date(2017, time.October, 28, 0, 0, 0, 0, time.UTC).Unix()))
	end := strconv.Itoa(int(time.Date(2017, time.October, 30, 0, 0, 0, 0, time.UTC).Unix()))
	req, err := http.NewRequest("GET", "/events/from/"+start+"/to/"+end, nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	getRouter().ServeHTTP(rr, req)

	if status:= rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	expected := 147
	result, n := getResponse(rr, t)
	t.Log(result.String())
	x := strings.Split(string(result.Bytes()[:n]), ",")
	if expected != len(x) {
		t.Errorf("Expected %d, got %d", expected, len(x))
	}
}



func TestGetEventStats(t *testing.T) {
	ReadDeploys()
	start := strconv.Itoa(int(time.Date(2017, time.October, 28, 0, 0, 0, 0, time.UTC).Unix()))
	req, err := http.NewRequest("GET", "/eventStats/"+start, nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	getRouter().ServeHTTP(rr, req)

	if status:= rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	expected := 2
	result, n := getResponse(rr, t)
	t.Log(result.String())
	x := strings.Split(string(result.Bytes()[:n]), ",")
	if expected != len(x) {
		t.Errorf("Expected %d, got %d", expected, len(x))
	}
}


func TestGetEventStatsAll(t *testing.T) {
	ReadDeploys()
	req, err := http.NewRequest("GET", "/eventStats", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	getRouter().ServeHTTP(rr, req)

	if status:= rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	expected := 16
	result, n := getResponse(rr, t)
	t.Log(result.String())
	x := strings.Split(string(result.Bytes()[:n]), ",")
	if expected != len(x) {
		t.Errorf("Expected %d, got %d", expected, len(x))
	}
}