package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

const (
	dbDriver = "mysql"
	// dbSource = "local:!Elnino1903@tcp(localhost:3306)/simple_bank"

	dbSource = "root:secret@tcp(localhost:3356)/simple_bank?parseTime=true"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to dbL:", err)
	}
	testQueries = New(testDB)
	log.Output(1, "connect successfully")

	os.Exit(m.Run())
}
