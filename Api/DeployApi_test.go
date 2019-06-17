package Api

import (
	"bytes"
	"context"
	"encoding/json"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"net/http/httptest"
	"slackProject/Response"
	"strings"
	"testing"
)

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
	handler := http.HandlerFunc(GetEngineers)

	handler.ServeHTTP(rr, req)

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
	req, err := http.NewRequest("GET", "/actions/engineer/{engineer}", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(GetActionsOfEngineer)

	// Populate the request's context with our test data.
	ctx := req.Context()
	ctx = context.WithValue(ctx, "engineer", "vincent")

	req = req.WithContext(ctx)

	handler.ServeHTTP(rr, req)

	if status:= rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	expected := 111
	result, n := getResponse(rr, t)
	t.Log(result.String())
	x := strings.Split(string(result.Bytes()[:n]), ",")
	if expected != len(x) {
		t.Errorf("Expected %d, got %d", expected, len(x))
	}
}