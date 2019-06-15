package Database

import (
	"database/sql"
	"log"
)

type DeployDB struct {
	database *sql.DB
	source string
}

func NewDataBase(conn string) (*DeployDB, error) {
	db, err := sql.Open("sqlite3", conn)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	newDB := DeployDB{db, conn}
	return &newDB, nil
}


func (db *DeployDB) Query(query string) (*sql.Rows, error) {
	result, err := db.database.Query(query)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return result, err
}


func (db *DeployDB) Close() {
	db.database.Close()
}