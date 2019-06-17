package Service

import (
	_ "github.com/mattn/go-sqlite3"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"
)

var once = sync.Once{}
var (
	dep *DeployService
)

func getDeployService(t * testing.T) (*DeployService){
	once.Do(func() {
		f, err := os.Open("../deploys.sqlite")
		if err != nil {
			t.Errorf("There should be a file existing instead got \n, %v", err)
		}
		d, err := PopulateDeploys(*f)
		dep = d
		if err != nil {
			t.Errorf("We should be able to populate the deploy service succesfully instead got \n, %v", err)
		}
		if len(dep.deploys) != 1000 {
			t.Errorf("There should be 100 records instead got %d", len(dep.deploys))
		}
	})
	return dep
}


func TestPopulateDeploys(t *testing.T) {
	getDeployService(t)
}


func TestDeployService_GetActionsOfEng(t *testing.T) {
	dep := getDeployService(t)
	ans1, err := dep.GetActionsOfEngineer("vincent")
	if err != nil {
		t.Errorf("Should be able to get the actions using deploy cache but instead got %v", err)
	}
	ans2, err := dep.GetActionsOfEng("vincent")
	if err != nil {
		t.Errorf("Should be able to get the actions using sql statements but instead got %v", err)
	}
	if len(ans1) != len(ans2) {
		t.Errorf("Should be the same result from both the cache and the DB but instead got %d, %d",
			len(ans1), len(ans2))
	}
}


func TestDeployService_GetDistinctEngineers(t *testing.T) {
	dep := getDeployService(t)
	ans1, err := dep.GetEngineers()
	if err != nil {
		t.Errorf("Should be able to get the engineers using deploy cache but instead got %v", err)
	}
	ans2, err := dep.GetDistinctEngineers()
	if err != nil {
		t.Errorf("Should be able to get the engineers using sql statements but instead got %v", err)
	}
	if len(ans1) != len(ans2) {
		t.Errorf("Should be the same result from both the cache and the DB but instead got %d, %d",
			len(ans1), len(ans2))
	}
}


func GetEventsHelper(start string, end string, t *testing.T) (*DeployService, error) {
	t.Log(start, end)
	dep := getDeployService(t)
	ans1, err := dep.GetEvents(start, end)
	if err != nil {
		return nil, err
	}
	ans2, err := dep.GetDistinctEngineers()
	if err != nil {
		return nil, err
	}
	if len(ans1) != len(ans2) {
		return nil, err
	}
	return dep, nil
}

func TestDeployService_GetEvents_Success(t *testing.T) {
	start := time.Date(2017, time.October, 28, 0, 0, 0, 0, time.UTC)
	end := time.Date(2017, time.October, 30, 0, 0, 0, 0, time.UTC)
	_, err := GetEventsHelper(
		strconv.Itoa(int(start.Unix())),
		strconv.Itoa(int(end.Unix())),
		t)
	if err != nil {
		t.Errorf("Should be able to get the events instead found %v", err)
	}
}

func TestDeployService_GetEvents_InvalidRange(t *testing.T) {
	start := time.Date(2017, time.November, 28, 0, 0, 0, 0, time.UTC)
	end := time.Date(2017, time.October, 30, 0, 0, 0, 0, time.UTC)
	_, err := GetEventsHelper(
		strconv.Itoa(int(start.Unix())),
		strconv.Itoa(int(end.Unix())),
		t)
	t.Logf("Error should not be null, got %v", err)
	if err == nil {
		t.Errorf("Should throw error complaining about time got %v", err)
	}
}


func TestDeployService_GetEvents_InvalidFormat(t *testing.T) {
	start := time.Date(2017, time.November, 28, 0, 0, 0, 0, time.UTC)
	end := time.Date(2017, time.October, 30, 0, 0, 0, 0, time.UTC)
	_, err := GetEventsHelper(start.String(), end.String(), t)
	t.Logf("Error should not be null, got %v", err)
	if err == nil {
		t.Errorf("Should throw error complaining about time got %v", err)
	}
}

func TestDeployService_GetEvents(t *testing.T) {
	dep := getDeployService(t)
	start := time.Date(2017, time.October, 28, 0, 0, 0, 0, time.UTC)
	end := time.Date(2017, time.October, 30, 0, 0, 0, 0, time.UTC)
	ans1, err := dep.GetEvents(
		strconv.Itoa(int(start.Unix())),
		strconv.Itoa(int(end.Unix())))
	if err != nil {
		t.Errorf("Should be able to get the events instead got %v",
			err)
	}
	ans2, err := dep.GetEventsDuring(
		strconv.Itoa(int(start.Unix())),
		strconv.Itoa(int(end.Unix())))
	if err != nil {
		t.Errorf("Should be able to get the events instead got %v",
			err)
	}

	if len(ans1) != len(ans2) {
		t.Errorf("Should see same result from both the cache and the DB instead got %d, %d", len(ans1), len(ans2))
	}
}


func TestDeployService_GetEventStats(t *testing.T) {
	dep := getDeployService(t)
	date := time.Date(2017, time.October, 29, 0, 0, 0, 0, time.UTC)
	ans1, err := dep.GetEventStats(strconv.Itoa(int(date.Unix())))
	if err != nil {
		t.Errorf("Should return the event stats instead got, %v", err)
	}
	if len(ans1) != 1 {
		t.Errorf("There should be 1 day worth of data instead got %d", len(ans1))
	}
	ans2, err := dep.GetEventStats("")
	if err != nil {
		t.Errorf("Should return the event stats instead got, %v", err)
	}
	if len(ans2) != 8 {
		t.Errorf("There should be 8 days worth of data instead got %d", len(ans2))
	}

}
