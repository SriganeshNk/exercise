package Database

import (
	_ "github.com/mattn/go-sqlite3"
	"testing"
)

func TestDeployDB_NewDataBase(t *testing.T) {
	db, err := NewDataBase("../deploys.sqlite")
	defer db.Close()
	if err != nil {
		t.Errorf("The Database should be present, instead got %v", err)
	}
}

func TestDeployDB_Query(t *testing.T) {
	db, err := NewDataBase("../deploys.sqlite")
	defer db.Close()
	if err != nil {
		t.Errorf("The Database should be present, instead got %v", err)
	}
	result, err := db.Query("select count(*) from deploys")
	var x int64
	for result.Next() {
		result.Scan(&x)
	}
	if x != 1000 {
		t.Errorf("There should be 100 records present %d", x)
	}
}
